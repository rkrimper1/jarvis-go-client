# Jarvis API Example

This example shows how to call the Jarvis API.

The code demonstrates a simple authenticate and business message call to your local environment.

## Prerequisites
You should have Go 1.22 >  installed.

## How to Run

We make the assumption here that the jarvis api microservice is up and running in your development environment to test the gRPC client

```
export APP_ENV=local
go run main.go
```