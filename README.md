# Audit-Log-Service

## Dependencies

Go 
ElasticSearch

## Setup

```
git status
git add
git commit


1. go mod init audit-log-service
2. go mod tidy
3. go get github.com/dgrijalva/jwt-go
4. go get "github.com/gorilla/mux"
5. go get "github.com/olivere/elastic/v7"
6. ./bin/elasticsearch -E xpack.security.enabled=false
7. go build .
8. ./audit-log-service
```
   

## Testing

API for login

You first need to call the /login API which will return the authentication token which is valid for 24 hours.

Below is the cURL for login API

```curl --location --request POST 'http://localhost:8080/login'```

API for creating event


API for querying on event

