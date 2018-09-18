Nordic 2018 Workshop
====================

1. Clone: `git clone https://github.com/TykTechnologies/nordic-2018-lab.git`
2. Configure: `cp .env.example.env .env` Change any secrets as per your requirements
3. Start: `docker-compose up` or `docker-compose up -d`

RMQ Management Console
----------------------
1. Admin Login: `http://localhost:15672`
2. Username: `.env(RABBITMQ_DEFAULT_USER)`
3. Password: `.env(RABBITMQ_DEFAULT_PASS)`

Mongo
-----
No auth has been set for MongoDB
