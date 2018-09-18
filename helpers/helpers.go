package helpers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func String(str string) *string {
	return &str
}

func StringValue(v *string) string {
	return *v
}

func FatalOnError(err error, msg *string) {
	if err != nil {
		if msg != nil {
			log.WithError(err).Fatal(StringValue(msg))
		} else {
			log.Fatal(err.Error())
		}
	}
}

type errorStruct struct {
	Error string `json:"error"`
}

type success struct {
	Message string
}

func ErrorJson(msg string) []byte {

	err := errorStruct{
		Error: msg,
	}
	bytes, _ := json.Marshal(err)

	return bytes
}

func SuccessJson(msg string) []byte {
	err := success{
		Message: msg,
	}
	bytes, _ := json.Marshal(err)

	return bytes
}
