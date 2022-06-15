package comparer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

// DifferItem store origin data
type DifferItem struct {
	Path       string   `json:"path,omitempty"`
	StructPath cmp.Path `json:"-"`
	Vx         string   `json:"vx,omitempty"`
	Vy         string   `json:"vy,omitempty"`
}

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path  cmp.Path
	Diffs []*DifferItem
}

// PushStep todo
func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

// Report todo
func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		var d DifferItem
		// d.Path = fmt.Sprintf("%#v", r.path)
		d.Path = r.path.GoString()
		d.StructPath = r.path
		if vx.Kind() != reflect.Invalid {
			d.Vx = fmt.Sprintf("%+v", vx)
		}
		if vy.Kind() != reflect.Invalid {
			d.Vy = fmt.Sprintf("%+v", vy)
		}

		r.Diffs = append(r.Diffs, &d)
	}
}

// PopStep todo
func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	res, err := json.MarshalIndent(r.Diffs, " ", "")
	if err != nil {
		return ""
	}
	return string(res)
}

func reflectDeepEqual(aobject any, bobject any) bool {
	t1 := time.Now()
	res := reflect.DeepEqual(aobject, bobject)
	elapsed := time.Since(t1)
	fmt.Printf("equal time elapsed:%d result %t \r\n", elapsed, res)
	return res
}

func googleDeepEqual(aobject any, bobject any) bool {
	tcomapre := time.Now()
	res := cmp.Equal(aobject, bobject)
	elapsed := time.Since(tcomapre)
	fmt.Printf("go-cmp compared time elapsed:%d us result: %t \r\n", elapsed.Nanoseconds(), res)
	return false
}

// GoCmpDiffByFile ("../baseMsgUn.txt","../testMsgUn.txt")
func GoCmpDiffByFile(afile string, bfile string) *DiffReporter {
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary

	atext, _ := ioutil.ReadFile(afile)
	jsonLeft := make(map[string]interface{})
	iterjson.Unmarshal(atext, &jsonLeft)

	btext, _ := ioutil.ReadFile(bfile)
	jsonRight := make(map[string]interface{})
	iterjson.Unmarshal(btext, &jsonRight)

	return GoCmpDiff(jsonLeft, jsonRight)
}

// GocmpDiffByJSON diff json
func GocmpDiffByJSON(afile string, bfile string) *DiffReporter {
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary

	jsonLeft := make(map[string]interface{})
	iterjson.Unmarshal([]byte(afile), &jsonLeft)

	jsonRight := make(map[string]interface{})
	iterjson.Unmarshal([]byte(bfile), &jsonRight)

	return GoCmpDiff(jsonLeft, jsonRight)
}

// DifferItemByGoCmp return different items
func DifferItemByGoCmp(aobject any, bobject any) *DiffReporter {
	tcomapre := time.Now()
	var r DiffReporter
	cmp.Diff(aobject, bobject, cmp.Reporter(&r))
	elapsed := time.Since(tcomapre)
	fmt.Printf("go-cmp Diff time elapsed:%d us \r\n", elapsed.Nanoseconds())
	return &r
}

// GoCmpDiff compare 2 json and return json result
func GoCmpDiff(jsonX any, jsonY any) *DiffReporter {
	res := DifferItemByGoCmp(jsonX, jsonY)
	return res
}

// CompareSDK s
type CompareSDK struct {
	BasicDiffMap *Diffs
}

// NewCompareSDK new
func NewCompareSDK() *CompareSDK {
	sdk := CompareSDK{}
	sdk.BasicDiffMap = NewDiffs()
	return &sdk
}

func (c *CompareSDK) storeResult(xpath Stack, curPath string, left string, right string) {
	if curPath == "" {
		c.BasicDiffMap.StoreDifferent(xpath, left, right)
	} else {
		xpath.Push(curPath)
		c.BasicDiffMap.StoreDifferent(xpath, left, right)
		xpath.Pop()
	}
}

