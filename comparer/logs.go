package comparer

import (
	"encoding/json"
	"fmt"
)

// Diff Store different string
type Diff struct {
	XPath Stack `json:"xpath"`
	// 基准报文.就是记录下来的报文
	BasicLog   string `json:"basiclog"`
	Alogs      string `json:"alog"`
	BBasicLogs string `json:"abasic,omitempty"`
	BLogs      string `json:"blog,omitempty"`
	asserted   bool
	ignored    bool
}

// Diffs all diff
type Diffs struct {
	maps map[string]Diff
}

// NewDiffs result
func NewDiffs() *Diffs {
	df := Diffs{maps: make(map[string]Diff)}
	return &df
}

// StoreDifferent store basic ,a
func (d *Diffs) StoreDifferent(xpath Stack, basiclog string, alog string) {
	key := xpath.ToString()
	if _, _ok := d.maps[key]; _ok {
		fmt.Printf("error: double key %s\n", key)
		return
	}

	tempDiff := Diff{XPath: xpath, BasicLog: basiclog, Alogs: alog, BLogs: "", BBasicLogs: "", asserted: true, ignored: false}
	d.maps[key] = tempDiff
}

// StoreDifferentAB store B
func (d *Diffs) StoreDifferentAB(xpath Stack, alog string, blog string) {
	key := xpath.ToString()
	if _, _ok := d.maps[key]; _ok {
		tempDiff := d.maps[key]
		tempDiff.BBasicLogs = alog
		tempDiff.BLogs = blog
		tempDiff.asserted = false
		d.maps[key] = tempDiff
		return
	}

	tempDiff := Diff{XPath: xpath, BasicLog: "", Alogs: "", BLogs: blog, BBasicLogs: alog, asserted: true, ignored: false}
	d.maps[key] = tempDiff
}

// CombineDifferentAB (basic<->A) union (A<-->B)
func (d *Diffs) CombineDifferentAB(o *Diffs) {
	for key, v := range o.maps {
		if _, _ok := d.maps[key]; _ok {
			tempDiff := d.maps[key]
			tempDiff.BBasicLogs = v.BasicLog
			tempDiff.BLogs = v.Alogs
			tempDiff.asserted = false
			d.maps[key] = tempDiff
			continue
		}

		tempDiff := Diff{XPath: v.XPath, BasicLog: "", Alogs: "", BLogs: v.BLogs, BBasicLogs: v.BasicLog, asserted: true, ignored: false}
		d.maps[key] = tempDiff
	}
}

// StoreAssert store Asert Info 废弃
func (d *Diffs) StoreAssert(xpath Stack, assertlog string) {
	key := xpath.ToString()
	if _, _ok := d.maps[key]; _ok {
		tempDiff := d.maps[key]
		tempDiff.BBasicLogs = assertlog
		d.maps[key] = tempDiff
		return
	}

	tempDiff := Diff{XPath: xpath, Alogs: "", BLogs: "", BBasicLogs: "", BasicLog: assertlog, asserted: true, ignored: false}
	d.maps[key] = tempDiff
}

// AssertTrue print
func (d *Diffs) AssertTrue() []byte {
	output := make(map[string]Diff)
	for k, v := range d.maps {
		diff := v
		if diff.asserted {
			output[k] = v
		}
	}
	jsonbyte, _ := json.MarshalIndent(output, "", " ")
	// fmt.Println(string(jsonbyte))
	return jsonbyte
}
