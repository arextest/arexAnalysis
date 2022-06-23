package arex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
// @Summary json-schema interface
// @Description get-all-jsonschema
// @Tags tags
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func getSchemas(c *gin.Context) {
	limit := c.Query("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 10
	}
	res := querySchemas(context.Background(), int64(intLimit))
	c.IndentedJSON(http.StatusOK, res)
}

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

// input key, jsonData
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

func postComparing(c *gin.Context) {
	var compare comparing
	if err := c.BindJSON(&compare); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "struct failed"})
		return
	}

	res := serviceDiff2JSON(compare.ValueX, compare.ValueY)
	c.IndentedJSON(http.StatusAccepted, res.Diffs)
}
