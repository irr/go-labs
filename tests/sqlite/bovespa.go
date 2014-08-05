package main

import (
    "fmt"
    "os"
    "gosqlite.googlecode.com/hg/sqlite"
)

func main() {
    conn, err := sqlite.Open("/home/irocha/git/R-trader/data/symbols.db")

    if err == nil {
        fmt.Printf("sqlite ok: %+v,%+v\n", conn, err)

        s, err := conn.Prepare(fmt.Sprintf("select s, d, c from %v limit 5", "symbols"))

        if err != nil {
            os.Exit(1)
        }

        defer s.Finalize()

        err = s.Exec()

        if err != nil {
            os.Exit(1)
        }

        sym, dt, cl := "", "", 0.0

        for s.Next() {
            err := s.Scan(&sym, &dt, &cl)

            if err != nil {
                os.Exit(1)
            }

            fmt.Printf("%+v, %+v, %f\n", sym, dt, cl)
        }

        conn.Close()
    } else {
        fmt.Printf("sqlite error: %+v", err)
    }
}
