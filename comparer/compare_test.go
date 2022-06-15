package comparer

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"encoding/hex"
	"encoding/json"
	encoding_json "encoding/json"

	gojson "github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"

	jsoniter "github.com/json-iterator/go"
)

func Test_Iterjson(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/baseMsgUn.json")
	if err != nil {
		fmt.Print(err)
	}
	t1 := time.Now()
	jsonObject := make(map[string]interface{})
	err = encoding_json.Unmarshal(data, &jsonObject)
	if err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)

	t1 = time.Now()
	jsonObject1 := make(map[string]interface{})
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary
	iterjson.Unmarshal(data, &jsonObject1)
	elapsed = time.Since(t1)
	fmt.Println("App elapsed-jsoniter: ", elapsed)
}

func Test_gojson(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/baseMsgUn.json")
	if err != nil {
		fmt.Print(err)
	}
	t1 := time.Now()
	jsonObject := make(map[string]interface{})
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary
	iterjson.Unmarshal(data, &jsonObject)
	// iterjson.UnmarshalFromString()
	it := iterjson.BorrowIterator(data)
	it.WhatIsNext()

	elapsed := time.Since(t1)
	fmt.Println("App elapsed-jsoniter: ", elapsed)

	t1 = time.Now()
	jsonObject1 := make(map[string]interface{})
	gojson.Unmarshal(data, &jsonObject1)
	elapsed = time.Since(t1)
	fmt.Println("App elapsed-gojson: ", elapsed)
}

func Test_Method1(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/baseMsgUn.json")
	if err != nil {
		fmt.Print(err)
	}
	t1 := time.Now()
	mapJSON := make(map[string]interface{})
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary
	iterjson.Unmarshal(data, &mapJSON)
	elapsed := time.Since(t1)
	fmt.Println("Unmarshal time: ", elapsed)

	t1 = time.Now()
	res1, err := iterjson.MarshalToString(mapJSON)
	elapsed = time.Since(t1)
	fmt.Println(res1)
	fmt.Println("Marshal time: ", elapsed)
}

func Test_JSONTree(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/baseMsgUn.json")
	if err != nil {
		fmt.Print(err)
	}
	t1 := time.Now()
	mapJSON := make(map[string]interface{})
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary
	iterjson.Unmarshal(data, &mapJSON)
	elapsed := time.Since(t1)
	fmt.Println("Unmarshal time: ", elapsed)

	t1 = time.Now()
	tree := NewJSONTree()
	err = tree.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(err)
	}
	elapsed = time.Since(t1)
	fmt.Println("Tree time: ", elapsed)
}

func Test_fnv(t *testing.T) {
	a := fnv.New64a()
	a.Write([]byte("hello"))
	fmt.Println(hex.EncodeToString(a.Sum(nil)))
	// fmt.Println(a.Sum64)
}

func Test_typeof(t *testing.T) {
	var x float64 = 3.4
	fmt.Println("type:", reflect.TypeOf(x))

	b := []int{1, 2, 3, 4}
	fmt.Println("type b:", reflect.TypeOf(b))
}

func Test_jsonsort(t *testing.T) {
	maparoot := make(map[string]any)
	a := make([][]int, 5)
	a[0] = []int{10, 0}
	a[1] = []int{1, 1}
	a[2] = []int{3, 1}
	a[3] = []int{1, 3}
	a[4] = []int{3, 5}
	maparoot["a"] = a
	bytea, _ := json.Marshal(maparoot)
	fmt.Println(string(bytea))
	mapa := make(map[string]any)
	json.Unmarshal(bytea, &mapa)

	mapbroot := make(map[string]any)
	b := make([][]int, 5)
	b[0] = []int{1, 1}
	b[1] = []int{0, 10}
	b[2] = []int{5, 3}
	b[3] = []int{1, 3}
	b[4] = []int{3, 1}
	mapbroot["a"] = b
	byteb, _ := json.Marshal(mapbroot)
	fmt.Println(string(byteb))
	mapb := make(map[string]any)
	json.Unmarshal(byteb, &mapb)

	if reflect.DeepEqual(mapa, mapb) {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}

	if cmp.Equal(mapa, mapb) {
		fmt.Println("1 equal")
	} else {
		fmt.Println("1 not equal")
	}

}

