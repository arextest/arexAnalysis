package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"testing"

	"github.com/arextest/arexAnalysis/comparer"
)

func Test_gojson(t *testing.T) {
	sch, err := Compile("../testdata/grafana_schema_2.json")
	if err != nil {
		fmt.Printf("%#v", err)
		log.Fatalf("%#v", err)
	}

	data, err := ioutil.ReadFile("../testdata/grafana.json")
	if err != nil {
		log.Fatal(err)
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatal(err)
	}

	if err = sch.Validate(v); err != nil {
		log.Fatalf("%#v", err)
	}

}

func Test_SchemaCompare(t *testing.T) {
	schx, err := Compile("../testdata/grafana_schema_1.json")
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	schy, err := Compile("../testdata/grafana_schema.json")
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	res := comparer.DifferItemByGoCmp(*schx, *schy)
	fmt.Println(res)

}

func Test_GenerateCase(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/grafana.json")
	if err != nil {
		log.Fatal(err)
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatal(err)
	}

}

func Test_GenerateSchema(test *testing.T) {
	data, err := ioutil.ReadFile("../testdata/grafana.json")
	if err != nil {
		log.Fatal(err)
	}

	f, err := ParseJson(data)
	if err != nil {
		log.Fatal(err)
	}
	schemaDoc, _ := SchemaGenerateGo(f, "Grafana")
	schemaText, err := schemaDoc.Document.String()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("json-Schema:\n%v\n", schemaText)
	ioutil.WriteFile("../testdata/grafana_schema_1.json", []byte(schemaText), fs.ModePerm)
}

func Test_GenerateSchema1(test *testing.T) {
	datax, err := ioutil.ReadFile("../testdata/grafana_schema_1.json")
	if err != nil {
		log.Fatal(err)
	}
	var schemaX = SchemaDocument{}
	json.Unmarshal(datax, &schemaX)

	datay, err := ioutil.ReadFile("../testdata/grafana_schema.json")
	if err != nil {
		log.Fatal(err)
	}
	var schemaY = SchemaDocument{}
	json.Unmarshal(datay, &schemaY)
	err = schemaX.MergeSchemaDocument(&schemaY)
	if err != nil {
		return
	}

	nData, err := json.MarshalIndent(schemaX, " ", " ")
	fmt.Printf("json-Schema:\n%v\n", string(nData))
	ioutil.WriteFile("../testdata/grafana_schema_2.json", nData, fs.ModePerm)
}
