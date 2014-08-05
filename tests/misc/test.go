package main

import (
    "fmt"
    "sync"
    "io/ioutil"
)

func dump1(s string, a []int) {
    for k, v := range a {
        fmt.Printf("%s[%d] = %d\n", s, k, v)
    }
}

func dump1ref(s string, a []int) {
    for i := 0; i < len(a); i++ {
        a[i] *= 10
        fmt.Printf("%s[%d] = %d\n", s, i, a[i])
    }
}

func dump1cpy(s string, a [5]int) {
    for i := 0; i < len(a); i++ {
        a[i] *= 10
        fmt.Printf("%s[%d] = %d\n", s, i, a[i])
    }
}

func dump2(s string, m map[string]int) {
    for k, v := range m {
        fmt.Printf("%s[%s] = %d\n", s, k, v)
    }
}

func test1(p ...interface{}) {
    for i := range p {
        fmt.Printf("%d: %v\n", i, p[i])
    }
}

type sortInterface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

func sort(data sortInterface) {
    for i := 1; i < data.Len(); i++ {
        for j := i; j > 0 && data.Less(j, j-1); j-- {
            data.Swap(j, j-1)
        }
    }
}

func files(d string) string {
    fs, _ := ioutil.ReadDir(d)
    s := ""
    for _, v := range fs {
        base := d + v.Name + "/"
        s += base + "\n" + files(base)
    }
    return s
}

type intArray []int

func (p intArray) Len() int           { return len(p) }
func (p intArray) Less(i, j int) bool { return p[i] < p[j] }
func (p intArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type day struct {
    num       int
    shortName string
    longName  string
}

type dayArray struct {
    data []*day
}

func (p *dayArray) Len() int           { return len(p.data) }
func (p *dayArray) Less(i, j int) bool { return p.data[i].num < p.data[j].num }
func (p *dayArray) Swap(i, j int)      { p.data[i], p.data[j] = p.data[j], p.data[i] }

func ping() chan string {
    ch := make(chan string)
    go func() {
        ch <- "ping"
        fmt.Printf("ping sent ok\n")
    }()
    return ch
}

func pong(in chan string) chan string {
    ch := make(chan string)
    go func() {
        s := <-in
        fmt.Printf("pong received: %v\n", s)
        ch <- "pong"
    }()
    return ch
}

type request struct {
    a, b   int
    replyc chan int
}

type binOp func(a, b int) int

func run(op binOp, req *request) {
    reply := op(req.a, req.b)
    fmt.Printf("adding: %d + %d\n", req.a, req.b)
    req.replyc <- reply
}

func server(op binOp, service chan *request, quit chan bool) {
    for {
        select {
        case req := <-service:
            go run(op, req)
        case <-quit:
            return
        }
    }
}

func startServer(op binOp) (service chan *request, quit chan bool) {
    service = make(chan *request)
    quit = make(chan bool)
    go server(op, service, quit)
    return service, quit
}

func testi1(a []interface{}) {
    for i := range a {
        fmt.Printf("a[%d]=%v\n", i, a[i])
    }
}

func main() {
    a := [5]int{10, 20, 30, 40, 50}

    // first-class value (copy, by-value)
    dump1cpy("a-cpy", a)

    // slice (by-reference)
    dump1("a", a[:])
    dump1ref("a-ref", a[:])
    dump1("a", a[:])

    b := make([]int, 3)
    b[0] = 1000
    b[1] = 2000
    b[2] = 3000
    dump1("b", b[:])

    c := b
    c[0] = 1001
    dump1("c", c[:])

    // maps: by-reference
    m := map[string]int{"one": 1, "two": 2}
    dump2("m", m)

    mm := map[string]int{}
    mm["three"] = 3
    mm["four"] = 4
    dump2("mm", mm)

    mmm := make(map[string]int)
    mmm["five"] = 5
    mmm["six"] = 6
    dump2("mmm", mmm)

    var mx map[string]int = m
    dump2("mx", mx)

    mmx := mm
    dump2("mmx", mmx)
    mmx["ref"] = 100
    dump2("mm", mm)

    type data struct {
        s       string
        a, b, c int
    }

    var t *data = new(data)
    (*t).a = 11
    t.b = 22
    t.s = "ivan"

    tt := &data{"", 11, 12, 13}

    if tt.s == "" {
        tt.s = "none"
    }

    d1 := 3 / 2.0
    d2 := 3 / 2

    fmt.Printf("%+v (%+v) [%.2f,%d]\n", *t, *tt, d1, d2)

    // types & interfaces
    type testT struct {
        intArray
    }

    testType := testT{intArray{1, 2, 3}}
    fmt.Printf("testT: %+v\n", testType)

    test1([]int{10, 20, 30}[:], 100, "test", 1.0)

    vec := []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
    va := intArray(vec)
    sort(va)
    fmt.Printf("%+v\n", va)

    days := new(dayArray)
    days.data = []*day{&day{10, "tue", "tuesday"}, &day{11, "wed", "wednesday"}, &day{8, "sun", "sunday"}}

    sort(days)

    for i, v := range days.data {
        fmt.Printf("%d: %+v\n", i, v)
    }

    // goroutines
    in := ping()
    s := <-pong(in)
    fmt.Printf("main received: %v\n", s)

    // server
    lock := &sync.WaitGroup{}
    adder, quit := startServer(func(a, b int) int { lock.Done(); return a + b })
    for n := 1; n < 5; n++ {
        req := &request{n, n * 10, make(chan int)}
        lock.Add(1)
        adder <- req
        res := <-req.replyc
        fmt.Printf("result: %v\n", res)
    }
    lock.Wait()
    fmt.Printf("exiting...\n")
    quit <- true

    c1 := make(chan int, 1)
    c2 := make(chan int, 1)

    go func() {
        fmt.Println("sender 1 started")
        c1 <- 1
    }()

    go func() {
        fmt.Println("sender 2 started")
        c2 <- 2
    }()

    var i1, i2 int

    it := 0

    for it < 3 {
        select {
        case i1 = <-c1:
            fmt.Printf("i1 received %d\n", i1)
            it += i1
        case i2 = <-c2:
            fmt.Printf("i2 received %d\n", i2)
            it += i2
        }
    }

    close(c1)
    close(c2)

    fmt.Printf("i1 + i2 = %d\n", it)

    a1 := []int{1, 2, 3, 4, 5}
    a2 := make([]interface{}, len(a1))
    for i1, v1 := range a1 {
        a2[i1] = v1
    }

    testi1(a2)

    fmt.Printf("%s", files("/home/irocha/.mozilla/plugins/"))
}
