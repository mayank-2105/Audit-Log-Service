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

### API for creating event

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

<img width="1063" alt="Screenshot 2023-08-16 at 12 30 39 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/f2435c39-2f51-4c02-8575-351c5f44bc43">

Here is another example you can use with different data.

```
curl --location 'http://localhost:8080/event' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIxODIxMDd9.la6KCY__xjgGtAZQkfh_OXVEG0zSvjIQeNy1JpcpZ28' \
--header 'Content-Type: application/json' \
--data '{
         "type": "ERROR",
         "action": "account_deleted",
         "data": {
             "msg": "user does not exist",
             "maxretry": 3
         }
}'
```
<img width="1063" alt="Screenshot 2023-08-16 at 12 28 24 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/0cf9249c-ad58-4ab2-8efe-34114d6be7ee">


### API for querying on event

