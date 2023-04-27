# Mini Aspire challenge

Created by Christopher Tok for the use of Aspire job application. 

## Dev environment
- MacOS 11.2 Big Sur
- Apple Silicon (arm64)

## Pre-requisites
- Latest Go
- SQL database (I used postgresql)
- rest api client (curl/postman/insomnia) I included insomnia json for request collections

## How to use
- run ``` go run main.go development migrate ``` for db migration
- run ``` go run main.go development seed ``` for seed admin data (email: admin@admin.com, password: admin)
- run ``` go run main.go development server ```
- connect to ``` localhost:8000 ``` using your rest api client

### API
- register (POST /user/register)
- login (POST /user/login)
- new loan (POST /loan)
- approve loan (PUT /loan/approve) , admin only
- pay loan (POST /loan/pay)
- get loan (GET /loan)

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
