package drivers

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_URI = ""
)

func TestMongoDBs(t *testing.T) {
	app, err := NewMongoApp(MONGO_URI) // NewMongoApp(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}

	dbs, err := app.GetAllDatabases()

	if err != nil {
		panic(err)
	}
	for i := range dbs {
		log.Println(dbs[i])
	}

}

func TestMongoUpdates(t *testing.T) {
	app, err := NewMongoApp(MONGO_URI) // NewMongoApp(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	b := new(bool)
	*b = false
	updateOpts := &options.UpdateOptions{Upsert: b}
	dbName := ""
	collectionName := ""
	app.UpdateCollection("$set", dbName, collectionName, map[string]interface{}{
		"owner_id": "TEST",
	}, map[string]interface{}{
		"owner_id": "TEST",
		"id":       uuid.New().String(),
		"data": map[string]interface{}{
			"key": "abc",
		},
	}, false, updateOpts)

	*b = true

	app.UpdateCollection("$set", dbName, collectionName, map[string]interface{}{
		"owner_id": "TEST",
	}, map[string]interface{}{
		"owner_id": "TEST",
		"id":       uuid.New().String(),
		"data": map[string]interface{}{
			"key": "xyz",
		},
	}, true, updateOpts)
}
