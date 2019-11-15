# Puppy API rip off (exercise)
A simple api that serves puppy recipes. All puppy recipes were retrieved from [recipepuppy](http://www.recipepuppy.com/)

## Configuration
After cloning you need to create a configuration file based on config.json.dist and name it config.json 

## API usage
Example:
```
http://127.0.0.1:8080/api/?i=onions&i=garlic&q=omelet&p=3
```
 
Parameter explanation:
- i : ingredients
- q : text search query
- p : page number


