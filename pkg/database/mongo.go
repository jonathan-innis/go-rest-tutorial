package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB(collection *mongo.Collection) *MongoDB {
	return &MongoDB{collection: collection}
}

func (db *MongoDB) Create(ctx context.Context, item interface{}) error {
	_, err := db.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) List(ctx context.Context, model interface{}) ([]interface{}, error) {
	cur, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []interface{}
	for cur.Next(ctx) {
		err := cur.Decode(model)
		if err != nil {
			return nil, err
		}
		res = append(res, model)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
