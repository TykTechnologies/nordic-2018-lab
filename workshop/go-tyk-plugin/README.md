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

