package main

import (
	"net"
	"os"

	"grpc_using_go/book"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/*
sudo apt install protobuf-compiler
go get -u github.com/golang/protobuf/protoc-gen-go

go install github.com/fullstorydev/grpcurl/cmd/grpcurl

> to create book.pb.go
protoc --go_out=plugins=grpc:book book.proto
*/

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	var port string
	var ok bool
	port, ok = os.LookupEnv("PORT")
	if ok {
		log.WithFields(log.Fields{
			"PORT": port,
		}).Info("PORT env var defined")

	} else {
		port = "9000"
		log.WithFields(log.Fields{
			"PORT": port,
		}).Warn("PORT env var not defined. Going with default")

	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("Failed to listen")
	}

	s := book.Server{}

	grpcServer := grpc.NewServer()

	/*
		grpcurl -plaintext localhost:9000 list
		grpcurl -plaintext localhost:9000 describe
		grpcurl -d '{"name": "To Kill a Mockingbird", "isbn": 12345}' -plaintext localhost:9000 book.BookService/GetBook
	*/
	reflection.Register(grpcServer)

	book.RegisterBookServiceServer(grpcServer, &s)

	log.Info("gRPC server started at ", port)
	if err := grpcServer.Serve(l); err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("Failed to serve")
	}

}
