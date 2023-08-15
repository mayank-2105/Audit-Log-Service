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

### API for login

> You first need to call the /login API which will return the authentication token which is valid for 24 hours.

Below is the cURL for login API

```curl --location --request POST 'http://localhost:8080/login'```

<img width="1060" alt="Screenshot 2023-08-15 at 11 15 58 PM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/82372fae-ab3c-4f66-94f9-2c4ea8a57b8e">

API for creating event

Below is the curl to create an event, here data field(JSON) corresponds to event specific data and can vary across events.

**Make sure to replace with the correct auth token you receive after calling login API**


```
curl --location 'http://localhost:8080/event' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIxODIxMDd9.la6KCY__xjgGtAZQkfh_OXVEG0zSvjIQeNy1JpcpZ28' \
--header 'Content-Type: application/json' \
--data '{
         "type": "INFO",
         "action": "account_created",
         "data": {
             "age": 46,
             "weight": 68,
             "occupation" : "painter"
         }
}'
```

API for querying on event

