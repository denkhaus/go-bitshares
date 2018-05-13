package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
	"github.com/stretchr/objx"
)

var (
	sampleDataTemplate *template.Template
	samplesDir         = "samples"
)

func main() {

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

	dataStore := NewOpDataStore()
	if err := dataStore.Init(data.OpSampleMap); err != nil {
		handleError(errors.Annotate(err, "init datastore"))
	}

	//TODO: save last scanned block and reapply
	block := uint64(230000)

	for {
		resp, err := api.CallWsAPI(0, "get_block", block)
		if err != nil {
			handleError(errors.Annotate(err, "GetBlock"))
		}

		m := objx.New(resp)

		trxs := m.Get("transactions")
		trxs.EachInter(func(_ int, trx interface{}) bool {
			ops := objx.New(trx).Get("operations")
			ops.EachInter(func(_ int, o interface{}) bool {
				op := o.([]interface{})
				opType := types.OperationType(types.Int8(op[0].(float64)))
				opData := objx.New(op[1])

				blob := NewOperationBlob(opData)
				ok, err := dataStore.Evaluate(opType, blob, block)
				if err != nil {
					handleError(errors.Annotate(err, "Evaluate"))
				}

				if ok {
					if err := GenerateSampleData(opType, opData); err != nil {
						handleError(errors.Annotate(err, "GenerateSampleData"))
					}
				}

				return true
			})

			return true
		})

		block++
	}
}

func GenerateSampleData(opType types.OperationType, opData objx.Map) error {
	opName := opType.OperationName()
	fileName := fmt.Sprintf("%s/%s.go", samplesDir, opName)
	fileName = strings.ToLower(fileName)

	f, err := os.Create(fileName)
	if err != nil {
		return errors.Annotate(err, "Evaluate")
	}

	defer f.Close()

	sampleDataJSON, err := json.MarshalIndent(opData, "", "  ")
	if err != nil {
		return errors.Annotate(err, "MarshalIndent")
	}

	sampleData := fmt.Sprintf("`%s`", sampleDataJSON)

	err = sampleDataTemplate.Execute(f, struct {
		SampleDataOpType  string
		SampleData        interface{}
		SampleDataVarName string
	}{
		SampleDataOpType:  opType.String(),
		SampleData:        template.HTML(sampleData),
		SampleDataVarName: fmt.Sprintf("sampleData%s", opName),
	})

	if err != nil {
		return errors.Annotate(err, "Execute")
	}

	return nil
}

func handleError(err error) {
	fmt.Println(errors.ErrorStack(err))
	os.Exit(1)
}
