package jsonschema

import (
	"encoding/json"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// DefaultSchemaText text for schema
const DefaultSchemaText = "http://json-schema.org/schema#"

// SchemaDocument schema struct
type SchemaDocument struct {
	Schema string `json:"$schema,omitempty"`

	property
}

// Reads the variable structure into the JSON-Schema Document
func (d *SchemaDocument) Read(variable interface{}) {
	d.setDefaultSchema()

	value := reflect.ValueOf(variable)
	d.read(value.Type(), tagOptions(""))
}

// MergeSchemaDocument merget y to current Schema Document
func (d *SchemaDocument) MergeSchemaDocument(y *SchemaDocument) error {
	if cmp.Equal(d.property, y.property) {
		return nil
	}
	return mergeProperty(&d.property, &y.property)

}

func mergeProperty(x *property, y *property) error {
	if cmp.Equal(*x, *y) {
		return nil
	}

	if y.Type != "" {
		x.Type = y.Type
	}
	if y.Format != "" {
		x.Format = y.Format
	}

	if y.Items != nil {
		if x.Items == nil {
			x.Items = y.Items
		} else {
			err := mergeProperty(x.Items, y.Items)
			if err != nil {
				return err
			}
		}
	}

	if y.Properties != nil {
		if x.Properties == nil {
			x.Properties = y.Properties
		}
		for key, value := range y.Properties {
			if _, ok := x.Properties[key]; ok {
				err := mergeProperty(x.Properties[key], value)
				if err != nil {
					return err
				}
			} else {
				x.Properties[key] = value
			}
		}
	}

	if y.Required != nil {
		if x.Required == nil {
			x.Required = y.Required
		} else {
			x.Required = intersect(x.Required, y.Required)
		}
	}

	return nil
}

// Compile the json data to json-schema
func (d *SchemaDocument) Compile(variable interface{}) {
	d.setDefaultSchema()

	value := reflect.ValueOf(variable)
	d.read(value.Type(), tagOptions(""))
}

func (d *SchemaDocument) setDefaultSchema() {
	if d.Schema == "" {
		d.Schema = DefaultSchemaText
	}
}

// Marshal returns the JSON encoding of the Document
func (d *SchemaDocument) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "    ")
}

// String return the JSON encoding of the Document as a string
func (d *SchemaDocument) String() (string, error) {
	json, err := d.Marshal()
	return string(json), err
}

//
// Type存放类型
// format存放格式数据
// items 存放array的类型
// properties 存放object的字段类型
// required 存放当前节点下
// Draft 6 中的新内容 examples
// Draft 6 中的新内容 const关键字被用于限制值为一个常量值。
// Draft 7 中的新内容 布尔类型的关键字readOnly和writeOnly
// Draft 7 中的新内容 $comment关键字严格用于向模式添加注释。它的值必须始终是一个字符串。
// Draft2019-09的新内容 deprecated关键字是一个布尔值
// Items有单例和多例[]*property,暂时只实现了单例
type property struct {
	Location string `json:"-"` // absolute location

	// dynamicAnchors []*property

	// annotations. captured only when Compiler.ExtractAnnotations is true.
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Comment     string        `json:"comment,omitempty"`
	Examples    []interface{} `json:"examples,omitempty"`
	Deprecated  bool          `json:"deprecated,omitempty"`
	ReadOnly    bool          `json:"readOnly,omitempty"`
	WriteOnly   bool          `json:"writeOnly,omitempty"`

	// type agnostic validations
	Format          string        `json:"format,omitempty"`
	Always          *bool         `json:"-"` // always pass/fail. used when booleans are used as schemas in draft-07.
	Ref             *property     `json:"-"`
	RecursiveAnchor bool          `json:"-"`
	RecursiveRef    *property     `json:"-"`
	DynamicAnchor   string        `json:"-"`
	DynamicRef      *property     `json:"-"`
	Type            string        `json:"type,omitempty"`
	Types           []string      `json:"-"`              // allowed types.
	Constant        []interface{} `json:"-"`              // first element in slice is constant value. note: slice is used to capture nil constant.
	Enum            []interface{} `json:"enum,omitempty"` // allowed values.
	// enumError       string        // error message for enum fail. captured here to avoid constructing error message every time.
	Not   *property   `json:"-"`
	AllOf []*property `json:"-"`
	AnyOf []*property `json:"-"`
	OneOf []*property `json:"-"`
	If    *property   `json:"-"`
	Then  *property   `json:"-"` // nil, when If is nil.
	Else  *property   `json:"-"` // nil, when If is nil.

	// object validations
	MinProperties         int                          `json:"-"`                  // -1 if not specified.
	MaxProperties         int                          `json:"-"`                  // -1 if not specified.
	Required              []string                     `json:"required,omitempty"` // list of required properties.
	Properties            map[string]*property         `json:"properties,omitempty"`
	PropertyNames         *property                    `json:"-"`
	RegexProperties       bool                         `json:"-"` // property names must be valid regex. used only in draft4 as workaround in metaschema.
	PatternProperties     map[*regexp.Regexp]*property `json:"-"`
	AdditionalProperties  interface{}                  `json:"additionalProperties,omitempty"` // nil or bool or *property.
	Dependencies          map[string]interface{}       `json:"-"`                              // map value is *property or []string.
	DependentRequired     map[string][]string          `json:"-"`
	DependentSchemas      map[string]*property         `json:"-"`
	UnevaluatedProperties *property                    `json:"-"`

	// array validations
	MinItems    int  `json:"-"` // -1 if not specified.
	MaxItems    int  `json:"-"` // -1 if not specified.
	UniqueItems bool `json:"-"`
	// Items            interface{} `json:"items,omitempty"` // nil or *property or []*property
	Items            *property   `json:"items,omitempty"`
	AdditionalItems  interface{} `json:"-"` // nil or bool or *property.
	PrefixItems      []*property `json:"-"`
	Items2020        *property   `json:"-"` // items keyword reintroduced in draft 2020-12
	Contains         *property   `json:"-"`
	ContainsEval     bool        `json:"-"` // whether any item in an array that passes validation of the contains schema is considered "evaluated"
	MinContains      int         `json:"-"` // 1 if not specified
	MaxContains      int         `json:"-"` // -1 if not specified
	UnevaluatedItems *property   `json:"-"`

	// string validations
	MinLength        int            `json:"-"` // -1 if not specified.
	MaxLength        int            `json:"-"` // -1 if not specified.
	Pattern          *regexp.Regexp `json:"-"`
	ContentEncoding  string         `json:"-"`
	ContentMediaType string         `json:"-"`
	// mediaType        func([]byte) error `json:"-"`

	// number validators
	Minimum          *big.Rat `json:"-"`
	ExclusiveMinimum *big.Rat `json:"-"`
	Maximum          *big.Rat `json:"-"`
	ExclusiveMaximum *big.Rat `json:"-"`
	MultipleOf       *big.Rat `json:"-"`

	// user defined extensions
	Extensions map[string]ExtSchema `json:"-"`
}

