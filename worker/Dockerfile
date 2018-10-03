FROM golang:1.10.0-alpine

ENV worker_path /go/src/github.com/TykTechnologies/nordic-2018-lab/worker
ENV PATH $PATH:$worker_path

WORKDIR $worker_path
COPY . .

RUN go get -u ./... go build -o foo .

ENTRYPOINT foo
