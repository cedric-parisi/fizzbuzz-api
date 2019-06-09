# fizzbuzz-api
[![Coverage Status](https://coveralls.io/repos/github/cedric-parisi/fizzbuzz-api/badge.svg?branch=master)](https://coveralls.io/github/cedric-parisi/fizzbuzz-api?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/cedric-parisi/fizzbuzz-api)](https://goreportcard.com/report/github.com/cedric-parisi/fizzbuzz-api)[![Build Status](https://travis-ci.com/cedric-parisi/fizzbuzz-api.svg?branch=master)](https://travis-ci.com/cedric-parisi/fizzbuzz-api)


Provides a REST API to get fizzbuzz sentences.

## Getting started

Before you start, make sure `golang` and `docker-compose` are installed.

This API was written with `go 1.12` and `GO111MODULE=on`.

### installation

```
git clone https://github.com/cedric-parisi/fizzbuzz-api.git
```

## local execution

The following command will launch the REST API along with its dependencies (db, jaeger):
```
make local
```
A swagger documentation is provided with this API, to help you find exposed endpoint. Open your browser and navigate to `http://localhost:8000/swaggerui/`.

## local development

To help on the development, a docker-compose with a postgres database and a jaeger client are launched using:
```
make dev
```
> Make sure to launch `go run migrations/migrate.go up` to migrate the DB schema. This operation is needed only when you launch the dev environment for the first time, or after you removed the docker images.


To stop and remove the docker images, launch:
```
make clean
```

### database migration

See [here](./migrations/README.md) for more details about the database migration.

### run the API

Use:
```
make run
```
> Make sure to follow `local development` first.

### test

To get an HTML representation of the code coverage, use:
```
make cover
```

To launch all unit test, use:
```
make test
```

> Unit tests are using mocked interfaces to work isolated from the outside world.
mocks are generated with [golang/mock](https://github.com/golang/mock)

### integration test

Two steps are required to launch integration tests:

- `make local` to setup a complete running API.
- in a 2nd bash window, launch `make integration`, tests will be perform against the local environment.
