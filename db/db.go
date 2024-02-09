package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Record struct {
	Id           string `bson:"_id"`
	RefreshToken string `bson:"refreshToken"`
}

type DAO struct {
	c *mongo.Collection
}

var DAOInstance DAO

func NewDAO(client mongo.Client) {
	if DAOInstance.c == nil {
		DAOInstance.c = client.Database("core").Collection("things")
	}
}

func GetDAO() *DAO {
	return &DAOInstance
}

func (dao *DAO) Insert(record Record) error {
	_, err := DAOInstance.c.InsertOne(context.Background(), record)
	return err
}

func (dao *DAO) FindById(id string) (Record, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var record Record
	err := DAOInstance.c.FindOne(context.Background(), filter).Decode(&record)
	switch err {
	case nil:
		return record, nil
	default:
		return Record{"", ""}, err
	}
}

func (dao *DAO) Update(id, refreshToken string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := DAOInstance.c.ReplaceOne(context.Background(), filter, Record{
		Id:           id,
		RefreshToken: refreshToken,
	})
	return err
}
