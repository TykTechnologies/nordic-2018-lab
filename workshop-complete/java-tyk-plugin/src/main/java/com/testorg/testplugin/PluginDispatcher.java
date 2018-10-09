package com.testorg.testplugin;

import coprocess.DispatcherGrpc;
import coprocess.CoprocessObject;

public class PluginDispatcher extends DispatcherGrpc.DispatcherImplBase {

    @Override
    public void dispatch(CoprocessObject.Object request,
            io.grpc.stub.StreamObserver<CoprocessObject.Object> responseObserver) {
        CoprocessObject.Object modifiedRequest = null;

        switch (request.getHookName()) {
            case "MyPreMiddleware":
                modifiedRequest = MyPreHook(request);
            default:
            // Do nothing, the hook name isn't implemented!
        }

        // Return the modified request (if the transformation was done):
        if (modifiedRequest != null) {
            responseObserver.onNext(modifiedRequest);
        };

        responseObserver.onCompleted();
    }

    CoprocessObject.Object MyPreHook(CoprocessObject.Object request) {
        CoprocessObject.Object.Builder builder = request.toBuilder();
        builder.getRequestBuilder().putSetHeaders("customheader", "customvalue");
        return builder.build();
    }
}

