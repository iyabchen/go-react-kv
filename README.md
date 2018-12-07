# Web app demo with golang and react

A REST web server to key-value pair CRUD operations.
A React UI is served with the web server at root directory.

## Build

`make clean` - Clean dependencies and test results for server.
`make build` - Build server binary in `./server/bin` folder.
`make test` - Run unit test on golang.  
`make docker` - Build docker image with ui for server application.

## REST API

Rest API can be found in swagger.yml.

## Usage

Run with docker:

```bash
$ make docker
$ docker run --name=app -p 18080:8080 -d demo:1.0
```

Demo can be found at https://murmuring-ridge-87096.herokuapp.com/.
