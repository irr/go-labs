package main

/*
go get -v github.com/bradfitz/latlong
*/

import "fmt"

import "github.com/bradfitz/latlong"

func main() {
	zone := latlong.LookupZoneName(-23.643414, -46.759600)

	if len(zone) > 0 {
		fmt.Printf("Zone=%s\n", zone)
	}
}
