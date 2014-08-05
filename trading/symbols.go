package main

import (
    "os"
    "mysql"
    "fmt"
)

const (
    MYSQL_HOST = "127.0.0.1:3306"
    MYSQL_USER = "root"
    MYSQL_PASS = "mysql"
    MYSQL_DB   = "bovespa"
)

type OHLC struct {
    Date   string
    Open   float64
    High   float64
    Low    float64
    Close  float64
    Volume float64
}

func GetSymbols(a ...interface{}) ([]OHLC, os.Error) {
    db, err := mysql.DialTCP(MYSQL_HOST, MYSQL_USER, MYSQL_PASS, MYSQL_DB)
    if err != nil {
        return nil, err
    }

    var stmt *mysql.Statement
    var ohlcs []OHLC
    var sql string

    defer func() {
        if db != nil && stmt != nil {
            db.FreeResult()
        }
    }()

    switch {
    case len(a) == 3:
        sql = "SELECT D, O, H, L, C, V FROM symbols WHERE S = ? AND D BETWEEN ? AND ? ORDER BY D"
        stmt, err = db.Prepare(sql)
        if err != nil {
            return nil, err
        }
        err = stmt.BindParams(a[0], a[1], a[2])
        if err != nil {
            return nil, err
        }
    case len(a) == 2:
        sql = "SELECT D, O, H, L, C, V FROM symbols WHERE S = ? AND D >= ? ORDER BY D"
        stmt, err = db.Prepare(sql)
        if err != nil {
            return nil, err
        }
        err = stmt.BindParams(a[0], a[1])
        if err != nil {
            return nil, err
        }
    default:
        sql = "SELECT D, O, H, L, C, V FROM symbols WHERE S = ? ORDER BY D"
        stmt, err = db.Prepare(sql)
        if err != nil {
            return nil, err
        }
        err = stmt.BindParams(a[0])
        if err != nil {
            return nil, err
        }
    }

    fmt.Printf("executing \"%s\" with arguments %v...\n", sql, a)

    err = stmt.Execute()
    if err != nil {
        return nil, err
    }

    err = stmt.StoreResult()
    if err != nil {
        return nil, err
    }

    c := stmt.RowCount()

    if c <= 0 {
        return make([]OHLC, 0), nil
    }

    var ohlc OHLC
    stmt.BindResult(&ohlc.Date, &ohlc.Open, &ohlc.High, &ohlc.Low, &ohlc.Close, &ohlc.Volume)

    fmt.Printf("found %d records...\n", c)

    ohlcs = make([]OHLC, c, c)

    for i := 0; ; i++ {
        eof, err := stmt.Fetch()
        if err != nil {
            return nil, err
        }

        if eof {
            break
        }

        ohlcs[i] = ohlc
    }

    fmt.Printf("%d records added\n", len(ohlcs))

    return ohlcs, nil
}
