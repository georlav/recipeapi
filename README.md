![](https://github.com/georlav/recipeapi/workflows/Test/badge.svg)
[![codecov](https://codecov.io/gh/georlav/recipeapi/branch/master/graph/badge.svg)](https://codecov.io/gh/georlav/recipeapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/georlav/recipeapi)](https://goreportcard.com/report/github.com/georlav/recipeapi)
[![](https://github.com/golangci/golangci-web/blob/master/src/assets/images/badge_a_plus_flat.svg)](https://golangci.com/r/github.com/georlav/recipeapi)
[![](https://img.shields.io/badge/unicorn-approved-ff69b4.svg)](https://www.youtube.com/watch?v=9auOCbH5Ns4)

# Puppy API rip off (exercise)
Simple api that serves recipes for puppies. This project is a step by step guide on how to create a simple api using
Go programming language. The purpose of the project is to demonstrate to newcommers the language basic features 
and concepts.

The exercise is based on [recipepuppy](http://www.recipepuppy.com/) free api. All the puppy recipes used were also 
retrieved from there.

### Project is divided into nine steps, each step has its own branch.
 * 01-Server
 * 02-Configuration
 * 03-Logging
 * 04-Handlers
 * 05-Routes
 * 06-Repository
 * 07-Handler-real-data
 * 08-Middleware

There are no notes/comments for each branch. All the above branches have been merged into master

### Prerequisites
 * Go
 * Docker
 * MongoDB 3.6

API uses mongoDB as its main database, a docker compose file is provided for faster setup. You can start a container of
mongo using the following command  
```bash
docker-compose up -d
```

### Importing data
When on master branch or after reaching 07-Handler-real-data branch you can import data by running the following command
```bash
make import
``` 

### Configuration
Most of the project values can be configured by editing config.json. File is located at the project root folder

### Running the tests
```bash
make test
```

### Running the project
```bash
make run
```

### Building the project
```bash
make build
``` 

### Usage example
Following link should work when you reach 07-Handler-real-data branch or when on master
```
http://127.0.0.1:8080/api/?i=onions&i=garlic&q=omelet&p=1
```

Available Parameters explanation:
- i : ingredients
- q : text search query
- p : page number

## Authors
* **George Lavdanis** - *Initial work* - [georlav](https://github.com/georlav)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

