package Database

import (
	"TianTing/Logger"
	"TianTing/Settings"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)


const (
	authenticationStringTemplate = "%s:%s@"
	connectionStringTemplate     = "%s://%s%s/%s?"
)

type DocDB struct {
	IDatabase
	client *mongo.Client
	config *Settings.DocDbConf
}

func fetchCaFileFromPath(url string) error {
	getUrl, err := http.Get(url)
	if err != nil {
		return nil
	}
	out, err := os.Create("rds-ca.pem")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, getUrl.Body)
	return err
}

func ConnectWithDocDB(config *Settings.DocDbConf) (*DocDB, error) {
	// if ssl enabled, start download from path.
	if config.SSL {
		err := fetchCaFileFromPath(config.CaFilePath)
		if err != nil {
			return nil, err
		}
	}

	docDb := &DocDB{
		config: config,
	}

	authenticationURI := ""
	if config.Username != "" {
		authenticationURI = fmt.Sprintf(
			authenticationStringTemplate,
			config.Username,
			config.Password,
		)
	}

	connectionURI := fmt.Sprintf(
		connectionStringTemplate,
		config.Protocol,
		authenticationURI,
		config.Host,
		config.DefaultDb,
	)

	connectUri, _ := url.Parse(connectionURI)
	connectQuery, _ := url.ParseQuery(connectUri.RawQuery)

	if config.SSL {
		connectQuery.Add("ssl", "true")
		connectQuery.Add("sslcertificateauthorityfile", "rds-ca.pem")
	}

	if config.ReplicaSet != "" {
		connectQuery.Add("replicaSet", config.ReplicaSet)
	}

	if config.ReadPreference != "" {
		connectQuery.Add("readpreference", config.ReadPreference)
	}

	connectUri.RawQuery = connectQuery.Encode()

	Logger.SysLog.Infof("[DocumentDB] Try to connect document db `%s`", connectUri)

	clientOptions := options.Client().ApplyURI(connectUri.String())
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to Create New Client, %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConnectTimeoutMs)*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		Logger.SysLog.Errorf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, readpref.SecondaryPreferred())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to ping cluster: %s", err))
	}

	Logger.SysLog.Info("[DocumentDB] Connected to DocumentDB!")

	docDb.client = client

	return docDb, nil
}

func (db *DocDB) PopulateIndex(database, collection, key string, sort int32, unique bool) {
	c := db.client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(3 * time.Second)
	index := db.yieldIndexModel(
		[]string{key,}, []int32{sort,}, unique, nil,
	)
	_, err := c.Indexes().CreateOne(context.Background(), index, opts)
	if err != nil {
		Logger.SysLog.Errorf("[DocumentDb] Ensure Index Failed, %s", err)
	}
}

func (db *DocDB) PopulateTTLIndex(database, collection, key string, sort int32, unique bool, ttl int32) {
	c := db.client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(3 * time.Second)
	index := db.yieldIndexModel(
		[]string{key,}, []int32{sort,}, unique,
		options.Index().SetExpireAfterSeconds(ttl),
	)
	_, err := c.Indexes().CreateOne(context.Background(), index, opts)
	if err != nil {
		Logger.SysLog.Errorf("[DocumentDb] Ensure TTL Index Failed, %s", err)
	}
}

func (db *DocDB) PopulateMultiIndex(database, collection string, keys []string, sorts []int32, unique bool) {
	if len(keys) != len(sorts) {
		Logger.SysLog.Warnf("[DocumentDb] Ensure Indexes Failed, %s", "Please provide some item length of keys/sorts")
		return
	}
	c := db.client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(3 * time.Second)
	index := db.yieldIndexModel(keys, sorts, unique, nil)
	_, err := c.Indexes().CreateOne(context.Background(), index, opts)
	if err != nil {
		Logger.SysLog.Errorf("[DocumentDb] Ensure TTL Index Failed, %s", err)
	}
}

func (db *DocDB) yieldIndexModel(keys []string, sorts []int32, unique bool, indexOpt *options.IndexOptions) mongo.IndexModel {
	SetKeysDoc := bsonx.Doc{}
	for index, _ := range keys {
		key := keys[index]
		sort := sorts[index]
		SetKeysDoc = SetKeysDoc.Append(key, bsonx.Int32(sort))
	}
	if indexOpt == nil {
		indexOpt = options.Index()
	}
	indexOpt.SetUnique(unique)
	index := mongo.IndexModel{
		Keys:    SetKeysDoc,
		Options: indexOpt,
	}
	return index
}

func (db *DocDB) ListIndexes(database, collection string) {
	c := db.client.Database(database).Collection(collection)
	duration := 10 * time.Second
	batchSize := int32(10)
	cur, err := c.Indexes().List(context.Background(), &options.ListIndexesOptions{&batchSize, &duration})
	if err != nil {
		Logger.SysLog.Fatalf("[DocumentDB] Something went wrong listing %v", err)
	}
	for cur.Next(context.Background()) {
		index := bson.D{}
		_ = cur.Decode(&index)
	}
}

func (db *DocDB) GetClient() *mongo.Client {
	return db.client
}
