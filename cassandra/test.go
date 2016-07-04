package main

import (
    "fmt"
    "log"
    "github.com/gocql/gocql"
)

func main() {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "system_schema"
    cluster.Consistency = gocql.Quorum
    cluster.ProtoVersion = 4
    session, _ := cluster.CreateSession()
    defer session.Close()

    iter := session.Query("SELECT keyspace_name FROM keyspaces").Iter()

    var keyspace string

    for iter.Scan(&keyspace) {
        fmt.Println(keyspace)
    }

    if err := iter.Close(); err != nil {
        log.Fatal(err)
    }
}
