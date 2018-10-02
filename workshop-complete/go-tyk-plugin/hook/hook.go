package hook

import (
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type Todo struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	User      string        `json:"user"`
	Todo      string        `json:"todo"`
	Complete  bool          `json:"complete"`
	CreatedAt time.Time     `json:"created_at"`
}

func TodoRPC(rabbitChannel *amqp.Channel, routingKey string, bodyBytes []byte, obj *coprocess.Object) {

	// declare a random queue to receive response
	replyQ, err := rabbitChannel.QueueDeclare(
		"",    // name
		false, // durable
		true, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "HERE"))
	}

	// publish to the todos exchange the request
	err = rabbitChannel.Publish(
		"todos",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			ReplyTo: replyQ.Name,
		},
	)

	if err != nil {
		log.Fatal(errors.Wrap(err, "THERE"))
	}

	msgs, err := rabbitChannel.Consume(
		replyQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for d := range msgs {
		obj.Request.ReturnOverrides.ResponseCode = http.StatusOK
		obj.Request.ReturnOverrides.ResponseError = string(d.Body)
		obj.Request.ReturnOverrides.Headers = make(map[string]string)
		obj.Request.ReturnOverrides.Headers["Content-Type"] = "application/json"
		return
	}
}
