package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/gojson"
	"github.com/juju/errors"
	"github.com/stretchr/objx"
	"gopkg.in/tomb.v2"

	// import this because of initialization of data.OpSampleMap
	_ "github.com/denkhaus/bitshares/gen/samples"
)

type Unmarshalable interface {
	UnmarshalJSON(input []byte) error
}

var (
	sampleDataTemplate *template.Template
	samplesDir         = "samples"
	operationsDir      = "operations"
	genChan            = make(chan GenData, 40)
	tb                 = tomb.Tomb{}

	// do not change order here
	knownTypes = []Unmarshalable{
		&types.AccountOptions{},
		&types.Asset{},
		&types.AssetAmount{},
		&types.AssetFeed{},
		&types.AssetOptions{},
		&types.GrapheneID{},
		&types.BitAssetDataOptions{},
		&types.Authority{},
		&types.Memo{},
		&types.Price{},
		&types.PriceFeed{},
		&types.Votes{},
		&types.PublicKey{},
		&types.Account{},
	}
)

type GenData struct {
	Type types.OperationType
	Data objx.Map
}

func main() {

	defer close(genChan)

	api := api.New(tests.WsFullApiUrl, tests.RpcApiUrl)
	if err := api.Connect(); err != nil {
		handleError(errors.Annotate(err, "Connect"))
	}

	api.OnError(func(err error) {
		handleError(errors.Annotate(err, "OnError"))
	})

	//init templates
	tmpl, err := template.ParseFiles("templates/opsampledata.go.tmpl")
	if err != nil {
		handleError(errors.Annotate(err, "ParseFiles"))
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

	//TODO: save last scanned block and reapply
	block := uint64(1248202)

	fmt.Println("loop blocks")

	for tb.Alive() {
		resp, err := api.CallWsAPI(0, "get_block", block)
		if err != nil {
			handleError(errors.Annotate(err, "GetBlock"))
		}

		m := objx.New(resp)

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
				ok, err := dataStore.Evaluate(opType, blob, block)
				if err != nil {
					handleError(errors.Annotate(err, "Evaluate"))
				}

				if ok && tb.Alive() {
					genChan <- GenData{
						Type: opType,
						Data: opData,
					}
				}

				return true
			})

			return true
		})

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
			fmt.Println("generate!")
			if err := generateSampleData(data); err != nil {
				return errors.Annotate(err, "generateSampleData")
			}

			// blocking?
			if err := generateOpData(data); err != nil {
				return errors.Annotate(err, "generateOpData")
			}

		case <-tb.Dying():
			return nil
		default:
		}
	}
}

func generateOpData(d GenData) error {
	sample := data.GetSampleByType(d.Type)
	if sample == "" {
		return nil
	}

	sample, err := strconv.Unquote(sample)
	if err != nil {
		return errors.Annotate(err, "Unquote")
	}

	fmt.Printf("generate struct by sample %+v\n", sample)

	buf, err := gojson.GenerateWithTypeGuessing(
		strings.NewReader(sample),
		gojson.ParseJson, d.Type.OperationName(),
		"operations", []string{"json"}, true, true,
		guessStructType,
	)

	if err != nil {
		return errors.Annotate(err, "GenerateWithTypeGuessing")
	}

	fmt.Println("generated struct ", string(buf))
	return nil
}

func generateSampleData(d GenData) error {
	opName := d.Type.OperationName()
	fileName := fmt.Sprintf("%s/%s.go", samplesDir, opName)
	fileName = strings.ToLower(fileName)

	f, err := os.Create(fileName)
	if err != nil {
		return errors.Annotate(err, "Evaluate")
	}

	defer f.Close()

	sampleDataJSON, err := json.MarshalIndent(d.Data, "", "  ")
	if err != nil {
		return errors.Annotate(err, "MarshalIndent")
	}

	sampleData := fmt.Sprintf("`%s`", sampleDataJSON)

	//update sample map too
	data.OpSampleMap[d.Type] = sampleData

	// formatted, err := format.Source([]byte(src))
	// if err != nil {
	// 	err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	// }

	err = sampleDataTemplate.Execute(f, struct {
		SampleDataOpType  string
		SampleData        interface{}
		SampleDataVarName string
	}{
		SampleDataOpType:  d.Type.String(),
		SampleData:        template.HTML(sampleData),
		SampleDataVarName: fmt.Sprintf("sampleData%s", opName),
	})

	if err != nil {
		return errors.Annotate(err, "Execute")
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

func guessStructType(value interface{}, suggestedType string) (string, error) {
	util.DumpJSON("valueToGuess", value)
	//	util.DumpJSON("suggestedType", suggestedType)

	for _, typ := range knownTypes {
		v, err := ffjson.Marshal(value)
		if err != nil {
			return "", errors.Annotate(err, "Marshal")
		}

		if err := typ.UnmarshalJSON(v); err == nil {
			name := reflect.ValueOf(typ).Type().Name()
			util.Dump("solved", typ)
			util.Dump("typeName", name)
			return name, nil
		}
	}

	return suggestedType, nil
}
