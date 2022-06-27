package arex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arextest/arexAnalysis/jsonschema"
	"github.com/gin-gonic/gin"
)

// InstallHandler setup handle
func InstallHandler(engine *gin.Engine) {
	engine.GET("/schemas", middleware, getSchemas)
	engine.GET("/schema/:key", middleware, getSchemaByKey)
	engine.POST("/schema/:key", middleware, postSchema)
	engine.PUT("/schema/:key", middleware, putSchema)
	engine.PATCH("/schema/:key", middleware, patchSchema)
	engine.DELETE("/schema/:key", middleware, deleteSchema)

	engine.GET("/validation/:key", middleware, getValidation)
	engine.POST("/validation", middleware, postValidation)

	engine.POST("/comparing", middleware, postComparing)

	engine.GET("/testcases/postman/:appid", middleware, getTestCasesOfPostman)
	engine.GET("/testcases/golang/:appid", middleware, getTestCasesOfGolang)
}

func middleware(c *gin.Context) {
	// token := c.Query("token")
	// if token == "" {
	// 	log.Println("Url Param 'token' is missing")
	// 	c.JSON(200, gin.H{
	// 		"error": "please input correct token.",
	// 	})
	// 	c.Abort()
	// 	return
	// }

	// if c.Request.Method == "OPTIONS" {
	// 	c.AbortWithStatus(204)
	// 	return
	// }

	c.Next()
}

// AsyncHandle async config
func AsyncHandle(handle func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		cp := c.Copy()
		go func() {
			handle(cp)
		}()
		c.JSON(200, gin.H{
			"success": "true",
			"time":    time.Now(),
		})
	}
}

// getSchemas get all json-schema
// @Summary      Query all json-schema format json
// @Description  ?limit=10 limit the max range
// @Description  http Get /schemas
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        limit  path int  true  "query limit count"
// @Security     ApiKeyAuth
// @Success      200  {string} string  "[]json-schemas"
// @Router       /schemas [get]
func getSchemas(c *gin.Context) {
	limit := c.Query("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 10
	}
	res := querySchemas(context.Background(), int64(intLimit))
	c.IndentedJSON(http.StatusOK, res)
}

// getSchemaByKey get the special key json of json-schema
// @Summary      query json-schema by key
// @Description  Query one json-schema by key
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        key  path  string  true  "schema key name"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "{}"
// @Fail         400  {string}  string "---"
// @Router       /schema/{key} [get]
func getSchemaByKey(c *gin.Context) {
	key := c.Param("key")
	res := querySchema(context.Background(), key)
	if res == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "schema not found"})
		return
	}
	jsonData := make(map[string]interface{})
	json.Unmarshal([]byte(res.Schema), &jsonData)

	c.IndentedJSON(http.StatusOK, jsonData)
	return

}

// postSchema    postSchema json-schema
// @Summary      store json-schema to database by key
// @Description  post data to store. path /keyName. Body {}json-schema
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        key  path  string  true  "restapiApplication-L2FjdHVhdG9yL21hcHBpbmdz"
// @Param        body body  string  true  "{}"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "---"
// @Fail         400  {string}  string "---"
// @Router       /schema/{key} [post]
func postSchema(c *gin.Context) {
	saveSchemaByKey := func(key string, data []byte) {
		jsonData := make(map[string]interface{})
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "json struct failed"})
			return
		}
		schemajson, err := json.Marshal(jsonData)

		var ss schemaStore
		ss.Key = key
		ss.Schema = string(schemajson)
		saveSchema(context.Background(), ss)
	}
	key := c.Param("key")
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "put failed:" + err.Error()})
		return
	}
	saveSchemaByKey(key, data)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "success"})
}

// putSchema putSchema json-schema
// @Summary      input json and parse json to schema, then save the schema by key
// @Description  post /schema-key body contain origin json string {}
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        key  path  string  true  "schema key name"
// @Param        body body  string  true  "{json}"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "---"
// @Fail         400  {string}  string "---"
// @Router       /schema/{key} [put]
func putSchema(c *gin.Context) {
	compareSchemaToSave := func(key string, data []byte) *jsonschema.SchemaDocument {
		res, err := serviceGenerateSchema(data)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		var ss schemaStore
		ss.Key = key
		storeData, err := json.Marshal(res.Document)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		ss.Schema = string(storeData)
		saveSchema(context.Background(), ss)
		return res.Document
	}

	key := c.Param("key")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "put failed:" + err.Error()})
		return
	}
	doc := compareSchemaToSave(key, jsonData)
	c.IndentedJSON(http.StatusAccepted, doc)
}

// patchSchema   merge schema to existed schema
// @Summary      patchSchema to merge new json to schema and merge to existed json-schema
// @Description  post new json and parse it to merge existed json-schema
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        key  path  string  true  "schema key name"
// @Param        body body  string  true  "{}"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "---"
// @Fail         400  {string}  string "---"
// @Router       /schema/{key} [patch]
func patchSchema(c *gin.Context) {
	mergeSchemaByKey := func(key string, jsonData []byte) *jsonschema.SchemaDocument {
		oldSchema := querySchema(context.Background(), key)
		newschema, err := serviceUpdateSchema(oldSchema.Schema, jsonData)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		var ss schemaStore
		ss.Key = key
		storeData, err := json.Marshal(newschema)
		if err != nil {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "schema marshal error"})
		}
		ss.Schema = string(storeData)
		saveSchema(context.Background(), ss)
		return newschema
	}

	key := c.Param("key")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	newschema := mergeSchemaByKey(key, jsonData)
	if newschema == nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "put failed:" + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, newschema)
}

