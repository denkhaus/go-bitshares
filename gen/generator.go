package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/objx"
	tomb "gopkg.in/tomb.v2"

	// importing this initializes sample data fetching
	"github.com/denkhaus/bitshares/gen/samples"
)

const (
	samplesDir    = "samples"
	operationsDir = "operations"
)

var (
	sampleDataTemplate *template.Template
	sampleMainTemplate *template.Template
	sampleMetaTemplate *template.Template
	genChan            = make(chan GenData, 200)
	tb                 = tomb.Tomb{}
)

type GenData struct {
	Type      types.OperationType
	Block     int
	SampleIdx int
	Data      objx.Map
}

func main() {
	defer close(genChan)

	logging.Info("connect api")
	api := bitshares.NewWebsocketAPI(tests.WsFullApiUrl)
	if err := api.Connect(); err != nil {
		handleError(errors.Annotate(err, "Connect"))
	}

	api.OnError(func(err error) {
		handleError(errors.Annotate(err, "OnError"))
	})

	logging.Info("parse templates")

	tmpl, err := template.ParseFiles("templates/meta.go.tmpl")
	if err != nil {
		handleError(errors.Annotate(err, "ParseFiles [meta]"))
	}
	sampleMetaTemplate = tmpl

	tmpl, err = template.ParseFiles("templates/opsamplemain.go.tmpl")
	if err != nil {
		handleError(errors.Annotate(err, "ParseFiles [main]"))
	}
	sampleMainTemplate = tmpl

	tmpl, err = template.ParseFiles("templates/opsampledata.go.tmpl")
	if err != nil {
		handleError(errors.Annotate(err, "ParseFiles [data]"))
	}
	sampleDataTemplate = tmpl

	// start generate goroutine
	tb.Go(func() error {
		return generate(genChan)
	})

	dataStore := NewOpDataStore()
	if err := dataStore.Init(data.OpSampleMap, genChan); err != nil {
		handleError(errors.Annotate(err, "init datastore"))
	}

	block := uint64(samples.LastScannedBlock)

	logging.Infof("loop blocks, starting from %d", block)

	for tb.Alive() {
		resp, err := api.CallWsAPI(0, "get_block", block)
		if err != nil {
			handleError(errors.Annotate(err, "GetBlock"))
		}

		var data map[string]interface{}
		if ffjson.Unmarshal(*resp, &data); err != nil {
			handleError(errors.Annotate(err, "Unmarshal [resp]"))
		}

		m := objx.New(data)
		trxs := m.Get("transactions")

		// enumerate Transactions
		trxs.EachInter(func(_ int, trx interface{}) bool {
			ops := objx.New(trx).Get("operations")
			// enumerate Operations
			ops.EachInter(func(_ int, o interface{}) bool {
				op := o.([]interface{})
				opType := types.OperationType(types.Int8(op[0].(float64)))
				opData := objx.New(op[1])

				blob := NewOperationBlob(opData)
				idx, err := dataStore.Insert(opType, blob, block)
				if err != nil {
					handleError(errors.Annotate(err, "Evaluate"))
				}

				if idx >= 0 && tb.Alive() {
					genChan <- GenData{
						Type:      opType,
						SampleIdx: idx,
						Data:      opData,
					}
				}

				return true
			})

			return true
		})

		if err := generateMetaFile(block); err != nil {
			handleError(errors.Annotate(err, "generateMetaFile"))
		}

		block++
	}

	if err := tb.Err(); err != nil {
		handleError(errors.Annotate(err, "main"))
	}
}

func generate(ch chan GenData) error {
	for {
		select {
		case data := <-ch:
			if err := generateSampleData(data); err != nil {
				return errors.Annotate(err, "generateSampleData")
			}

			// if err := generateOpData(data); err != nil {
			// 	return errors.Annotate(err, "generateOpData")
			// }

		case <-tb.Dying():
			return nil
		default:
		}
	}
}

func generateSampleDataFile(d GenData, sampleData string) error {
	opName := d.Type.OperationName()

	buf := bytes.NewBuffer(nil)
	err := sampleDataTemplate.Execute(buf, struct {
		SampleDataOpType  string
		SampleData        interface{}
		SampleDataVarName string
		SampleDataIdx     int
	}{
		SampleDataOpType:  d.Type.String(),
		SampleData:        template.HTML(sampleData),
		SampleDataVarName: fmt.Sprintf("sampleData%s", opName),
		SampleDataIdx:     d.SampleIdx,
	})

	if err != nil {
		return errors.Annotate(err, "Execute")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.Annotate(err, "Source")
	}

	fileName := strings.ToLower(fmt.Sprintf("%s/%s_%d.go", samplesDir, opName, d.SampleIdx))
	if err := ioutil.WriteFile(fileName, formatted, 0622); err != nil {
		return errors.Annotate(err, "WriteFile")
	}

	return nil
}

func generateMetaFile(block uint64) error {
	buf := bytes.NewBuffer(nil)
	err := sampleMetaTemplate.Execute(buf, struct {
		LastScannedBlock uint64
	}{
		LastScannedBlock: block,
	})

	if err != nil {
		return errors.Annotate(err, "Execute")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.Annotate(err, "Source")
	}

	fileName := strings.ToLower(fmt.Sprintf("%s/meta.go", samplesDir))
	if err := ioutil.WriteFile(fileName, formatted, 0622); err != nil {
		return errors.Annotate(err, "WriteFile")
	}

	return nil
}

func generateSampleMainFile(d GenData) error {
	opName := d.Type.OperationName()

	buf := bytes.NewBuffer(nil)
	err := sampleMainTemplate.Execute(buf, struct {
		SampleDataOpType  string
		SampleDataVarName string
	}{
		SampleDataOpType:  d.Type.String(),
		SampleDataVarName: fmt.Sprintf("sampleData%s", opName),
	})

	if err != nil {
		return errors.Annotate(err, "Execute")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.Annotate(err, "Source")
	}

	fileName := strings.ToLower(fmt.Sprintf("%s/%s.go", samplesDir, opName))
	if err := ioutil.WriteFile(fileName, formatted, 0622); err != nil {
		return errors.Annotate(err, "WriteFile")
	}

	return nil
}

func generateSampleData(d GenData) error {
	if d.SampleIdx == 0 {
		if err := generateSampleMainFile(d); err != nil {
			return errors.Annotate(err, "generateSampleMainFile")
		}
	}

	sampleDataJSON, err := json.MarshalIndent(d.Data, "", "  ")
	if err != nil {
		return errors.Annotate(err, "MarshalIndent")
	}

	sampleData := fmt.Sprintf("`%s`", sampleDataJSON)

	//update sample map too
	if data.OpSampleMap[d.Type] == nil {
		data.OpSampleMap[d.Type] = make(map[int]string)
	}
	data.OpSampleMap[d.Type][d.SampleIdx] = sampleData

	if err := generateSampleDataFile(d, sampleData); err != nil {
		return errors.Annotate(err, "generateSampleDataFile")
	}

	return nil
}

func handleError(err error) {
	fmt.Println("error: ", errors.ErrorStack(err))

	if tb.Alive() {
		//kill generator goroutine and wait
		tb.Kill(err)
		tb.Wait()
	}

	os.Exit(1)
}
