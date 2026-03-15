# Jarvis API gRPC Example

This example shows how to call the Jarvis API.

The code demonstrates a simple authenticate and business message call to your local environment.

## Prerequisites
You should have Go 1.25 installed.

## How to run the generator
```
make
```
NOTE: The generator expects that you have already checked out a fresh copy of jarvis alongside jarvis-go-client, Note that prep_protos.sh requires that you install our protoc-gen-stripper plugin.


## How to Run The Example gRPC Client

```
export APP_ENV=local
go run main.go
```

## References