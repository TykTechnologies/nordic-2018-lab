Nordic APIs - Tyk Workshop 2
============================

## Introduction

`main.go`

The entrypoint to the Tyk gRPC middleware plugin.

We first establish a connection to RabbitMQ, then start a gRPC server listening on `tcp://0.0.0.0:9000`

This gRPC server registers a DispatcherServer which we pass in the connection to Rabbit.

`server.go`

This is where you will be building out all your functionality.


