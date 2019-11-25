# Puppy API rip off (exercise)
Simple api that serves recipes for puppies. This project is a step by step guide on how to create a simple api using
Go programming language. The purpose of the project is to demonstrate to newcommers the language basic features 
and concepts.

The api is based on [recipepuppy](http://www.recipepuppy.com/) free api, All puppy recipes used were also retrieved 
from there.

## Getting Started
Clone the project
```git
git clone https://github.com/georlav/recipeapi
```




After cloning you need to create a configuration file based on config.json.dist and name it config.json 

### API usage
Example:
```
http://127.0.0.1:8080/api/?i=onions,garlic&q=omelet&p=3
```
 
Parameter explanation:
- i : comma delimited ingredients
- q : text search query
- p : page number


