package jsonschema

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type propertySuite struct{}

var _ = Suite(&propertySuite{})

type ExampleJSONBasic struct {
	Omitted    string  `json:"-,omitempty"`
	Bool       bool    `json:",omitempty"`
	Integer    int     `json:",omitempty"`
	Integer8   int8    `json:",omitempty"`
	Integer16  int16   `json:",omitempty"`
	Integer32  int32   `json:",omitempty"`
	Integer64  int64   `json:",omitempty"`
	UInteger   uint    `json:",omitempty"`
	UInteger8  uint8   `json:",omitempty"`
	UInteger16 uint16  `json:",omitempty"`
	UInteger32 uint32  `json:",omitempty"`
	UInteger64 uint64  `json:",omitempty"`
	String     string  `json:",omitempty"`
	Bytes      []byte  `json:",omitempty"`
	Float32    float32 `json:",omitempty"`
	Float64    float64
	Interface  interface{}
	Timestamp  time.Time `json:",omitempty"`
}

func (p *propertySuite) TestLoad(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONBasic{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type:     "object",
			Required: []string{"Float64", "Interface"},
			Properties: map[string]*property{
				"Bool":       {Type: "boolean"},
				"Integer":    {Type: "integer"},
				"Integer8":   {Type: "integer"},
				"Integer16":  {Type: "integer"},
				"Integer32":  {Type: "integer"},
				"Integer64":  {Type: "integer"},
				"UInteger":   {Type: "integer"},
				"UInteger8":  {Type: "integer"},
				"UInteger16": {Type: "integer"},
				"UInteger32": {Type: "integer"},
				"UInteger64": {Type: "integer"},
				"String":     {Type: "string"},
				"Bytes":      {Type: "string"},
				"Float32":    {Type: "number"},
				"Float64":    {Type: "number"},
				"Interface":  {},
				"Timestamp":  {Type: "string", Format: "date-time"},
			},
		},
	})
}

type ExampleJSONBasicWithTag struct {
	Bool bool `json:"test"`
}

func (p *propertySuite) TestLoadWithTag(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONBasicWithTag{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type:     "object",
			Required: []string{"test"},
			Properties: map[string]*property{
				"test": {Type: "boolean"},
			},
		},
	})
}

type SliceStruct struct {
	Value string
}

type ExampleJSONBasicSlices struct {
	Slice            []string      `json:",foo,omitempty"`
	SliceOfInterface []interface{} `json:",foo"`
	SliceOfStruct    []SliceStruct
}

func (p *propertySuite) TestLoadSliceAndContains(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONBasicSlices{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Slice": {
					Type:  "array",
					Items: &property{Type: "string"},
				},
				"SliceOfInterface": {
					Type: "array",
				},
				"SliceOfStruct": {
					Type: "array",
					Items: &property{
						Type:     "object",
						Required: []string{"Value"},
						Properties: map[string]*property{
							"Value": {
								Type: "string",
							},
						},
					},
				},
			},

			Required: []string{"SliceOfInterface", "SliceOfStruct"},
		},
	})
}

type ExampleJSONNestedStruct struct {
	Struct struct {
		Foo string
	}
}

func (self *propertySuite) TestLoadNested(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONNestedStruct{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Struct": {
					Type: "object",
					Properties: map[string]*property{
						"Foo": {Type: "string"},
					},
					Required: []string{"Foo"},
				},
			},
			Required: []string{"Struct"},
		},
	})
}

type EmbeddedStruct struct {
	Foo string
}

type ExampleJSONEmbeddedStruct struct {
	EmbeddedStruct
}

func (self *propertySuite) TestLoadEmbedded(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONEmbeddedStruct{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Foo": {Type: "string"},
			},
			Required: []string{"Foo"},
		},
	})
}

type ExampleJSONBasicMaps struct {
	Maps           map[string]string `json:",omitempty"`
	MapOfInterface map[string]interface{}
}

func (self *propertySuite) TestLoadMap(c *C) {
	j := &SchemaDocument{}
	j.Read(&ExampleJSONBasicMaps{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Maps": {
					Type: "object",
					Properties: map[string]*property{
						".*": {Type: "string"},
					},
					AdditionalProperties: false,
				},
				"MapOfInterface": {
					Type:                 "object",
					AdditionalProperties: true,
				},
			},
			Required: []string{"MapOfInterface"},
		},
	})
}

func (self *propertySuite) TestLoadNonStruct(c *C) {
	j := &SchemaDocument{}
	j.Read([]string{})

	c.Assert(*j, DeepEquals, SchemaDocument{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type:  "array",
			Items: &property{Type: "string"},
		},
	})
}

func (self *propertySuite) TestString(c *C) {
	j := &SchemaDocument{}
	j.Read(true)

	expected := "{\n" +
		"    \"$schema\": \"http://json-schema.org/schema#\",\n" +
		"    \"type\": \"boolean\"\n" +
		"}"

	res, err := j.String()
	if err != nil {
		fmt.Printf("%#v", err)
	}
	c.Assert(res, Equals, expected)
}

func (self *propertySuite) TestMarshal(c *C) {
	j := &SchemaDocument{}
	j.Read(10)

	expected := "{\n" +
		"    \"$schema\": \"http://json-schema.org/schema#\",\n" +
		"    \"type\": \"integer\"\n" +
		"}"

	json, err := j.Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(json), Equals, expected)
}

type EmbeddedType struct {
	Zoo string
}

type Item struct {
	Value string
}

type ExampleBasic struct {
	Foo bool   `json:"foo"`
	Bar string `json:",omitempty"`
	Qux int8
	Baz []string
	EmbeddedType
	List []Item
}

func Test_Generate(test *testing.T) {
	s := &SchemaDocument{}
	s.Read(&ExampleBasic{})
	fmt.Println(s)
}
