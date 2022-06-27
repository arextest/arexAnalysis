package arex

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// Title, describes what the test intends to verify and includes an identification number;
// Description, explains the test objective;
// Preconditions, including any configurations, data setup, tests QA professionals must execute first, etc.;
// References, including links to user stories to ensure traceability;
// Detailed steps, clearly written and easy to follow; and
// Expected results, clearly document how the application should respond at each step.
type testcase struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Service       string `json:"service"`
	InterfaceName string `json:"api"`

	Protocol      string                 `json:"protocol,omitempty"`
	Methods       string                 `json:"method,omitempty"`
	URI           string                 `json:"uri,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
	Authorization string                 `json:"author,omitempty"`
	Headers       map[string]interface{} `json:"header,omitempty"`
	Body          string                 `json:"body,omitempty"`
	BodyFormat    string                 `json:"format,omitempty"`
	PreScript     string                 `json:"pre-script,omitempty"`
	Tests         string                 `json:"tests,omitempty"`
	Settings      []string               `json:"settings,omitempty"`
}

func (t *testcase) ToCaseText() string {
	tmpl, err := template.ParseFiles("../template/golang.tmpl")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, t)
	return buf.String()
}

func getTestCases(appid, start string) []*testcase {
	convertSevletToTestCase := func(s *servletmocker) *testcase {
		var tc testcase
		tc.Title = s.AppID + "_" + strings.ReplaceAll(s.Path, "/", "_")
		tc.Description = s.ID + ":" + s.CreateTime.GoString()
		tc.Service = s.Method + "_" + strings.ReplaceAll(s.Path, "/", "_")
		return &tc
	}
	var startTime time.Time
	startTime, err := time.Parse("2022-02-22", start)
	if err != nil {
		startTime = time.Time{}
	}

	rl := queryServletmocker(context.TODO(), appid, startTime)
	tcs := make([]*testcase, 0)
	for _, oneSevlet := range rl {
		oCase := convertSevletToTestCase(oneSevlet)
		tcs = append(tcs, oCase)
	}
	return tcs
}
