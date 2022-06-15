package arex

import (
	"encoding/json"
	"errors"
	"fmt"
	"local/arex-reporter/comparer"
	"local/arex-reporter/jsonschema"
)

// serviceGenerateSchema input: json output: json-schema
func serviceGenerateSchema(data interface{}) (*jsonschema.SchemaModel, error) {
	f, err := jsonschema.ParseJson(data.([]byte))
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	schemaDoc, _ := jsonschema.SchemaGenerateGo(f, "")
	// schemaText, err := schemaDoc.Document.String()
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// 	return "", err
	// }
	return schemaDoc, nil
}

// serviceValidate2JSONBySchema compare 2 json, wether are those jsons same.
func serviceValidate2JSONBySchema(dataX string, dataY string) (bool, error) {
	return false, nil
}

// serviceValidate2Schema compare 2 json-schema, wether those are same.
func serviceValidate2Schema(schemaX string, schemaY string) (bool, error) {
	return false, nil
}

// serviceValidateJSONBySchema valid json by schema
func serviceValidateJSONBySchema(dataSchema string, data string) (bool, error) {
	schema, err := jsonschema.CompileString("jason-schema", dataSchema)
	if err != nil {
		return false, err
	}
	var someInterface interface{}
	json.Unmarshal([]byte(data), &someInterface)

	err = schema.Validate(someInterface)
	if err == nil {
		return true, nil
	}
	return false, err
}

// serviceUpdateSchema update schema by new json return new schema
func serviceUpdateSchema(jsonSchema string, args ...string) (string, error) {
	if jsonSchema == "" {
		return "", errors.New("empty schema")
	}
	var schema jsonschema.SchemaDocument
	json.Unmarshal([]byte(jsonSchema), &schema)
	schemaChan := make(chan jsonschema.SchemaDocument, len(args))
	for _, arg := range args {
		go func(jsonData *string) {
			res, err := jsonschema.SchemaGenerateGo(*jsonData, "")
			if err != nil {
				fmt.Printf("error %v\n", err)
				return
			}
			schemaChan <- *res.Document
		}(&arg)
	}

	for oneschema := range schemaChan {
		err := schema.MergeSchemaDocument(&oneschema)
		if err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
	close(schemaChan)

	return schema.String()
}

// serviceDiff2JSON compare 2 json and return json result
func serviceDiff2JSON(dataX, dataY string) *comparer.DiffReporter {
	dx := make(map[string]interface{})
	json.Unmarshal([]byte(dataX), &dx)
	dy := make(map[string]interface{})
	json.Unmarshal([]byte(dataY), &dy)

	jsonResult := comparer.GoCmpDiff(dx, dy)
	return jsonResult
}