func (p *property) read(t reflect.Type, opts tagOptions) {
	jsType, format, kind := getTypeFromMapping(t)
	if jsType != "" {
		p.Type = jsType
	}
	if format != "" {
		p.Format = format
	}

	switch kind {
	case reflect.Slice:
		p.readFromSlice(t)
	case reflect.Map:
		p.readFromMap(t)
	case reflect.Struct:
		p.readFromStruct(t)
	case reflect.Ptr:
		p.read(t.Elem(), opts)
	}
}

func (p *property) readFromSlice(t reflect.Type) {
	jsType, _, kind := getTypeFromMapping(t.Elem())
	if kind == reflect.Uint8 {
		p.Type = "string"
	} else if jsType != "" {
		oneItem := &property{}
		oneItem.read(t.Elem(), tagOptions(""))
		p.Items = oneItem
	}
}

func (p *property) readFromMap(t reflect.Type) {
	jsType, format, _ := getTypeFromMapping(t.Elem())

	if jsType != "" {
		p.Properties = make(map[string]*property, 0)
		p.Properties[".*"] = &property{Type: jsType, Format: format}
	} else {
		p.AdditionalProperties = true
	}
}

func (p *property) readFromStruct(t reflect.Type) {
	p.Type = "object"
	p.Properties = make(map[string]*property, 0)
	p.AdditionalProperties = false

	count := t.NumField()
	for i := 0; i < count; i++ {
		field := t.Field(i)

		tag := field.Tag.Get("json")
		name, opts := parseTag(tag)
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue
		}

		if field.Anonymous {
			embeddedProperty := &property{}
			embeddedProperty.read(field.Type, opts)

			for name, property := range embeddedProperty.Properties {
				p.Properties[name] = property
			}
			p.Required = append(p.Required, embeddedProperty.Required...)

			continue
		}

		p.Properties[name] = &property{}
		p.Properties[name].read(field.Type, opts)

		if !opts.Contains("omitempty") {
			p.Required = append(p.Required, name)
		}
	}
}

var formatMapping = map[string][]string{
	"time.Time": {"string", "date-time"},
}

var kindMapping = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
	reflect.Slice:   "array",
	reflect.Struct:  "object",
	reflect.Map:     "object",
}

func getTypeFormatByMapping(typeT string) (string, string) {
	if v, ok := formatMapping[typeT]; ok {
		return v[0], v[1]
	}
	return "", ""
}

func getTypeFromMapping(t reflect.Type) (string, string, reflect.Kind) {
	if v, ok := formatMapping[t.String()]; ok {
		return v[0], v[1], reflect.String
	}

	if v, ok := kindMapping[t.Kind()]; ok {
		return v, "", t.Kind()
	}

	return "", "", t.Kind()
}

type tagOptions string

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

//求交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 slice1-并集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
