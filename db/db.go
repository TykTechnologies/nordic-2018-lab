package db

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	. "github.com/TykTechnologies/nordic-2018-lab/helpers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

const (
	todosCollection = "todos"
)

func init() {
	sess, err := mgo.Dial(os.Getenv("MONGO_URL"))
	FatalOnError(err, String(fmt.Sprintf("unable to dial mongo `%s`", os.Getenv("MONGO_URL"))))
	db = sess.DB("todos")
}

type Todo struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	User      string        `json:"user"`
	Todo      string        `json:"todo"`
	Complete  bool          `json:"complete"`
	CreatedAt time.Time     `json:"created_at"`
}

func Insert(t Todo) (*Todo, error) {

	t.CreatedAt = time.Now()
	t.ID = bson.NewObjectId()
	if err := db.C("todos").Insert(t); err != nil {
		return nil, err
	}

	return &t, nil
}

func GetByUser(user string) ([]Todo, error) {

	var todos []Todo
	if err := db.C(todosCollection).Find(bson.M{
		"user": user,
	}).All(&todos); err != nil {
		return nil, errors.Wrap(err, "unable to GetByUser:All")
	}

	if len(todos) == 0 {
		return nil, errors.New(fmt.Sprintf("todos for user %s not found", user))
	}

	return todos, nil
}

func GetByUserAndId(user string, id bson.ObjectId) (Todo, error) {
	todo := Todo{}

	if err := db.C(todosCollection).Find(bson.M{
		"_id":  id,
		"user": user,
	}).One(&todo); err != nil {
		return todo, err
	}

	return todo, nil
}

func DeleteByUserAndId(user string, id bson.ObjectId) error {

	if err := db.C(todosCollection).Remove(bson.M{
		"_id":  id,
		"user": user,
	}); err != nil {
		return err
	}

	return nil
}

func Update(user string, id bson.ObjectId, todo Todo) error {

	if err := db.C(todosCollection).Update(bson.M{
		"_id":  id,
		"user": user,
	}, bson.M{
		"$set": bson.M{"complete": todo.Complete},
	}); err != nil {
		return err
	}

	return nil
}
