![AvitoTech Logo](https://avatars.githubusercontent.com/u/13049122?s=200&v=4)

# AvitoTech Backend Assignment

<!--toc:start-->
- [AvitoTech Backend Assignment](#avitotech-backend-assignment)
  - [Quickstart](#quickstart)
  - [Structure](#structure)
  - [Docs](#docs)
  - [Build / Run](#build-run)
    - [Targets](#targets)
  - [Configuration](#configuration)
  - [What's implemented](#whats-implemented)
    - [Primary task](#primary-task)
    - [Extra 1 - CSV Audit](#extra-1-csv-audit)
    - [Extra 2 - Segments expiration](#extra-2-segments-expiration)
    - [Extra 3 - Automatically assign % users to the segment](#extra-3-automatically-assign-users-to-the-segment)
<!--toc:end-->

Assignment for AvitoTech 2023 backend internship

## Quickstart

```bash
# If you use mage
# https://magefile.org/
mage docker:run

# Otherwise...
docker compose up
```

*This will spin up...*

- The server itself at port `1234` with OpenAPI Web-UI at path `/docs/`
- [PGWeb](https://github.com/sosedoff/pgweb) at port `8081` - Web-UI client for postgres
- PostgreSQL at port `5432` - Database

## Structure

API is developed in a design-first approach.
That is, the API is designed with OpenAPI[^1] spec defined in [openapi.yaml](./openapi.yaml) file
and then ran through the code generation with [oapi-codegen](https://github.com/deepmap/oapi-codegen/)

It implements a complete spec in [openapi.gen.go](./server/api/openapi.gen.go) file
with RPC inspired interface implemented in [server.go](./server/server.go).

The main points of this RPCish interface is to automate some parsing, abstract user code
from server specific code, and also to force user code to comply with the schema.

## Docs

- Web-UI is available at [localhost:1234/docs/](http://localhost:1234/docs/)
- OpenAPI Specification file is available at [openapi.yaml](./openapi.yaml)

## Build / Run

This project uses [Mage](https://magefile.org/) as the build tool

<details>
<summary>Why?</summary>

From the [Mage](https://magefile.org/) website...

> Makefiles are hard to read and hard to write. Mostly because makefiles are
> essentially fancy bash scripts with significant white space and
> additional make-related syntax.
>
> Mage lets you have multiple magefiles, name your magefiles whatever
> you want, and they’re easy to customize for multiple operating systems.
> Mage has no dependencies (aside from go) and runs just fine on all major
> operating systems, whereas make generally uses bash which is not well
> supported on Windows. Go is superior to bash for any non-trivial task
> involving branching, looping, anything that’s not just straight line
> execution of commands. And if your project is written in Go, why
> introduce another language as idiosyncratic as bash?
> Why not use the language your contributors are already comfortable with?

</details>

To run the server with auxiliary docker containers run...

```bash
mage docker:run

# or using docker compose
docker compose up
```

You can also run supplimentary containers only, without the server itself.
This can be useful you wan't to quickly test some new changes without restarting
other containers.

```bash
mage docker:dev

# then, in some other window you can run
mage run
```

See [Targets](#targets) for all available targets

### Targets

To show the available targets run...

```bash
mage -l
```

```
Targets:
  docker:dev     Start docker compose only with auxiliary containers (database, web ui) without the server
  docker:run     Rebuild Dockerfile and start docker compose with the server itself
  docker:test    Spin up docker containers and run tests
  generate       Run code generation
  run            Start the server
  test           Run tests
```

## Configuration

This project uses environment variables for configuration. `.env` files
are supported, but will be overwritten by the existing environment variables.
Consider values from `.env` file as a sensible defaults.

You can use [template.env](./template.env) as a template for `.env`

```env
# template.env

# Server port to use. Can omitted. In that case, 1234 will be used
SERVER_PORT=1234

# PostgreSQL connection URI
# https://www.prisma.io/dataguide/postgresql/short-guides/connection-uris
SERVER_DB_POSTGRES=postgresql://postgres:postgres@localhost:5432/db
```

You can use this template like that...

```bash
cp template.env .env
```

> [!NOTE]  
> [`docker-compose.yml`](./docker-compose.yml) already includes required
> environment variables, so that you don't need to configure anything.

## What's implemented

Everything. Primary task and all extra tasks.

### Primary task

- [x] Method for creating a user

```bash
curl --request POST \
  --url http://0.0.0.0:1234/users/1234 \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json'
```

- [x] Method for creating a segment

```bash
curl --request POST \
  --url http://0.0.0.0:1234/segments/AVITO_TEST \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json'
```

- [x] Method for deleting a segment

```bash
curl --request DELETE \
  --url http://0.0.0.0:1234/segments/AVITO_TEST \
  --header 'Accept: application/json'
```

- [x] Method for assigning a segment to a user

```bash
curl --request POST \
  --url http://0.0.0.0:1234/users/1234/segments/AVITO_TEST \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json'
```

- [x] Method for depriving a segment from a user

```bash
curl --request DELETE \
  --url http://0.0.0.0:1234/users/1234/segments/AVITO_TEST \
  --header 'Accept: application/json'
```

- [x] Method for getting active segments of a user

```bash
curl --request GET \
  --url http://0.0.0.0:1234/users/1234/segments \
  --header 'Accept: application/json'
```

### Extra 1 - CSV Audit

Request:

```bash
curl --request GET \
  --url http://0.0.0.0:1234/audit \
  --header 'Accept: text/csv'
```

Response:

```csv
user_id,segment_slug,action,stamp
1234,AVITO_TEST,ASSIGN,2023-08-28T22:29:13+03:00
1234,AVITO_TEST,DEPRIVE,2023-08-28T22:29:16+03:00
```

It will track both types of changes

- Manual deprivation *&* assignment. Done through api endpoints directly.
- Expiration and auto-assignment (from the [third task](#extra-3-automatically-assign-users-to-the-segment))

Implmented with postgres trigger that watches for `INSERT` and `DELETE` operations
on the given table and logs changes to the special audit table.
That is, the whole process is done automatically.

### Extra 2 - Segments expiration

[![asciicast](https://asciinema.org/a/ZoMo8mrVnfj3luLtk95EHt5VI.svg)](https://asciinema.org/a/ZoMo8mrVnfj3luLtk95EHt5VI)

The rows itself are not deleted from the database when they expire.
Instead, they are just filtered out on `SELECT` query.

```bash
curl --request POST \
  --url http://0.0.0.0:1234/users/1234/segments/AVITO_TEST \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{
  "expiresAt": "2017-07-21T17:32:28+03:00"
}'
```

### Extra 3 - Automatically assign % users to the segment

```bash
curl --request POST \
  --url http://0.0.0.0:1234/segments/AVITO_TEST \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{
  "outreach": 0.42
}'
```

That would assign `AVITO_TEST` segment to the 42% (`"outreach": 0.42`)
of the existing and new users. When the segment is deleted, it would be deprived
from all users automatically.

[^1]: The OpenAPI Specification is a specification language for HTTP APIs that provides a standardized means to define your API to others. https://www.openapis.org/
