package arex

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	dog "github.com/DataDog/zstd"
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

func Test_Unzip(t *testing.T) {
	mybytes, err := base64.StdEncoding.DecodeString("H4sIAAAAAAAAAF2OwQ6CMBBE/6VnbEAJgieNXrgYIyb2ZjZlU2pKIZSaEMK/u9WL8TSZmZednVmDULPdzLwmYY9bJUShnq7P3DbepCxi0mi0Y3milpzuSbOEJ3ydpDyJi0A0YC0aKspzpUekyIBVHhRShnYlRKD8MKCV0we7UuA6P8hA3A8Xsj1MLQ0d/zHTSTC/h2p8aYnfh5aIYQs6bINrwIEdfZ7vVci47Fq2vAGPNtrv4gAAAA==")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	res, err := gzipDecompress(mybytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func Test_task(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	batchGenerateByCtripAREX(context.TODO())
}

func Test_ZSTD(t *testing.T) {
	dataJSON, _ := ioutil.ReadFile("../testdata/zstd.txt")
	data, err := base64.StdEncoding.DecodeString(string(dataJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	unzipStr, err := dog.Decompress(nil, data)
	fmt.Println(unzipStr)

	batchGenerateSchema(context.Background(), time.Time{})

}

func Test_CaseGenerate(t *testing.T) {
	res := getTestCases()
	for _, oneCase := range res {
		fmt.Println(oneCase.ToCaseText())
	}

}
