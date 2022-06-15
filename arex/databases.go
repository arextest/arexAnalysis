package arex

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type schemaStore struct {
	ID         bson.ObjectId `json:"_id,omitempty"`
	Key        string        `json:"key"`
	Schema     string        `json:"schema"`
	LastUpdate time.Time     `json:"lastupdate"`
}

var mongoDatabase *mongo.Database

const schemaCollectionName string = "schemas"

// ConnectOfMongoDB get default mongodb
func ConnectOfMongoDB() *mongo.Database {
	if mongoDatabase != nil {
		return mongoDatabase
	}
	settings := map[string]string{
		"Database": `arex_storage_db`,  // Database name
		"Host":     `10.5.153.1:27017`, // Server IP or name
		"User":     `arex`,
		"Password": `iLoveArex`,
		"URL":      `mongodb://arex:iLoveArex@10.5.153.1:27017/arex_storage_db`,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := &options.ClientOptions{}
	clientOptions.ApplyURI(settings["URL"])
	clientOptions.SetDirect(true)
	clientOptions.SetMaxPoolSize(100)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}
	return client.Database(settings["Database"])
}

func saveSchema(ctx context.Context, item schemaStore) {
	db := ConnectOfMongoDB()
	opts := options.Update().SetUpsert(true)
	scs := db.Collection(schemaCollectionName)

	filter := bson.M{"key": item.Key}
	update := bson.M{"$set": bson.M{"schema": item.Schema, "lastupdate": time.Now()}}

	// result, err := scs.InsertOne(ctx, item)
	result, err := scs.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Printf("save new document failed %s\n", err)
		return
	}
	if result.MatchedCount > 1 {
		fmt.Println("error result.matchcount == 1")
	}
}

func querySchema(ctx context.Context, key string) *schemaStore {
	db := ConnectOfMongoDB()
	scs := db.Collection(schemaCollectionName)
	var schema schemaStore

	filter := bson.M{"key": key}
	var result bson.M
	scs.FindOne(ctx, filter).Decode(&result)
	// cursor, err := scs.Find(ctx, filter, nil)
	// if err != nil {
	// 	return nil
	// }
	// defer cursor.Close(ctx)

	// var results []bson.M
	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	return nil
	// }
	bsonBytes, _ := bson.Marshal(result)
	err := bson.Unmarshal(bsonBytes, &schema)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	return &schema
}

func querySchemas(ctx context.Context, limit int64) *map[string]*schemaStore {
	db := ConnectOfMongoDB()
	scs := db.Collection(schemaCollectionName)

	filter := bson.M{}
	opts := options.Find().SetLimit(limit)
	cursor, err := scs.Find(ctx, filter, opts)
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil
	}
	schemaMap := make(map[string]*schemaStore)
	for _, oneM := range results {
		var schema schemaStore
		bsonBytes, _ := bson.Marshal(oneM)
		err := bson.Unmarshal(bsonBytes, &schema)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		schemaMap[schema.Key] = &schema
	}
	return &schemaMap
}

func delteSchemaData(ctx context.Context, key string) bool {
	db := ConnectOfMongoDB()
	scs := db.Collection(schemaCollectionName)

	filter := bson.M{"key": key}
	res, err := scs.DeleteOne(ctx, filter)
	if err != nil {
		return false
	}
	fmt.Printf("delete count %d\n", res.DeletedCount)
	return true
}
