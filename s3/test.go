package main

import (
	"fmt"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

func add(b *s3.Bucket, s string, n string) {
	data := []byte(s)
	err := b.Put(n, data, "text/plain", s3.BucketOwnerFull)
	if err != nil {
		log.Fatal(err)
	}
}

func list(b *s3.Bucket, p string, m string, n int) {
	res, err := b.List(p, "/", m, n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n\n", res)
}

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}

	s := s3.New(auth, aws.SAEast)
	bucket := s.Bucket("irrs3")

	fmt.Printf("S3: %#v\n\n", s)

	add(bucket, "test1", "test/sample1.txt")
	add(bucket, "test2", "test/sample2.txt")

	list(bucket, "test/", "", 1)
	list(bucket, "test/", "test/sample1.txt", 100)
}
