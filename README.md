# AvitoTech Backend Assignment

<!--toc:start-->
- [AvitoTech Backend Assignment](#avitotech-backend-assignment)
  - [Quickstart](#quickstart)
  - [Structure](#structure)
  - [Docs](#docs)
  - [Build / Run](#build-run)
    - [Targets](#targets)
  - [Configuration](#configuration)
<!--toc:end-->

Assignment for AvitoTech 2023 backend internship

## Quickstart

```bash
# If you use mage
# https://magefile.org/
mage docker:all

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
mage docker:all

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
  docker:all    Rebuild Dockerfile and start docker compose
  docker:dev    Start docker compose only with auxiliary containers (database, web ui) without the server itself
  generate      Run code generation
  run           Start the server
  test          Spin up docker containers and run tests
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

[^1]: The OpenAPI Specification is a specification language for HTTP APIs that provides a standardized means to define your API to others. https://www.openapis.org/
