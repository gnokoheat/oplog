# oplog
MongoDB oplog tailing by Golang


## Install

``` go
go get -u github.com/gnokoheat/oplog
```

## Usage example

- main.go

``` go
package main

import (
	"log"
	"github.com/gnokoheat/oplog"
)

func main() {
	// mongodb://username:password@127.0.0.1:27017,127.0.0.1:27018/local?replicaSet=rs01&authSource=admin
	var o = &oplog.Options{
		Addrs:      []string{"127.0.0.1:27017", "127.0.0.1:27018"},
		Username:   "username",
		Password:   "password",
		ReplicaSet: "rs01",
		DB:         "myDB",
		Collection: "myCollection",
		Events:     []string{"insert", "update", "delete"},
	}

	l := make(chan *[]oplog.Log)
	e := make(chan error)
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
```

- Result
```
2019/11/08 16:13:57 [Oplog Tail Start]  2019-11-08 07:13:57.485633 +0000 UTC
2019/11/08 16:14:04 [Result]  &[{2019-11-08 16:14:02.744 +0900 ... ]
2019/11/08 16:14:08 [Result]  &[{2019-11-08 16:14:08.554 +0900 ... ]
2019/11/08 16:16:57 [Result]  &[{2019-11-08 16:16:57.364 +0900 ... ]
```
