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

	s := s3.New(auth, aws.Region{Name: "uol", S3Endpoint: "http://api.uolos.com.br"})
	fmt.Printf("S3: %#v\n\n", s)
	bucket := s.Bucket("irr")
	list(bucket, "", "", 100)
}
