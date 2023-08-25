# KODE Test task

Test task for the [KODE](https://kode.ru/)

- [Openapi Schema](./openapi.yaml)
- [Postman Collection](./postman_collection.json)

## Task

> Golang a service that provides a REST API interface with methods
> 
> - Add notes
> - Output list of notes
> 
> Data should be stored in PostgreSQL. 
> When saving notes it is necessary to validate spelling errors using
> [Yandex.Speller](https://yandex.ru/dev/speller/) service (add integration with the service).
> 
> It is also necessary to implement authentication and authorization. Users should have access only to their notes.

## Run

This project uses [Mage](https://magefile.org/) as the build tool.

```
Targets:
  docker      Rebuild Dockerfile and start docker compose
  generate    Run code generation
  test        Spin up docker containers and run tests
```

To start the server run...

```bash
mage docker

# or... 
docker compose up
```

This will spin up

- Redis
- Postgres
- [Adminer](https://www.adminer.org/) - DB Web UI
- Swagger UI on port `8082`. You can visit on `http://localhost:8082`. However, "Try it out" feature won't work due to CORS.

## Configuration

The server is configured through environment variabled.
See [template.env](./template.env) for reference.

```env
# port to listen on
SERVER_PORT=1234

# jwt secret to sign tokens
SERVER_JWT_SECRET=top-secret

# postgres uri
SERVER_DB_POSTGRES=postgresql://postgres:postgres@localhost:5432/db

# redis uri
SERVER_DB_REDIS=redis://localhost:6379?protocol=3
```

You can use this template like this...

```bash
cp template.env .env
```

## Test

Note, that it requires a valid configuration.
You can do that by simply copying [`template.env`](./template.env) into `.env`


```bash
cp template.env .env
mage test

# or...
cp template.env .env
docker compose up -d --no-deps --build server
docker compose up -d
sleep 5
go test ./...
docker compose down
```

This will run [`api_test.go`](./server/api/api_test.go)

## About

Server is developed schema-first with [OpenAPI](https://www.openapis.org/). That is, the REST API is designed in [openapi.yaml](./openapi.yaml)
file and then ran through the code generation with [oapi-codegen](https://github.com/deepmap/oapi-codegen/)

It uses [chi router](https://github.com/go-chi/chi) under the hood.

### Authentication

Authentication is implemented with JWT with Access *&* Refresh tokens.
Access token is valid for 5 minutes after it was issued. Refresh token is
valid for 24 hours.

Refresh tokens are stored in the Redis, while the users (username, password)
are stored in the Postgres DB.