func Test_objectsort(t *testing.T) {
	maparoot := make(map[string]any)
	a := make([]Diff, 5)
	a0 := Diff{}
	a0.XPath = Stack{}
	a0.XPath.Push("1")
	a0.XPath.Push("2")
	a0.BasicLog = "basic0"
	a0.Alogs = "alog0"
	a0.BBasicLogs = "bbasic0"
	a0.BLogs = "blog0"
	a[0] = a0
	a1 := Diff{}
	a1.XPath = Stack{}
	a1.XPath.Push("1")
	a1.XPath.Push("2")
	a1.BasicLog = "basic1"
	a1.Alogs = "alog1"
	a1.BBasicLogs = "bbasic1"
	a1.BLogs = "blog1"
	a[1] = a1
	a2 := Diff{}
	a2.XPath = Stack{}
	a2.XPath.Push("a")
	a2.XPath.Push("b")
	a2.BasicLog = "basic2"
	// a2.Alogs = "alog2"
	a2.Alogs = "alog2xxxxxxxxxxxxx"
	a2.BBasicLogs = "bbasic2"
	a2.BLogs = "blog2"
	a[2] = a2
	a3 := Diff{}
	a3.XPath = Stack{}
	a3.XPath.Push("a")
	a3.XPath.Push("b")
	a3.BasicLog = "basic3"
	a3.Alogs = "alog3"
	a3.BBasicLogs = "bbasic3"
	a3.BLogs = "blog3"
	a[3] = a3
	a4 := Diff{}
	a4.XPath = Stack{}
	a4.XPath.Push("a")
	a4.XPath.Push("b")
	a4.BasicLog = "basic4"
	a4.Alogs = "alog4"
	a4.BBasicLogs = "bbasic4"
	a4.BLogs = "blog4"
	a[4] = a4
	maparoot["a"] = a
	bytea, _ := json.Marshal(maparoot)
	fmt.Println(string(bytea))
	mapa := make(map[string]any)
	json.Unmarshal(bytea, &mapa)

	mapbroot := make(map[string]any)
	b := make([]Diff, 5)
	b0 := Diff{}
	b0.XPath = Stack{}
	b0.XPath.Push("1")
	b0.XPath.Push("2")
	b0.BasicLog = "basic0"
	b0.Alogs = "alog0"
	b0.BBasicLogs = "bbasic0"
	b0.BLogs = "blog0"
	b1 := Diff{}
	b1.XPath = Stack{}
	b1.XPath.Push("1")
	b1.XPath.Push("2")
	b1.BasicLog = "basic1"
	b1.Alogs = "alog1"
	b1.BBasicLogs = "bbasic1"
	b1.BLogs = "blog1"
	b2 := Diff{}
	b2.XPath = Stack{}
	b2.XPath.Push("a")
	b2.XPath.Push("b")
	b2.BasicLog = "basic2"
	b2.Alogs = "alog2"
	b2.BBasicLogs = "bbasic2"
	b2.BLogs = "blog2"
	b3 := Diff{}
	b3.XPath = Stack{}
	b3.XPath.Push("a")
	b3.XPath.Push("b")
	b3.BasicLog = "basic3"
	b3.Alogs = "alog3"
	b3.BBasicLogs = "bbasic3"
	b3.BLogs = "blog3"
	b4 := Diff{}
	b4.XPath = Stack{}
	b4.XPath.Push("a")
	b4.XPath.Push("b")
	b4.BasicLog = "basic4"
	b4.Alogs = "alog4"
	b4.BBasicLogs = "bbasic4"
	b4.BLogs = "blog4"
	b[0] = b2
	b[1] = b3
	b[2] = b4
	b[3] = b1
	b[4] = b0
	mapbroot["a"] = b
	byteb, _ := json.Marshal(mapbroot)
	fmt.Println(string(byteb))
	mapb := make(map[string]any)
	json.Unmarshal(byteb, &mapb)

	if reflect.DeepEqual(mapa, mapb) {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}

	if cmp.Equal(mapa, mapb) {
		fmt.Println("1 equal")
	} else {
		fmt.Println("1 not equal")
	}

}

func Test_Compare2Json(t *testing.T) {
	res := GoCmpDiffByFile("../testdata/baseMsgUn.json", "../testdata/testMsgUn.json")
	fmt.Println(res)
}

func Test_Compare2Schema(t *testing.T) {
	res := GoCmpDiffByFile("../testdata/grafana_schema.json", "../testdata/grafana_schema_1.json")
	fmt.Println(res)
}
