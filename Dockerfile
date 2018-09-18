FROM golang:1.10.0-alpine

ENV worker_path /go/src/github.com/TykTechnologies/nordic-2018-lab
ENV PATH $PATH:$worker_path

WORKDIR $worker_path
COPY . .

RUN go build -o foo ./worker

ENTRYPOINT foo
