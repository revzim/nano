package drivers

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	AZMongoApp struct {
		client *mongo.Client
		// streamChannels map[string]chan *
	}
)

var (
	AZMongoErrUnInitialized = errors.New("mongo driver - not initialized")
)

func NewMongoApp(uri string) (*AZMongoApp, error) {

	// var err error
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	log.Println("mongo driver init")
	return &AZMongoApp{client: mongoClient}, nil

}

func (app *AZMongoApp) IsInit() bool {
	return app.client == nil
}

func (app *AZMongoApp) GetAllDatabases() ([]string, error) {
	if !app.IsInit() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		databases, err := app.client.ListDatabaseNames(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		return databases, nil
	}
	return nil, AZMongoErrUnInitialized
}

func (app *AZMongoApp) GetCollection(dbName, collectionName string) (*mongo.Collection, error) {
	var collection *mongo.Collection
	if !app.IsInit() {
		db := app.client.Database(dbName)
		collection = db.Collection(collectionName)

		return collection, nil
	}
	return collection, AZMongoErrUnInitialized
}

func (app *AZMongoApp) InsertIntoCollection(dbName, collectionName string, data map[string]interface{}, postOpDisconnect bool, opts *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if !app.IsInit() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if postOpDisconnect {
			defer app.client.Disconnect(ctx)
		}

		collection, err := app.GetCollection(dbName, collectionName)
		if err != nil {
			return nil, err
		}
		return collection.InsertOne(ctx, InitGenericMongoDoc(data), opts)
	}
	return nil, AZMongoErrUnInitialized
}
func (app *AZMongoApp) UpdateCollection(opKey, dbName, collectionName string, filter interface{}, data map[string]interface{}, postOpDisconnect bool, opts *options.UpdateOptions) (*mongo.UpdateResult, error) {
	if !app.IsInit() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if postOpDisconnect {
			defer app.client.Disconnect(ctx)
		}

		collection, err := app.GetCollection(dbName, collectionName)
		if err != nil {
			return nil, err
		}
		return collection.UpdateOne(ctx, filter, InitGenericSetMongoDoc(opKey, data), opts)
	}
	return nil, AZMongoErrUnInitialized
}

func (app *AZMongoApp) ConcurrentMongoStreamParse(dbName, collectionName string, closeOnFinish bool, streamCallback func(data bson.M)) {
	if !app.IsInit() {
		if closeOnFinish {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			defer app.client.Disconnect(ctx)
		}

		db := app.client.Database(dbName)
		collection := db.Collection(collectionName)

		var wg sync.WaitGroup
		currStream, err := collection.Watch(context.TODO(), mongo.Pipeline{})
		if err != nil {
			panic(err)
		}

		wg.Add(1)

		routineCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go app.parseStreamChanges(routineCtx, wg, currStream, streamCallback)

		wg.Wait()
	}

}

func (app *AZMongoApp) parseStreamChanges(ctx context.Context, wg sync.WaitGroup, stream *mongo.ChangeStream, f func(data bson.M)) {
	if !app.IsInit() {
		defer stream.Close(ctx)

		defer wg.Done()

		for stream.Next(ctx) {
			var data bson.M
			if err := stream.Decode(&data); err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("%v\n", data)
			f(data)
		}
	}
}

func InitGenericSetMongoDoc(key string, m map[string]interface{}) bson.M {
	genMongoDoc := bson.M{}
	for k, v := range m {
		genMongoDoc[k] = v
	}
	return bson.M{key: genMongoDoc}
}

func InitGenericMongoDoc(m map[string]interface{}) bson.M {
	genMongoDoc := bson.M{}
	for k, v := range m {
		genMongoDoc[k] = v
	}
	return genMongoDoc
}

// func parseStreamChanges(ctx context.Context, wg sync.WaitGroup, stream *mongo.ChangeStream) {
// 	checkMongoInit()
// 	defer stream.Close(ctx)

// 	defer wg.Done()

// 	for stream.Next(ctx) {
// 		var data bson.M
// 		if err := stream.Decode(&data); err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%v\n", data)
// 	}
// }

// func checkMongoInit() {
// 	if mongoMasterClient == nil {
// 		log.Fatal("mongo client not initialized! need to call `InitMongoClient()` | check out mongo.go file for info")
// 	}
// }
