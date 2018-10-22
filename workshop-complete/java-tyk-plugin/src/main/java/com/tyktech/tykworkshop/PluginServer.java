package com.tyktech.tykworkshop;

import coprocess.DispatcherGrpc;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;

import java.net.*;
import java.util.concurrent.TimeoutException;

import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.Channel;

public class PluginServer {

    private static final Logger logger = Logger.getLogger(PluginServer.class.getName());
    static Server server;
    static int port = 9000;

    public static void main(String[] args) throws IOException, InterruptedException{
        System.out.println("Initializing gRPC server.");

        ConnectionFactory  factory = new ConnectionFactory();
        try {
            factory.setUri("amqp://abc:abc@server/test");
        } catch(Exception e) {
        }

        try {
            // Start the AMQP connection:
            Connection conn = factory.newConnection();

            // Initialize a new dispatcher:
            PluginDispatcher dispatcher = new PluginDispatcher();
            dispatcher.setAMQPConnection(conn);
            
            // Build and start the service:
            server = ServerBuilder.forPort(port)
                .addService(dispatcher)
                .build()
                .start();
            blockUntilShutdown();
        } catch(TimeoutException e) {
            System.out.println("timeout exception");
        };
    }

    static void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }
}

