package cp_nosql

import (
	"warehouse/v5-go-component/cp_log"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"fmt"
	"strings"
	"time"
	"warehouse/v5-go-component/cp_dc"
	"warehouse/v5-go-component/cp_error"
)

type Mongo struct {
	Client *mongo.Client
	Ctx	context.Context
}

type MongoCollection struct {
	c *mongo.Collection
	ctx context.Context
}

func (this *MongoCollection) InsertOne(opt interface{}) (*mongo.InsertOneResult, error) {
	return this.c.InsertOne(this.ctx, opt)
}

func (this *Mongo) NewCollection(db, collection string) *MongoCollection {
	m := &MongoCollection{}
	m.c = this.Client.Database(db).Collection(collection)
	m.ctx = this.Ctx
	return m
}

func NewMongoClient(config *cp_dc.DcNosqlConfig) (*Mongo, error) {
	var err error

	if len(config.Mongo.Hosts) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mc := &Mongo{
		Ctx: ctx,
	}

	connUri := fmt.Sprintf("mongodb://%s", strings.Join(config.Mongo.Hosts, ","))
	cp_log.Info(connUri)

	mc.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(connUri))
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	err = mc.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	return mc, nil
}




