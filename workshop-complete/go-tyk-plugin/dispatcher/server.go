package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"

	"github.com/TykTechnologies/nordic-2018-lab/workshop-complete/go-tyk-plugin/hook"
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

		path := strings.Replace(obj.Request.Url, "/todos", "", -1)
		routingKey := "index"
		todo := hook.Todo{}

		switch obj.Request.Method {
		case "GET":
			if path != "/" {
				routingKey = "show"

				idString := strings.Replace(path, "/", "", 1)

				todo.ID = bson.ObjectIdHex(idString)
			}
			todo.User = obj.Session.Alias
			log.Printf("METHOD: GET path := %s, TODO: %#v", path, todo)
		case "POST":
			routingKey = "store"
			todo.User = obj.Session.Alias

			var tmpTodo hook.Todo
			_ = json.Unmarshal([]byte(obj.Request.Body), &tmpTodo)
			todo.Todo = tmpTodo.Todo

		case "DELETE":
			routingKey = "delete"
			todo.User = obj.Session.Alias
			idString := strings.Replace(path, "/", "", 1)
			todo.ID = bson.ObjectIdHex(idString)
		case "PATCH":
			// NOT IMPLEMENTED YET

		default:
			return obj, errors.New("unsupported method")
		}

		channel, _ := s.RabbitConn.Channel()
		defer channel.Close()

		bodyBytes, _ := json.Marshal(todo)

		hook.TodoRPC(channel, routingKey, bodyBytes, obj)
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
