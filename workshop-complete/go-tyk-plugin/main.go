package main

import (
	"log"
	"net"

	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

const (
	listenAddress = ":9000"
)

func main() {

	conn, err := amqp.Dial("amqp://tyk-nordic:tyk-nordic@localhost:5672")
	fatalOnError(err)
	defer conn.Close()

	listener, err := net.Listen("tcp", listenAddress)
	fatalOnError(err)

	log.Printf("listening on tcp://%s", listenAddress)

	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Server{
		RabbitConn: conn,
	})

	fatalOnError(s.Serve(listener))
}

func fatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
