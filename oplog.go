package oplog

import (
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Options : MongoDB connection information for oplog tailing
type Options struct {
	Addrs      []string
	Username   string
	Password   string
	ReplicaSet string
	DB         string
	Collection string
	Events     []string
}

// Log : Oplog document
type Log struct {
	Timestamp    time.Time              `json:"wall" bson:"wall"`
	HistoryID    int64                  `json:"h" bson:"h"`
	MongoVersion int                    `json:"v" bson:"v"`
	Operation    string                 `json:"op" bson:"op"`
	Namespace    string                 `json:"ns" bson:"ns"`
	Doc          map[string]interface{} `json:"o" bson:"o"`
	Update       map[string]interface{} `json:"o2" bson:"o2"`
}

// MgoConn : MongoDB connect
func (o *Options) MgoConn(e chan error) (*mgo.Session, *mgo.Collection) {
	var mgoSess *mgo.Session
	var mgoColl *mgo.Collection

	m := &mgo.DialInfo{
		Addrs:          o.Addrs,
		Database:       "local",
		Direct:         false,
		FailFast:       false,
		Username:       o.Username,
		Password:       o.Password,
		ReplicaSetName: o.ReplicaSet,
		Source:         "admin",
		Mechanism:      "",
		Timeout:        time.Duration(0),
		PoolLimit:      0,
	}
	sess, err := mgo.DialWithInfo(m)
	if err != nil {
		e <- err
		return mgoSess, mgoColl
	}
	mgoSess = sess
	mgoColl = mgoSess.DB("local").C("oplog.rs")

	return mgoSess, mgoColl
}

// Tail : MongoDB oplog tailing start
func (o *Options) Tail(l chan *[]Log, e chan error) {
	bsonInsert := bson.M{}
	bsonUpdate := bson.M{}
	bsonDelete := bson.M{}
	for _, v := range o.Events {
		if v != "insert" && v != "update" && v != "delete" {
			e <- errors.New("Events type must be insert, update, delete")
		}
		if v == "insert" {
			bsonInsert = bson.M{"op": "i"}
		} else if v == "update" {
			bsonUpdate = bson.M{"op": "u"}
		} else if v == "delete" {
			bsonDelete = bson.M{"op": "d"}
		}
	}
	var events = bson.M{
		"$or": []bson.M{
			bsonInsert,
			bsonUpdate,
			bsonDelete,
		},
	}

	sTime := time.Now().UTC()
	mgoSess, mgoColl := o.MgoConn(e)
	defer mgoSess.Close()
	log.Println("[Oplog Tail Start] ", sTime)

	for {
		var fetchedLog = []Log{}
		err := mgoColl.Find(bson.M{
			"$and": []bson.M{
				bson.M{"ns": o.DB + "." + o.Collection},
				bson.M{"wall": bson.M{"$gt": sTime}},
				events,
			},
		}).Sort("wall:-1").All(&fetchedLog)
		if err != nil {
			e <- err
		}
		if len(fetchedLog) != 0 {
			sTime = fetchedLog[len(fetchedLog)-1].Timestamp
			l <- &fetchedLog
		}
	}
}
