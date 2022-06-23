package arex

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	dog "github.com/DataDog/zstd"
	"github.com/klauspost/compress/zstd"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type servletmocker struct {
	ID             string            `bson:"_id"`
	AppID          string            `bson:"appId,omitempty"`
	CreateTime     time.Time         `bson:"createTime,omitempty"`
	Env            int64             `bson:"env,omitempty"`
	Method         string            `bson:"method,omitempty"`
	Path           string            `bson:"path,omitempty"`
	Pattern        string            `bson:"pattern,omitempty"`
	RequestHeaders map[string]string `bson:"requestHeaders,omitempty"`
	Request        string            `bson:"request,omitempty"`
	Response       []byte            `bson:"response,omitempty"`
}

const servletmockerCollectionName = "ServletMocker"

func queryServletmocker(ctx context.Context, lastTime time.Time) []*servletmocker {
	db := ConnectOfMongoDB()
	scs := db.Collection(servletmockerCollectionName)

	filter := bson.M{}
	if !lastTime.IsZero() {
		filter = bson.M{
			"createTime": bson.M{
				"$gt": lastTime,
			},
		}
	}
	opts := options.Find().SetLimit(1000)

	cursor, err := scs.Find(ctx, filter, opts)
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil
	}
	sliceMockers := make([]*servletmocker, 0)
	for _, oneM := range results {
		var mocker servletmocker
		bsonBytes, _ := bson.Marshal(oneM)
		err := bson.Unmarshal(bsonBytes, &mocker)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		sliceMockers = append(sliceMockers, &mocker)
	}
	return sliceMockers
}

//
func unzipBase64andTstdString(in string) ([]byte, error) {
	unzipData, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// data, err := zstdDecompress(oneServlet.Response)
	// data, err := gozstd.Decompress(nil, oneServlet.Response)
	oriData, err := dog.Decompress(nil, unzipData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !json.Valid(oriData) {
		return nil, errors.New("json invalid")
	}
	return oriData, nil
}

func zstdDecompress(src []byte) ([]byte, error) {
	var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	return decoder.DecodeAll(src, nil)
}

func unzipBase64andGzipString(in string) ([]byte, error) {
	mybytes, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}

	res, err := gzipDecompress(mybytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func gzipDecompress(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func getAREXKey(serviceName, apiName string) string {
	var key strings.Builder
	key.WriteString(serviceName)
	key.WriteString("-")
	// key.WriteString(oneServlet.Path)
	key.WriteString(apiName)
	// uniKey := url.PathEscape(key.String())
	uniKey := key.String()
	return uniKey
}

// spiderAREXSchemaData(oneServlet.AppID, base64.URLEncoding.EncodeToString([]byte(oneServlet.Path)),string(oneServlet.Response)
func spiderAREXSchemaData(ctx context.Context, serviceName, apiName string, jsonStr string) {
	uniKey := getAREXKey(serviceName, apiName)

	curSchemaStore := querySchema(ctx, uniKey)
	if curSchemaStore == nil || curSchemaStore.Key == "" {
		curSchemaStore = &schemaStore{}
		curSchemaStore.Key = uniKey

		m, err := serviceGenerateSchema([]byte(jsonStr))
		if err != nil {
			fmt.Println(err)
			return
		}
		storeData, err := json.Marshal(m.Document)
		curSchemaStore.Schema = string(storeData)
	} else {
		b, err := serviceUpdateSchema(curSchemaStore.Schema, []byte(jsonStr))
		if err != nil {
			fmt.Println(err)
			return
		}
		storeData, err := json.Marshal(b)
		curSchemaStore.Schema = string(storeData)
	}
	saveSchema(ctx, *curSchemaStore)
}

// Opensource AREX DATA
func batchGenerateSchema(ctx context.Context, lastTime time.Time) {
	servletArrays := queryServletmocker(ctx, lastTime)
	if len(servletArrays) == 0 {
		return
	}

	for _, oneServlet := range servletArrays {
		if oneServlet.AppID == "" || oneServlet.Path == "" {
			continue
		}
		bytes, err := unzipBase64andTstdString(string(oneServlet.Response))
		if err != nil {
			fmt.Println(err)
			continue
		}
		spiderAREXSchemaData(ctx, oneServlet.AppID,
			base64.URLEncoding.EncodeToString([]byte(oneServlet.Path)),
			string(bytes))
	}
}

// Ctrip AREX Data
func batchGenerateByCtripAREX(ctx context.Context) {
	veriftyKeyAndJSON := func(unikey string, valTest interface{}) {
		if valTest == nil {
			return
		}

		testBytes, err := unzipBase64andGzipString(valTest.(string))
		if err != nil {
			log.Println(err)
			return
		}
		testJSON := make(map[string]interface{})
		err = json.Unmarshal(testBytes, &testJSON)
		if err != nil {
			log.Println(err)
			return
		}
		status, err := validateSchema(unikey, testBytes)
		if err != nil {
			log.Printf("validate fail %v %s", unikey, status)
		}
	}

	connectAREXFAT := func() *mongo.Collection {
		clientOptions := &options.ClientOptions{}
		clientOptions.ApplyURI(`mongodb://t_fltreplaydata:aUJxgKpsje8c_RWYdkSi@fltreplaydata01.mongo.db.fat.qa.nt.ctripcorp.com:55111/?replicaSet=fatfltpub01&authSource=fltreplaydatadb`)
		clientOptions.SetDirect(true)
		clientOptions.SetMaxPoolSize(100)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Panic(err)
		}
		db := client.Database("fltreplaydatadb")
		return db.Collection("ReplayCompareResultNew")
	}

	ds := connectAREXFAT()
	filter := bson.M{}
	// serverlist, err := ds.Distinct(ctx, "service", filter)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(serverlist...)
	opts := options.Find().SetLimit(1000)
	cursor, err := ds.Find(ctx, filter, opts)
	if err != nil {
		log.Panic(err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
	}

	for _, oneM := range results {
		if oneM["service"] == nil || oneM["resultname"] == nil {
			continue
		}

		val := oneM["basemsg"]
		baseBytes, err := unzipBase64andGzipString(val.(string))
		if err != nil {
			log.Println(err)
			continue
		}
		baseJSON := make(map[string]interface{})
		err = json.Unmarshal(baseBytes, &baseJSON)
		if err != nil {
			log.Println(err)
			continue
		}

		serviceName := oneM["service"].(string)
		apiName := oneM["resultname"].(string)
		spiderAREXSchemaData(ctx, serviceName, apiName, string(baseBytes))

		uniKey := getAREXKey(serviceName, apiName)
		valTest := oneM["testmsg"]
		veriftyKeyAndJSON(uniKey, valTest)
	}
}
