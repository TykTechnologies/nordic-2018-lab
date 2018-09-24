package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TykTechnologies/nordic-2018-lab/worker/db"
	. "github.com/TykTechnologies/nordic-2018-lab/worker/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	rabbitConnectionStringFormat = "amqp://%s:%s@%s:%s/"
	rabbitConnectionDialError    = "unable to dial rabbit - check connection string"
	rabbitChannelOpenError       = "unable to open channel"
	rabbitExchangeDeclareError   = "unable to declare exchange"
	todosExchange                = "todos"
)

func main() {

	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("starting worker service")

	rmqUser := os.Getenv("RABBITMQ_USER")
	rmqPass := os.Getenv("RABBITMQ_PASS")
	rmqHost := os.Getenv("RABBITMQ_HOST")
	rmqPort := os.Getenv("RABBITMQ_PORT")

	actions := []string{
		"index",
		"show",
		"store",
		"update",
		"delete",
	}

	// for simplicity, we have one connection per action / one worker per queue
	// In a real environment, these would all be separate services, individually auto-scalable.
	for _, action := range actions {

		conn, err := amqp.Dial(fmt.Sprintf(rabbitConnectionStringFormat, rmqUser, rmqPass, rmqHost, rmqPort))
		FatalOnError(err, String(rabbitConnectionDialError))
		defer conn.Close()

		channel, err := conn.Channel()
		FatalOnError(err, String(rabbitChannelOpenError))
		defer channel.Close()

		err = channel.ExchangeDeclare(todosExchange, "direct", true, false, false, false, nil)
		FatalOnError(err, String(rabbitExchangeDeclareError))

		q, err := channel.QueueDeclare(action, true, false, false, false, nil)
		FatalOnError(err, nil)

		err = channel.QueueBind(action, action, todosExchange, false, nil)
		FatalOnError(err, nil)

		log.Infof("binding: queue `%s` exchange `%s` routing-key `%s`", q.Name, todosExchange, action)

		go func(action string) {
			msgs, err := channel.Consume(action, "", true, false, false, false, nil)
			if err != nil {
				log.WithError(err).Fatal("unable to consume channel")
			}

			for msg := range msgs {
				log.WithField("queue", action).Debugf("rawMsg: %s", string(msg.Body))

				todo := db.Todo{}
				err := json.Unmarshal(msg.Body, &todo)
				if err != nil {
					log.WithError(err).Errorf("junk supplied - dropping: %s", string(msg.Body))
					continue
				}

				log.WithField("queue", action).Debugf("structMsg: %#v %d", todo.ID.String(), len(todo.ID.String()))

				var replyBytes []byte

				switch action {
				case "index":
					todos, err := db.GetByUser(todo.User)
					if err != nil {
						log.WithField("queue", action).WithError(err).Error("some problem getting todos")
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)

						continue
					}

					replyBytes, err = json.Marshal(todos)
					if err != nil {
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)
						continue
					}

				case "store":
					t, err := db.Insert(todo)
					if err != nil {
						log.WithField("queue", action).WithError(err).Error("problem storing todo")
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)
						continue
					}

					replyBytes, _ = json.Marshal(t)

				case "delete":
					if err := db.DeleteByUserAndId(todo.User, todo.ID); err != nil {
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)
						log.WithError(err).Error("problem deleting todo by user and id")
						continue
					}

					replyBytes = SuccessJson(fmt.Sprintf("successfully deleted todo for user %s id %s", todo.User, todo.ID.String()))

				case "show":
					t, err := db.GetByUserAndId(todo.User, todo.ID)
					if err != nil {
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)
						log.WithError(err).Error("problem getting todo by user and id")
						continue
					}

					replyBytes, _ = json.Marshal(t)

				case "update":
					if err := db.Update(todo.User, todo.ID, todo); err != nil {
						send(channel, ErrorJson(err.Error()), msg.ReplyTo, msg.CorrelationId, action)
						log.WithError(err).Error("problem updating todo")
						continue
					}

					replyBytes = SuccessJson(fmt.Sprintf("updated %s", string(msg.Body)))
				}

				send(channel, replyBytes, msg.ReplyTo, msg.CorrelationId, action)
			}
		}(action)
	}

	done := make(chan bool)
	<-done
}

func send(channel *amqp.Channel, msg []byte, replyTo string, correlationId string, action string) {
	log.WithField("action", action).Infof("sending to `%s`: %s", replyTo, string(msg))
	channel.Publish(
		"",
		replyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          msg,
		},
	)
}
