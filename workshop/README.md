Nordic APIs - Tyk Workshop 2
============================

Configuring your API:

1. Configure an API with the listen path `todos` via Tyk Dashboard to accept JWT auth
    - for simplicity, choose shared secret - and give a simple password.
    - set the subject to `sub`
    - set the policy to `pol`
2. Configure a Policy to grant access to your API.
3. Send an API call to your api without a token
4. Generate a JWT via https://jwt.io
    - You should add a couple of claims `exp`, `sub` and `pol`
    - Remember to select `HS256` and enter your shared secret to sign the JWT.
5. Test that your JWT works

---

## Building your plugin:

In this part of the workshop, we will be building a gRPC plugin, which will interface with
some todo's microservices sitting behind RabbitMQ.

We need to modify your API definition to tell the gateway to pass requests to the gRPC
server after the gateway has validated the JWT token.

Update the gateway configuration to point to your gRPC server and restart the service

```
"coprocess_options": {
  "enable_coprocess": true,
  "coprocess_grpc_server": "tcp://127.0.0.1:9000"
}
```

From the dashboard, go to your todos api, and click `Raw API Definition`

Change `custom_middleware.driver` from the empty string `""` to `"grpc"`.

Add the following object to the `custom_middleware.post_key_auth` array.

```json
{
  "name": "TodoRabbitHook"
}
```

Now, for every request to that API, the gRPC server will be invoked.

We now need to implement a couple of methods they should work as follows:

Save a TODO:
```
curl gateway:8080/todos/ -H "Authorization: YOUR_JWT" -d '{"todo": "Start this workshop"}'
{"id":"5bb3edf026f06900085d2771","user":"asoorm","todo":"Start this workshop","complete":false,"created_at":"2018-10-02T22:15:12.9646878Z"}
```

Save another TODO:
```
curl gateway:8080/todos/ -H "Authorization: YOUR_JWT" -d '{"todo": "Finish this workshop"}'
{"id":"5bb3ee3326f06900085d2772","user":"asoorm","todo":"Finish this workshop","complete":false,"created_at":"2018-10-02T22:15:12.9646878Z"}
```

List your user's TODOs
```
curl gateway:8080/todos/ -H "Authorization: YOUR_JWT" | python -mjson.tool               
[
    {
        "complete": false,
        "created_at": "2018-10-02T22:15:12.964Z",
        "id": "5bb3edf026f06900085d2771",
        "todo": "Finish this workshop",
        "user": "asoorm"
    },
    {
        "complete": false,
        "created_at": "2018-10-02T22:16:19.756Z",
        "id": "5bb3ee3326f06900085d2772",
        "todo": "Start this workshop",
        "user": "asoorm"
    }
]
```

Get a single TODO
```
curl gateway:8080/todos/5bb3ee3326f06900085d2772 -H "Authorization: YOUR_JWT" | python -mjson.tool
{
    "complete": false,
    "created_at": "2018-10-02T22:16:19.756Z",
    "id": "5bb3ee3326f06900085d2772",
    "todo": "Start this workshop",
    "user": "asoorm"
}
```

Update the TODO
```
curl -X PATCH gateway:8080/todos/5bb3ee3326f06900085d2772 -H "Authorization: YOUR_JWT" -d '{"complete": true}'
{
    "Message": "updated {\"id\":\"5bb3ee3326f06900085d2772\",\"user\":\"asoorm\",\"todo\":\"\",\"complete\":true,\"created_at\":\"0001-01-01T00:00:00Z\"}"
}
```

Delete the TODO
```
curl -X DELETE gateway:8080/todos/5bb3ee3326f06900085d2772 -H "Authorization: YOUR_JWT"
{
    "Message": "successfully deleted todo for user asoorm id ObjectIdHex(\"5bb3ee3326f06900085d2772\")"
}
```

Hints:

- The JWT claim `sub` contains your user's username. This is stored inside `obj.Session.Alias` of the CoProcess object.
- try dumping the coprocess object as json to the console to see what's available for the request:

```text
objJson, _ := json.MarshalIndent(obj, "", "  ")
log.Printf("%s", string(objJson))

{
  "hook_type": 3,
  "hook_name": "TodoRabbitHook",
  "request": {
    "headers": {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc29vcm0iLCJuYW1lIjoiQWhtZXQgU29vcm1hbGx5IiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjMwMDAwMDAwMDAsInBvbCI6IjViYjM3YzMxM2YwM2QzNDQxMTdhYmM3NyJ9.boZowZ6dx_Sg-9lGIXjgxb46XS4mxH_ztF1_dpQbrao",
      "Host": "gateway.ahmet:8080"
    }
    "url": "/todos/",
    "return_overrides": {
      "response_code": -1
    },
    "method": "GET",
    "request_uri": "/todos/"
  },
  "session": {
    "allowance": 1000,
    "rate": 1000,
    "per": 60,
    "expires": 3000000000,
    "quota_max": -1,
    "quota_renewal_rate": 3600,
    "access_rights": {
      "c6d940d9faeb455840d33a97c895525f": {
        "api_name": "todos",
        "api_id": "c6d940d9faeb455840d33a97c895525f",
        "versions": [
          "Default"
        ]
      }
    },
    "org_id": "5b296ceb3f03d310fffc9b9d",
    "basic_auth_data": {},
    "jwt_data": {},
    "monitor": {},
    "metadata": {
      "TykJWTSessionID": "5b296ceb3f03d310fffc9b9d857ce552473343e103c71b5e961d9601"
    },
    "alias": "asoorm",
    "last_updated": "1538489393",
    "apply_policies": [
      "5bb37c313f03d344117abc77"
    ]
  },
  "metadata": {
    "TykJWTSessionID": "5b296ceb3f03d310fffc9b9d857ce552473343e103c71b5e961d9601"
  },
  "spec": {
    "APIID": "c6d940d9faeb455840d33a97c895525f",
    "OrgID": "5b296ceb3f03d310fffc9b9d",
    "config_data": "{}"
  }
}
```
