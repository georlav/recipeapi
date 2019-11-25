# Puppy API rip off (exercise)
Simple api that serves recipes for puppies. This project is a step by step guide on how to create a simple api using
Go programming language. The purpose of the project is to demonstrate to newcommers the language basic features 
and concepts.

The exercise is based on [recipepuppy](http://www.recipepuppy.com/) free api. All the puppy recipes used were also 
retrieved from there.

### Prerequisites
 * Go
 * Docker
 * MongoDB 3.6

API uses mongoDB as its main database, a docker compose file is provided for faster setup. You can start a container of
mongo using the following command  
```bash
docker-compose up -d
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

