package arex

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func Test_BuildArexReport(t *testing.T) {
	dataJSON, err := ioutil.ReadFile("../testdata/grafana.json")
	SchemaJSON, err := ioutil.ReadFile("../testdata/grafana_schema_2.json")
	if err != nil {
		log.Fatal(err)
	}

	var val validation
	val.Key = ""
	val.Schema = string(SchemaJSON)
	val.Input = string(dataJSON)
	res, err := json.Marshal(val)
	fmt.Printf("\n%s\n", res)
}

func Test_NewSchema(t *testing.T) {
	var cp comparing
	dataJSON, _ := ioutil.ReadFile("../testdata/grafana.json")
	dataJSONy, _ := ioutil.ReadFile("../testdata/grafana1.json")

	cp.ValueX = string(dataJSON)
	cp.ValueY = string(dataJSONy)
	cp.Options = ""
	res, _ := json.Marshal(cp)
	fmt.Printf("\n%s\n", res)
}

func Test_mongodb(t *testing.T) {
	var item schemaStore
	item.Key = "aaaaaaaaa"
	item.Schema = "{\"aa\":\"bb\"}"
	saveSchema(context.Background(), item)
}