// deleteSchema  delete json-schema by key
// @Summary      delete json-schema by key
// @Description  send DELETE method http by jsonschema key
// @Tags         JSON-Schema
// @Accept       application/json
// @Produce      application/json
// @Param        key  path  string  true  "schema key name"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "---"
// @Fail         400  {string}  string "---"
// @Router       /schemas/{key} [delete]
func deleteSchema(c *gin.Context) {
	key := c.Param("key")
	res := delteSchemaData(context.Background(), key)
	if !res {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "delete failed"})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "delete Success"})
}

type validation struct {
	Key    string `json:"key"`
	Schema string `json:"schema"`
	Input  string `json:"input"`
	Result string `json:"result"`
}

// getValidation Validate json by json-schema that stored in database
// @Summary      Validate json by json-schema that stored in database
// @Description  get by keyname and body (Json format), then valid json by the keyname's json-schema
// @Tags         Validate by json-schema
// @Accept       application/json
// @Produce      application/json
// @Param        key   path  string  true  "schema key name"
// @Param        body  body  string  true  "{}"
// @Security     ApiKeyAuth
// @Success      200  {string}  string "---"
// @Fail         400  {string}  string "---"
// @Router       /validation/{key} [get]
func getValidation(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "schema key cannot be empty"})
		return
	}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "validation failed." + err.Error()})
		return
	}

	msg, err := validateSchema(key, jsonData)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "validation not ok:" + err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": msg})
	return
}

func validateSchema(key string, jsonData []byte) (string, error) {
	ss := querySchema(context.Background(), key)
	data := ss.Schema
	var sd jsonschema.SchemaDocument
	err := json.Unmarshal([]byte(data), &sd)
	if err != nil {
		return "", errors.New("json struct error")
	}

	return serviceValidateJSONBySchema(string(data), string(jsonData))
}

// postValidation  execute validate result
// @Summary      valid json by json-schema (input: validation struct)
// @Description  post struct that include schema's key and json that will be valid. return valid result
// @Description  if key is not exist, then it return nil
// @Tags         Validate by json-schema
// @Accept       application/json
// @Produce      application/json
// @Param        validation body  validation   true  "struct validation{}"
// @Security     ApiKeyAuth
// @Success      200   {string} string "{result}"
// @Failure      400   {string} string "{result}"
// @Router       /validation [post]
func postValidation(c *gin.Context) {
	var valid validation
	if err := c.BindJSON(&valid); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "schema not found"})
		return
	}
	if valid.Key == "" && valid.Schema == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "schema not found"})
		return
	}

	var schemaText string

	if valid.Key != "" {
		var sd *schemaStore
		sd = querySchema(context.TODO(), valid.Key)
		if sd == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "key not found"})
			return
		}
		// schema, err := json.Marshal(data.(*jsonschema.SchemaDocument))
		// if err != nil {
		// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "json struct error"})
		// 	return
		// }
		schemaText = string(sd.Schema)
	} else {
		schemaText = valid.Schema
	}

	msg, err := serviceValidateJSONBySchema(schemaText, valid.Input)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "validation failed." + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, gin.H{"message": msg})
	return
}

type comparing struct {
	ValueX  string `json:"vx"`
	ValueY  string `json:"vy"`
	Options string `json:"options"`
}

// postComparing  compare two json and get compared result
// @Summary      compare json
// @Description  post 2 json and return the difference
// @Tags         Comparing JSON
// @Accept       application/json
// @Produce      application/json
// @Param        body  body  comparing  true  "comparing struct"
// @Security     ApiKeyAuth
// @Success      201  {string}  string "[]object"
// @Failure      400  {string}  string "---"
// @Router       /comparing [post]
func postComparing(c *gin.Context) {
	var compare comparing
	if err := c.BindJSON(&compare); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "struct failed"})
		return
	}

	res := serviceDiff2JSON(compare.ValueX, compare.ValueY)
	c.IndentedJSON(http.StatusCreated, res.Diffs)
}

// getTestCasesOfPostman generate testcase that has postman format
// @Summary      Query testcases json of postman
// @Description  appid/?start=2022-2-22 limit the beggining
// @Description  http Get /schemas
// @Tags         Testcases
// @Accept       application/json
// @Produce      application/json
// @Param        start  path string  true  "start date"
// @Security     ApiKeyAuth
// @Success      200  {string} string  "json data"
// @Router       /testcases/postman/{appid} [get]
func getTestCasesOfPostman(c *gin.Context) {
	appid := c.Param("appid")
	start := c.Query("start")

	val := exportAREXToPostman(appid, start)
	c.IndentedJSON(http.StatusOK, val)
}

// getTestCasesOfGolang generate testcase that has golang format
// @Summary      Query testcases json of golang
// @Description  appid/?start=2022-2-22 limit the beggining
// @Description  http Get /schemas
// @Tags         Testcases
// @Accept       application/json
// @Produce      text/plain
// @Param        start  path string  true  "start date"
// @Security     ApiKeyAuth
// @Success      200  {string} string  "json data"
// @Router       /testcases/golang/{appid} [get]
func getTestCasesOfGolang(c *gin.Context) {
	appid := c.Param("appid")
	start := c.Query("start")

	val := getTestCases(appid, start)
	var caseText strings.Builder
	for _, oneCase := range val {
		caseText.WriteString(oneCase.ToCaseText())
		caseText.WriteString("\r\n\r\n")
	}
	c.String(http.StatusOK, caseText.String())
}
