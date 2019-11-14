package oplog_test

import (
	"log"

	"github.com/gnokoheat/oplog"
)

func ExampleTail() {
	var o = &oplog.Options{
		// (e.g. mongodb://username:password@127.0.0.1:27017,127.0.0.1:27018/local?replicaSet=rs01&authSource=admin)
		Addrs:      []string{"127.0.0.1:27017", "127.0.0.1:27018"}, // replicaset host and port
		Username:   "username",                                     // admin db username
		Password:   "password",                                     // admin db user password
		ReplicaSet: "rs01",                                         // replicaset name
		DB:         "myDB",                                         // tailing target db
		Collection: "myCollection",                                 // tailing target collection
		Events:     []string{"insert", "update", "delete"},         // tailing target method
	}

	l := make(chan *[]oplog.Log) // Oplog Channel
	e := make(chan error)        // Error Channel

	// Oplog tailing start !
	go o.Tail(l, e)

	for {
		select {
		case err := <-e:
			log.Println("[Error] ", err)
			return
		case op := <-l:
			// input oplog handling code
			log.Println("[Result] ", op)
			break
		}
	}
}
