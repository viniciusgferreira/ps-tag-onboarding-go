# Go Onboarding exercise - TAG

[Onboarding exercise](https://github.com/wexinc/ps-tag-onboarding) for Tag team.

## Installation
Make sure you have the following mandatory software installed:
- Docker && Docker-compose

To run the application, you need to run generate docs, build and run docker containers using the makefile.   
```make```   

Server is running default at localhost:8080, you can change in docker-compose.yml

## API Documentation
You can acess the Swagger API docs with the application running.   
Click here -> [Documentation](http://localhost:8080/swagger/index.html)

To generate docs, run:
```make docs```

## Technologies Used
- Go 1.21
- Gin
- MongoDB
- Swagger
- Docker