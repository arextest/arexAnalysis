package jsonschema

import (
	"errors"
	"fmt"
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
	m.fillProperType(p, "integer")
}

func (m *SchemaDataModel) parseNumber(vv float64, keyName string, p *property) {
	//json parser always returns a float for number values, check if it is an int value
	m.fillProperType(p, "number")
}

func (m *SchemaDataModel) parseBool(vv bool, keyName string, p *property) {
	m.fillProperType(p, "boolean")
}

func (m *SchemaDataModel) parseNull(keyName string, p *property) {
	p.Type = "null"
	m.fillObject(p, keyName, "null")
}

func (m *SchemaDataModel) parseString(vv string, keyName string, p *property) {
	subType, needConvert := parseStringFormat(vv)

	if needConvert {
		sType, sFormat := getTypeFormatByMapping(subType)
		if sType != "" {
			m.fillProperTypeFormat(p, sType, sFormat)
			return
		}
	}
	m.fillProperType(p, "string")
}

func (m *SchemaDataModel) parseArray(vv []interface{}, keyName string, p *property) {
	p.Type = "array"
	if len(vv) > 0 {
		subProp := &property{}
		p.Items = subProp
		m.parse(vv[0], keyName, subProp)
	}
}

// value范例值;typeText类型;
func (m *SchemaDataModel) fillProperType(p *property, typeText string) {
	p.Type = typeText
}

func (m *SchemaDataModel) fillProperTypeFormat(p *property, typeText string, formatText string) {
	// name := replaceName(keyName)
	p.Type = typeText
	p.Format = formatText
}

func (m *SchemaDataModel) fillObject(p *property, n string, t string) {
	// name := replaceName(n)
	p.Type = t
	// p.Format = append(p.Format, name)
}
