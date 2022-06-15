package jsonschema

import (
	"fmt"
)

// SchemaModel add reader
type SchemaModel struct {
	// Writer      io.Writer
	WithExample bool
	Data        interface{}
	Name        string
	Format      bool
	Convert     bool
	Document    *SchemaDocument
}

// NewSchemaModel create schema model
func NewSchemaModel(data interface{}, name string) *SchemaModel {
	if name == "" {
		name = "Data"
	}
	name = replaceName(name)

	var version = 2020
	_, textSchema := readDraft(version)
	return &SchemaModel{
		// Writer:      os.Stdout,
		WithExample: true,
		Data:        data,
		Name:        name,
		Format:      true,
		Convert:     true,
		Document: &SchemaDocument{
			Schema: textSchema,
		},
	}
}

// SchemaFromBytes create instacne
func SchemaFromBytes(bytes []byte, name string) (*SchemaModel, error) {
	f, err := ParseJson(bytes)
	if err != nil {
		return nil, err
	}
	return NewSchemaModel(f, name), nil
}

// SchemaGetModel from bytes
func SchemaGetModel(url string) (*SchemaModel, error) {
	b, name, err := Get(url)
	if err != nil {
		return nil, err
	}
	return SchemaFromBytes(b, name)
}

// SchemaGenerateGo use json to generate schema
// root entry
func SchemaGenerateGo(f interface{}, name string) (*SchemaModel, error) {
	m := NewSchemaModel(f, name)
	m.GenerateSchema()
	return m, nil
}

// GenerateSchema todo
func (m *SchemaModel) GenerateSchema() {
	fu := func(ms map[string]interface{}, p *property) { m.parseMap(ms, p) }
	m.generate(fu)
}

// json root 是数组的话,就取第一个; 不是的话,就是全部
// 当跟节点是[],怎么处理
func (m *SchemaModel) generate(fu func(map[string]interface{}, *property)) {
	var ma map[string]interface{}
	switch v := m.Data.(type) {
	case []interface{}:
		ma = v[0].(map[string]interface{})
	default:
		ma = m.Data.(map[string]interface{})
	}
	fu(ma, &m.Document.property)
}

// 解析json的根,或者解析节点是map的节点
func (m *SchemaModel) parseMap(ms map[string]interface{}, p *property) {
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

// data是节点, keyName是节点名,p是传入的可以供当前使用的property
func (m *SchemaModel) parse(data interface{}, keyName string, p *property) {
	switch vv := data.(type) {
	case string:
		subType, converted := parseStringType(vv)
		if converted {
			sType, sFormat := getTypeFormatByMapping(subType)
			if sType == "" {
				m.printType(p, keyName, vv, "string", false)
			} else {
				m.printTypeFormat(p, keyName, vv, sType, sFormat)
			}
		} else {
			m.printType(p, keyName, vv, "string", false)
		}
	case bool:
		m.printType(p, keyName, vv, "boolean", false)
	case float64:
		//json parser always returns a float for number values, check if it is an int value
		m.printType(p, keyName, vv, "number", false)
	case int64:
		m.printType(p, keyName, vv, "integer", false)
	case []interface{}:
		p.Type = "array"
		// p.Items = &property{}

		if len(vv) > 0 {
			subProp := &property{}
			p.Items = subProp
			switch vvv := vv[0].(type) {
			case string:
				subType, converted := parseStringType(vvv)
				if converted {
					sType, sFormat := getTypeFormatByMapping(subType)
					if sType == "" {
						m.printType(subProp, "", vv, "string", false)
					} else {
						m.printTypeFormat(subProp, "", vv, sType, sFormat)
					}
				} else {
					m.printType(subProp, "", vv, "string", false)
				}
			case float64:
				//json parser always returns a float for number values, check if it is an int value
				m.printType(subProp, keyName, vv, "number", false)
			case int64:
				m.printType(subProp, keyName, vv, "integer", false)
			case bool:
				m.printType(subProp, "", vv[0], "boolean", false)
			case []interface{}:
				subProp.Type = "array"
				if len(data.([]interface{})) > 0 {
					lastP := &property{}
					subProp.Items = lastP
					m.parse(vvv, "", lastP)
				} else {
					m.printType(subProp, keyName, nil, "interface{}", false)
				}
			case map[string]interface{}:
				subProp.Type = "object"
				subProp.Properties = make(map[string]*property, 0)
				m.parseMap(vv[0].(map[string]interface{}), subProp)
			default:
				fmt.Printf("unknown type: %T", vvv)
				m.printType(subProp, keyName, nil, "interface{}", false)
			}
		} else {
			// empty []
		}
	case map[string]interface{}:
		p.Type = "object"
		m.parseMap(vv, p)
	case nil:
		p.Type = "null"
		m.printObject(p, keyName, "null")
	default:
		//fmt.Printf("unknown type: %T", vv)
		m.printType(p, keyName, nil, "interface{}", false)
	}
}

// value范例值;typeText类型;
func (m *SchemaModel) printType(p *property, keyName string, value interface{}, typeText string, converted bool) {
	// name := replaceName(keyName)
	if converted {
		keyName += ",string"
	}
	p.Type = typeText
}

func (m *SchemaModel) printTypeFormat(p *property, keyName string, value interface{}, typeText string, formatText string) {
	// name := replaceName(keyName)
	p.Type = typeText
	p.Format = formatText
}

func (m *SchemaModel) printObject(p *property, n string, t string) {
	// name := replaceName(n)
	p.Type = t
	// p.Format = append(p.Format, name)
}
