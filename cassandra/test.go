package main

import (
    "database/sql"
    "fmt"
    _ "github.com/tux21b/gocql"
)

func main() {
    db, err := sql.Open("gocql", "localhost:9042 keyspace=system")
    if err != nil {
        fmt.Println("Open error:", err)
    }

    rows, err := db.Query("SELECT keyspace_name FROM schema_keyspaces")
    if err != nil {
        fmt.Println("Query error:", err)
    }

    for rows.Next() {
        var keyspace string
        err = rows.Scan(&keyspace)
        if err != nil {
            fmt.Println("Scan error:", err)
        }
        fmt.Println(keyspace)
    }

    if err = rows.Err(); err != nil {
        fmt.Println("Iteration error:", err)
        return
    }
}
