# Mini Aspire challenge

Created by Christopher Tok for the use of Aspire job application. 

## Dev environment
- MacOS 11.2 Big Sur
- Apple Silicon (arm64)

## Pre-requisites
- Internet Connection (for downloading docker images)
- Docker
- Make
- rest api client (curl/postman/insomnia) I included insomnia json for request collections
- golangci-lint

## How to use
- run ``` go run main.go development migrate ``` for db migration
- run ``` go run main.go development seed ``` for seeding admin data
- run ``` go run main.go development server ```
- connect to ``` localhost:8000 ``` using your rest api client
- run ```golangci-lint run .``` to see golangci lint result

## Architecture
repo architecture:
- Config (to be injected to any layer / resource initialization. consist of configurations)
- Resources (to be injected to repository layer, usually client for other dependency like database)
- Handler (layer for serialization, deserialization, request validation. call usecase)
- Usecase (layer for business logic, consists of repository)
- Repository (wrapper for other library, no unit test because dependencies not mockable)
 

## Experiences with tech stacks
- Go: ~3 years
- Postgresql: ~5 years