func (c *CompareSDK) mapcompare(xpath Stack, curPath string, amap map[string]any, bmap map[string]any) bool {
	if curPath != "" {
		xpath.Push(curPath)
		defer xpath.Pop()
	}

	samed := true
	akeys := mapset.NewSet()
	for k := range amap {
		akeys.Add(k)
	}

	bkeys := mapset.NewSet()
	for k := range bmap {
		bkeys.Add(k)
	}

	onlya := akeys.Difference(bkeys)
	for val := range onlya.Iterator().C {
		c.storeResult(xpath, val.(string), fmt.Sprintf("- %s%s/%s\n", xpath.ToString(), val.(string), val), "")
		samed = false
	}

	onlyb := bkeys.Difference(akeys)
	for val := range onlyb.Iterator().C {
		c.storeResult(xpath, val.(string), "", fmt.Sprintf("+ %s%s/%s\n", xpath.ToString(), val.(string), val))
		samed = false
	}

	samekeys := akeys.Intersect(bkeys)
	for val := range samekeys.Iterator().C {
		if !cmp.Equal(amap[val.(string)], bmap[val.(string)]) {
			vala := amap[val.(string)]
			valb := bmap[val.(string)]
			switch vala.(type) {
			case map[string]any:
				c.mapcompare(xpath, val.(string), vala.(map[string]any), valb.(map[string]any))
			case []any:
				c.arraycompare(xpath, val.(string), vala.([]any), valb.([]any))
			default:
				c.storeResult(xpath, val.(string),
					fmt.Sprintf("- %s/%s/%s\n", xpath.ToString(), val.(string), vala),
					fmt.Sprintf("+ %s/%s/%s\n", xpath.ToString(), val.(string), valb))
			}
			samed = false
		}
	}

	return samed
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (c *CompareSDK) arraycompareV2(xpath Stack, curPath string, aarray []any, barray []any) {
	if curPath != "" {
		xpath.Push(curPath)
		defer xpath.Pop()
	}

	alen := len(aarray)
	if alen == 0 {
		c.storeResult(xpath, "", "",
			fmt.Sprintf("+ %s/%s\n", xpath.ToString(), barray))
		return
	}
	blen := len(barray)
	if blen == 0 {
		c.storeResult(xpath, "",
			fmt.Sprintf("- %s/%s\n", xpath.ToString(), aarray), "")
		return
	}
	len := min(alen, blen)

	switch aarray[0].(type) {
	case map[string]any:

		for i := 0; i < len; i++ {
			c.mapcompare(xpath, fmt.Sprintf("[%d]", i), aarray[i].(map[string]any), barray[i].(map[string]any))
		}
	case []any:
		for i := 0; i < len; i++ {
			c.arraycompare(xpath, fmt.Sprintf("[%d]", i), aarray[i].([]any), barray[i].([]any))
		}
	default:
		seta := mapset.NewSetFromSlice(aarray)
		setb := mapset.NewSetFromSlice(barray)

		onlya := seta.Difference(setb)
		for val := range onlya.Iterator().C {
			c.storeResult(xpath, fmt.Sprintf("%s", val),
				fmt.Sprintf("- %s/%s\n", xpath.ToString(), val), "")
		}
		onlyb := setb.Difference(seta)
		for val := range onlyb.Iterator().C {
			c.storeResult(xpath, fmt.Sprintf("%s", val), "",
				fmt.Sprintf("+ %s/%s\n", xpath.ToString(), val))
		}
	}

}

func (c *CompareSDK) arraycompare(xpath Stack, curPath string, aarray []any, barray []any) {
	if curPath != "" {
		xpath.Push(curPath)
		defer xpath.Pop()
	}

	alen := len(aarray)
	if alen == 0 {
		c.storeResult(xpath, "", "",
			fmt.Sprintf("+ %s/%s\n", xpath.ToString(), barray))
		return
	}
	blen := len(barray)
	if blen == 0 {
		c.storeResult(xpath, "",
			fmt.Sprintf("- %s/%s\n", xpath.ToString(), aarray), "")
		return
	}
	len := min(alen, blen)

	switch aarray[0].(type) {
	case map[string]any:
		for i := 0; i < len; i++ {
			c.mapcompare(xpath, fmt.Sprintf("[%d]", i), aarray[i].(map[string]any), barray[i].(map[string]any))
		}
	case []any:
		for i := 0; i < len; i++ {
			c.arraycompare(xpath, fmt.Sprintf("[%d]", i), aarray[i].([]any), barray[i].([]any))
		}
	default:
		seta := mapset.NewSetFromSlice(aarray)
		setb := mapset.NewSetFromSlice(barray)

		onlya := seta.Difference(setb)
		for val := range onlya.Iterator().C {
			c.storeResult(xpath, fmt.Sprintf("%s", val),
				fmt.Sprintf("- %s/%s\n", xpath.ToString(), val), "")
		}
		onlyb := setb.Difference(seta)
		for val := range onlyb.Iterator().C {
			c.storeResult(xpath, fmt.Sprintf("%s", val), "",
				fmt.Sprintf("+ %s/%s\n", xpath.ToString(), val))
		}
	}

}

func (c *CompareSDK) compare(basicfile []byte, replayfile []byte) {
	var xpath Stack

	t1 := time.Now()
	jsonLeft := make(map[string]interface{})
	var iterjson = jsoniter.ConfigCompatibleWithStandardLibrary
	iterjson.Unmarshal(basicfile, &jsonLeft)
	elapsed := time.Since(t1)
	fmt.Println("read A file to JSON elapsed:", elapsed)

	t1 = time.Now()
	jsonRight := make(map[string]interface{})
	iterjson.Unmarshal(replayfile, &jsonRight)
	elapsed = time.Since(t1)
	fmt.Println("read B file to JSON elapsed: ", elapsed)

	// if deepequal(jsonLeft, jsonRight) {
	// 	fmt.Println("AB Test success.")
	// 	return
	// }

	if googleDeepEqual(jsonLeft, jsonRight) {
		fmt.Println("AB Test success.")
		return
	}

	t1 = time.Now()
	c.mapcompare(xpath, "", jsonLeft, jsonRight)
	elapsed = time.Since(t1)
	fmt.Println("mapcompare time elapsed: ", elapsed)

	// t1 = time.Now()
	// elapsed = time.Since(t1)
	// fmt.Println("!=: ", elapsed)
}

func comparefile(basicfile string, replayfile string, cfile string) []byte {
	data, err := ioutil.ReadFile(basicfile)
	if err != nil {
		fmt.Print(err)
	}

	data1, err := ioutil.ReadFile(replayfile)
	if err != nil {
		fmt.Print(err)
	}

	data2, err := ioutil.ReadFile(cfile)
	if err != nil {
		fmt.Print(err)
	}

	return compare(data, data1, data2)
}

// Compare func
func compare(basicfile []byte, afile []byte, bfile []byte) []byte {
	aSDK := NewCompareSDK()
	achan := make(chan Diffs)
	defer close(achan)
	go func() {
		aSDK.compare(basicfile, afile)
		dif := aSDK.BasicDiffMap
		achan <- *dif
	}()

	bSDK := NewCompareSDK()
	bchan := make(chan Diffs)
	defer close(bchan)
	go func() {
		bSDK.compare(afile, bfile)
		dif := bSDK.BasicDiffMap
		bchan <- *dif
	}()

	diffa := <-achan
	diffb := <-bchan
	diffa.CombineDifferentAB(&diffb)
	return diffa.AssertTrue()
}
