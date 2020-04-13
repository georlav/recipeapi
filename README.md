![Tests](https://github.com/georlav/recipeapi/workflows/Tests/badge.svg?branch=master)
![GolangCI](https://github.com/georlav/recipeapi/workflows/GolangCI/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/georlav/recipeapi/branch/master/graph/badge.svg)](https://codecov.io/gh/georlav/recipeapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/georlav/recipeapi)](https://goreportcard.com/report/github.com/georlav/recipeapi)
[![](https://img.shields.io/badge/unicorn-approved-ff69b4.svg)](https://www.youtube.com/watch?v=9auOCbH5Ns4)

# Puppy API rip off (exercise)
Simple api that serves recipes for puppies. This project is a step by step guide on how to create a simple api using
Go programming language. The purpose of the project is to demonstrate to new comers the language basic features and
concepts.

The exercise is based on [recipepuppy](http://www.recipepuppy.com/) free api. All the puppy recipes used were also 
retrieved from there.

### Project is divided into nine steps, each step has its own branch.
 * 00-CI
 * 01-Server
 * 02-Configuration
 * 03-Logging
 * 04-Handlers
 * 05-Routes
 * 06-Database
 * 07-Handler-real-data
 * 08-Middleware
 * 09-Authentication

### Prerequisites
 * Go
 * Docker

### Setup databases
Set up required databases. At each execution recreates database and re imports data.
```bash
make db
```

### Starting database container
```bash
docker-compose up -d
```

Set up/reset only dev-db
```bash
make db-dev
```
Set up/reset only test-db
```bash
make db-test
```

### Configuration
Most of the project values can be configured by editing config.yaml, config file is located under the project 
root folder. 

#####You can use the following config format types
* JSON 
* TOML
* YAML
* HCL
* envfile

You can override any config value by using environmental values, for example to override token.ttl you need to set an 
environmental variable with key RECIPE_TOKEN_TTL

Project is using [Viper config](https://github.com/spf13/viper) a complete configuration solution for Go 
applications

### Running tests
```bash
make test
```

### Running project
```bash
make run
```

### Building project
```bash
make build
```

### Running linter
```bash
make lint
```
if you are behind a corporate firewall using a custom certificate use
```bash
make lint-insecure
```

Project is using [golang lint](https://github.com/golang/lint) linter

### Usage examples

Get recipe
```
http://127.0.0.1:8080/api/recipe/1 [GET]
```

Get Recipes
```
http://127.0.0.1:8080/api/recipes?ingredient=onions&ingredient=garlic&term=omelet&page=1 [GET]
```

User Sign up
```
http://127.0.0.1:8080/api/user/signup [POST]

{
    "email":"email@email.com",
    "username":"username2",
    "password":"password",
    "repeatPassword":"password",
    "fullName":"test user"
}
```

User Sign in
```
http://127.0.0.1:8080/api/user/signin [POST][body]
{
    "username": "username1",
    "password": "password"
}
```

User Profile 
```
http://127.0.0.1:8080/api/user [GET]
```

### Postman
For your convenience Postman collection/environment files are available at
```
api/Recipes.postman_collection.json
api/Recipes.postman_environment.json
``` 

Available Parameters explanation:
- ingredient : list of ingredients
- term : text search in titles
- page : page number

### Swagger Docs
You can view swagger docs after running the app here [http://127.0.0.1:8080/swagger/index.html](http://127.0.0.1:8080/swagger/index.html)

Project is using [swaggo/swag](https://github.com/swaggo/swag) to generate/update swagger docs you can install 
the cli using the following command
```bash
go get github.com/swaggo/swag/cmd/swag
```

Update swagger docs
```bash
swag init -g cmd/api/main.go -o api/swagger
```
or
```bash
make swag
```

## Authors
* **George Lavdanis** - *Initial work* - [georlav](https://github.com/georlav)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

