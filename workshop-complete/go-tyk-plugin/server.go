package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"

	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

type Server struct {
	RabbitConn *amqp.Connection
}

func (s Server) Dispatch(ctx context.Context, obj *coprocess.Object) (*coprocess.Object, error) {

	switch obj.HookName {
	case "TodoRabbitHook":

		//objJson, _ := json.MarshalIndent(obj, "", "  ")
		//log.Printf("%s", string(objJson))

		// stripping the listen path
		// /todos/12345abc -> /12345abc
		path := strings.Replace(obj.Request.Url, "/todos", "", -1)

		// get the id from the path
		// /12345abc -> 12345abc
		idString := strings.Replace(path, "/", "", 1)

		routingKey := ""
		todo := Todo{}

		// Handle request routing and build the request
		switch obj.Request.Method {
		case http.MethodGet:
			if path == "/" {
				routingKey = "index"
				todo.User = obj.Session.Alias
			} else {
				routingKey = "show"

				if !bson.IsObjectIdHex(idString) {
					obj.Request.ReturnOverrides.ResponseCode = http.StatusBadRequest
					obj.Request.ReturnOverrides.ResponseError = `{"error": "invalid id"}`
					obj.Request.ReturnOverrides.Headers = map[string]string{
						"Content-Type": "application/json",
					}
					return obj, nil
				}
				todo.ID = bson.ObjectIdHex(idString)
			}
		case http.MethodPost:
			routingKey = "store"

			_ = json.Unmarshal([]byte(obj.Request.Body), &todo)
		case http.MethodDelete:
			routingKey = "delete"
			if !bson.IsObjectIdHex(idString) {
				obj.Request.ReturnOverrides.ResponseCode = http.StatusBadRequest
				obj.Request.ReturnOverrides.ResponseError = `{"error": "invalid id"}`
				obj.Request.ReturnOverrides.Headers = map[string]string{
					"Content-Type": "application/json",
				}
				return obj, nil
			}
			todo.ID = bson.ObjectIdHex(idString)
		case http.MethodPatch:
			routingKey = "update"

			if !bson.IsObjectIdHex(idString) {
				obj.Request.ReturnOverrides.ResponseCode = http.StatusBadRequest
				obj.Request.ReturnOverrides.ResponseError = `{"error": "invalid id"}`
				obj.Request.ReturnOverrides.Headers = map[string]string{
					"Content-Type": "application/json",
				}
				return obj, nil
			}

			_ = json.Unmarshal([]byte(obj.Request.Body), &todo)
		default:
			return obj, errors.New("unsupported method")
		}

		// regardless of what the user posted was their user,
		// we set user to that of the JWT sub claim
		todo.User = obj.Session.Alias

		channel, _ := s.RabbitConn.Channel()
		defer channel.Close()

		bodyBytes, _ := json.Marshal(todo)

		TodoRPC(channel, routingKey, bodyBytes, obj)
	default:
		log.Printf("hook not implemented %s", obj.HookName)
	}

	return obj, nil
}

func (s Server) DispatchEvent(ctx context.Context, obj *coprocess.Event) (*coprocess.EventReply, error) {
	log.Println("DispatchEvent called")

	fmt.Println(obj.Payload)

	unquoted, err := strconv.Unquote(obj.Payload)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("%s", unquoted)

	return &coprocess.EventReply{}, nil
}
