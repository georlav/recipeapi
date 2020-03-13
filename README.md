![Tests](https://github.com/georlav/migrate/workflows/Tests/badge.svg?branch=master)
![GolangCI](https://github.com/georlav/migrate/workflows/GolangCI/badge.svg?branch=master)
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
 * 01-Server
 * 02-Configuration
 * 03-Logging
 * 04-Handlers
 * 05-Routes
 * 06-Database
 * 07-Handler-real-data
 * 08-Middleware
 * 09-CI

### Prerequisites
 * Go
 * Docker

### Setup database and import schema
```bash
make db
```

### Starting database container
```bash
docker-compose up -d
```
   
### Importing data
When on master or after reaching 07-Handler-real-data branch you can import data by running the following 
commands. 

Start the API.
```bash
go run cmd/api/main.go
``` 

Then run the import cmd. Import reads recipes.json, creates recipe requests from read data and posts requests to recipe
 API.
```bash
go run cmd/import/main.go
``` 

### Configuration
Most of the project values can be configured by editing config.json, config file is located under the project 
root folder.

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

### Usage example
Following link should work when you reach 07-Handler-real-data branch or when on master
```
http://127.0.0.1:8080/api/?ingredient=onions&ingredient=garlic&term=omelet&page=1
```

Available Parameters explanation:
- ingredient : list of ingredients
- term : text search in titles
- page : page number

## Authors
* **George Lavdanis** - *Initial work* - [georlav](https://github.com/georlav)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details