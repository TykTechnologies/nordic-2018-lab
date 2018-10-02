package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/streadway/amqp"
)

type Todo struct {
	ID        string    `json:"id,omitempty"`
	User      string    `json:"user"`
	Todo      string    `json:"todo"`
	Complete  bool      `json:"complete"`
	CreatedAt time.Time `json:"created_at"`
}

func TodoRPC(rabbitChannel *amqp.Channel, routingKey string, bodyBytes []byte, obj *coprocess.Object) {

	// declare a random queue to receive response
	replyQ, _ := rabbitChannel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	// publish to the todos exchange the request
	_ = rabbitChannel.Publish(
		"todos",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			ReplyTo:     replyQ.Name,
		},
	)

	// consume from that temporary queue we just created
	msgs, _ := rabbitChannel.Consume(
		replyQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for d := range msgs {

		obj.Request.ReturnOverrides.ResponseCode = http.StatusBadRequest
		if !isErrorRes(d.Body) {
			obj.Request.ReturnOverrides.ResponseCode = http.StatusOK
		}

		obj.Request.ReturnOverrides.ResponseError = string(d.Body)
		obj.Request.ReturnOverrides.Headers = map[string]string{
			"Content-Type": "application/json",
		}

		return
	}
}

type ErrorStruct struct {
	Error string `json:"error"`
}

func isErrorRes(res []byte) bool {
	var errObj ErrorStruct

	err := json.Unmarshal(res, &errObj)
	if err != nil {
		log.Println(err.Error())
	}

	return errObj.Error != ""
}
