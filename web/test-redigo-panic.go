// curl -s "http://localhost:12345/test?u=root"
// go run test-redigo.go -l (enable syslog)
// tail -f /var/log/message

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"time"
)

var l *bool
var db *sql.DB
var rd *redis.Pool

var logger *syslog.Writer
var err error

func logPanics(function func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v", request.RemoteAddr, x)
				http.Error(writer, x.(string), http.StatusInternalServerError)
			}
		}()
		function(writer, request)
	}
}

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

func getValue(req *http.Request, n string, d string) string {
	v := req.FormValue(n)
	return (T((v == ""), d, v)).(string)
}

func checkError(err error, w http.ResponseWriter, msg string, status int) bool {
	if err != nil {
		log.Printf("%s: %s\n", msg, err.Error())
		http.Error(w, err.Error(), status)
	}
	return (err != nil)
}

func TestServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query, err := db.Query("SELECT Host FROM user WHERE User = ?", getValue(req, "u", "root"))
	if checkError(err, w, "db.Query", http.StatusInternalServerError) {
		return
	}
	defer query.Close()

	columns, err := query.Columns()
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(values))
	for i := range values {
		scans[i] = &values[i]
	}

	mr := make(map[string][]string)
	for query.Next() {
		err = query.Scan(scans...)
		if checkError(err, w, "query.Scan", http.StatusInternalServerError) {
			return
		}
		for _, col := range values {
			mr["mysql"] = append(mr["mysql"], T((col == nil), "NULL", string(col)).(string))
		}
	}

	conn := rd.Get()
	if checkError(err, w, "rd.Get", http.StatusInternalServerError) {
		return
	}
	defer conn.Close()

	/*
	   === INFO ===
	   info, err := redis.String(conn.Do("INFO"))
	   if (checkError(err, w, "conn.Do", http.StatusInternalServerError)) { return }

	   for _, s := range strings.Split(info, "\r\n") {
	       if s != "" && !strings.HasPrefix(s, "#") {
	           mr["redis"] = append(mr["redis"], s)
	       }
	   }
	*/

	/*
	   === MULTI/EXEC ===
	   conn.Send("MULTI")
	   conn.Send("ZADD", "family", 1, "alessandra")
	   conn.Send("ZADD", "family", 10, "babi")
	   conn.Send("ZADD", "family", 100, "lara")
	   conn.Send("ZADD", "family", 1000, "luma")
	   conn.Send("ZRANGEBYSCORE", "family", 1, 1000)
	   multi, err := redis.Values(conn.Do("EXEC"))
	   if (checkError(err, w, "redis.Values", http.StatusInternalServerError)) { return }

	   if len(multi) > 0 {
	       var i1, i2, i3, i4 int
	       var ranges []string
	       multi, err = redis.Scan(multi, &i1, &i2, &i3, &i4, &ranges)
	       mr["redis"] = ranges
	   }
	*/

	conn.Send("ZADD", "family", 1, "alessandra")
	conn.Send("ZADD", "family", 10, "babi")
	conn.Send("ZADD", "family", 100, "lara")
	conn.Send("ZADD", "family", 1000, "luma")
	conn.Send("ZRANGEBYSCORE", "family", 1, 1000)
	conn.Flush()
	n := 5
	for n > 0 {
		if n == 1 {
			v, err := redis.Strings(conn.Receive())
			if checkError(err, w, "conn.Receive(Strings)", http.StatusInternalServerError) {
				return
			}
			mr["redis"] = v
		} else {
			_, err := conn.Receive()
			if checkError(err, w, "conn.Receive(int64)", http.StatusInternalServerError) {
				return
			}
		}
		n--
	}

	j, err := json.Marshal(mr)
	if checkError(err, w, "json.Marshal", http.StatusInternalServerError) {
		return
	}
	fmt.Fprintf(w, "%s\n", j)
	if *l {
		logger.Info(string(j))
	}
}

func main() {
	l = flag.Bool("l", false, "syslog enabled/disabled")
	h := flag.Bool("h", false, "help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: test-redigo [-l][-h]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if *l {
		logger, err = syslog.New(syslog.LOG_INFO, "[test-redigo]")
		if err != nil {
			log.Fatal("syslog: ", err)
		}
	}
	db, err = sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/mysql?autocommit=true")
	if err != nil && db.Ping() != nil {
		log.Fatal("sql.Open/db.Ping: ", err)
	}
	defer db.Close()
	rd = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("GET", "FOO")
			return err
		},
	}
	http.HandleFunc("/test", logPanics(TestServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
