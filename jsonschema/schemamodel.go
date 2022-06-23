package jsonschema

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"strconv"
	"time"
)

// SchemaDataModel add reader
type SchemaDataModel struct {
	// Writer      io.Writer
	WithExample bool
	Data        []byte
	Name        string
	Format      bool
	Convert     bool
	Document    *SchemaDocument
}

// NewSchemaDataModel create schema model
func NewSchemaDataModel(data []byte, modelName string) *SchemaDataModel {
	if modelName == "" {
		modelName = "Data"
	}
	modelName = replaceName(modelName)

	var version = 2020
	_, textSchema := readDraft(version)
	return &SchemaDataModel{
		// Writer:      os.Stdout,
		WithExample: true,
		Data:        data,
		Name:        modelName,
		Format:      true,
		Convert:     true,
		Document: &SchemaDocument{
			Schema: textSchema,
		},
	}
}

// SchemaGetModel from bytes
func SchemaGetModel(url string) (*SchemaDataModel, error) {
	b, name, err := Get(url)
	if err != nil {
		return nil, err
	}
	return NewSchemaDataModel(b, name), nil
}

// GenerateSchemaDataModel use json to generate schema
// root entry
func GenerateSchemaDataModel(f []byte, name string) (*SchemaDataModel, error) {
	m := NewSchemaDataModel(f, name)
	err := m.generate()
	return m, err
}

func (m *SchemaDataModel) generate() error {
	if m.Data == nil {
		return errors.New("data is empty")
	}

	jsonData, err := ParseJson(m.Data)
	if err != nil {
		return err
	}

	return m.parse(jsonData, "", &m.Document.property)
}

// pase object.
func (m *SchemaDataModel) parse(data interface{}, keyName string, p *property) error {
	switch vv := data.(type) {
	case string:
		m.parseString(vv, keyName, p)
	case bool:
		m.parseBool(vv, keyName, p)
	case float64:
		m.parseNumber(vv, keyName, p)
	case int64:
		m.parseInteger(vv, keyName, p)
	case []interface{}:
		m.parseArray(vv, keyName, p)
	case map[string]interface{}:
		m.parseMap(vv, keyName, p)
	case nil:
		m.parseNull(keyName, p)
	default:
		return fmt.Errorf("unknown type: %T", vv)
	}
	return nil
}

// Parse map struct such as json root or json node
func (m *SchemaDataModel) parseMap(ms map[string]interface{}, keyName string, p *property) {
	p.Type = "object"
	p.Properties = make(map[string]*property, 0)

	keys := getSortedKeys(ms)
	for _, k := range keys {
		tempP := &property{}
		p.Properties[k] = tempP
		m.parse(ms[k], k, tempP)
		p.Required = append(p.Required, k)
	}
}

func (m *SchemaDataModel) parseInteger(vv int64, keyName string, p *property) {
	//json parser always returns a float for number values, check if it is an int value
	p.Type = "integer"
	p.Maximum = float64(vv)
	p.Minimum = float64(vv)
}

func (m *SchemaDataModel) parseNumber(vv float64, keyName string, p *property) {
	//json parser always returns a float for number values, check if it is an int value
	p.Type = "number"
	p.Maximum = vv
	p.Minimum = vv
}

func (m *SchemaDataModel) parseBool(vv bool, keyName string, p *property) {
	p.Type = "boolean"
}

func (m *SchemaDataModel) parseNull(keyName string, p *property) {
	p.Type = "null"
}

func (m *SchemaDataModel) parseString(vv string, keyName string, p *property) {
	fillingGeneralString := func(cv string, types string, format string, p *property) {
		p.Type = types
		p.Format = format
		p.MaxLength = len(vv)
		p.MinLength = len(vv)

		if p.MaxLength < constNotEnumMaxLength {
			p.Examples = append(p.Examples, vv)
		}
	}

	parseStringInlineFormat := func(value string) (string, string) {
		if ok := isDate(value); ok {
			return getTypeFormatByMapping("date")
		} else if ok := isDateTime(value); ok {
			return getTypeFormatByMapping("time.Time")
		} else if _, err := time.Parse(time.RFC3339, value); err == nil {
			return getTypeFormatByMapping("time.Time")
		} else if ip := net.ParseIP(value); ip != nil {
			if ip.To4() != nil {
				return getTypeFormatByMapping("ipv4")
			}
			return getTypeFormatByMapping("ipv6")
		} else if _, err := mail.ParseAddress(value); err == nil {
			return getTypeFormatByMapping("email")
		} else if _, err := url.Parse(value); err == nil {
			return getTypeFormatByMapping("url")
		} else if ok := isURI(value); ok {
			return getTypeFormatByMapping("uri")
		} else if ok := isURIReference(value); ok {
			return getTypeFormatByMapping("uri-reference")
		} else if ok := isJSONPointer(value); ok {
			return getTypeFormatByMapping("json-pointer")
		} else if ok := isRelativeJSONPointer(value); ok {
			return getTypeFormatByMapping("relative-json-pointer")
		} else if ok := isRegex(value); ok {
			return getTypeFormatByMapping("regex")
		} else if _, err := strconv.ParseBool(value); err == nil {
			return "string", ""
			// TODO return "string", "bool"=> enum("true","false")
		} else {
			return "string", ""
		}
	}

	subType, format := parseStringInlineFormat(vv)
	fillingGeneralString(vv, subType, format, p)
}

func (m *SchemaDataModel) parseArray(vv []interface{}, keyName string, p *property) {
	p.Type = "array"
	p.MinItems = len(vv)
	p.MaxItems = len(vv)
	if len(vv) > 0 {
		subProp := &property{}
		p.Items = subProp
		m.parse(vv[0], keyName, subProp)
	}
}

const constNotEnumMaxLength = 20

var formatMapping = map[string][]string{
	"time":                  {"string", "time"},
	"date":                  {"string", "date"},
	"time.Time":             {"string", "date-time"},
	"ipv4":                  {"string", "ipv4"},
	"ipv6":                  {"string", "ipv6"},
	"email":                 {"string", "email"},
	"idn-email":             {"string", "idn-email"},
	"hostname":              {"string", "hostname"},
	"idn-hostname":          {"string", "idn-hostname"},
	"uri":                   {"string", "uri"},
	"uri-reference":         {"string", "uri-reference"},
	"iri":                   {"string", "iri"},
	"iri-reference":         {"string", "iri-reference"},
	"uri-template":          {"string", "uri-template"},
	"json-pointer":          {"string", "json-pointer"},
	"relative-json-pointer": {"string", "relative-json-pointer"},
	"regex":                 {"string", "regex"},
}

func getTypeFormatByMapping(typeT string) (string, string) {
	if v, ok := formatMapping[typeT]; ok {
		return v[0], v[1]
	}
	return "string", ""
}
