# Mutants Go REST API

[![go](https://img.shields.io/badge/go-v1.13.X-cyan.svg)](https://golang.org/)

> A simple Golang project to verify if DNA is from Mutant or Human
>
>Developed by Camilo Díaz Jaimes

## Prerequisites

You will need the following things properly installed on your computer.

* [Git](http://git-scm.com/)
* [Go](https://golang.org/)
* [aws-cli](https://aws.amazon.com/es/cli/)

## Installation

Following you can find the instructions:

* `git clone https://github.com/camilodiazj/mutants.git` this repository.
* Change into the new directory `cd mutants`.
* Execute `go build`
* You can see the app binary in bin directory of you `$GOPATH`

## Build

Run `go build` to build the project. The build artifacts will be stored in the Root of the
directory.

## Running tests coverage

Run `go test ./... -cover` to execute the unit tests with coverage.

## Run
* If you are in Windows, is recommended to use Git Bash to run the bash commands.
* Into the root project `cd mutants`

- You need a connection to AWS DynamoDB to persist data, to achieved it You need the
  following credentials saved into your environment variables. 
  * (Attached policy with only following actions: dynamodb:PutItem, dynamodb:Scan)
  
```bash
export AWS_ACCESS_KEY_ID=provided-key-id
export AWS_SECRET_ACCESS_KEY=provided-secret-access-key
```
* Once you have credentials to DynamoDb connection, build the artifact:
```bash
go build
```
* Next, you need to change the access permission of file to be executable:
```bash
sudo chmod +x mutants
```
* Run mutants application
```bash
./mutants
```
> Is recommended to use [POSTMAN](https://www.postman.com/) to consume the Endpoints.
> Postman Collection in `mutants/postman_collection.json`  
