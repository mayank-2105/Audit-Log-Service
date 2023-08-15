# Audit-Log-Service

## Dependencies

* Go 
* ElasticSearch

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

**Important- You first need to call the login API before testing other APIs which will return the authentication token which is valid for 24 hours.**

Below is the cURL for login API

```curl --location --request POST 'http://localhost:8080/login'```

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



### API for querying on event

Query for all the events

```
curl --location 'http://localhost:8080/query' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIxODIxMDd9.la6KCY__xjgGtAZQkfh_OXVEG0zSvjIQeNy1JpcpZ28' \
--data ''
```

Query for common event data

```
curl --location 'http://localhost:8080/query?action=account_deleted' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIxODIxMDd9.la6KCY__xjgGtAZQkfh_OXVEG0zSvjIQeNy1JpcpZ28' \
--data ''
```

Query for specific event data

```
curl --location 'http://localhost:8080/query?data.age=44' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIxODIxMDd9.la6KCY__xjgGtAZQkfh_OXVEG0zSvjIQeNy1JpcpZ28' \
--data ''
```


## Here are some of the screenshots for all these APIs when tested on Postman

<img width="1063" alt="Screenshot 2023-08-16 at 12 30 39 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/f2435c39-2f51-4c02-8575-351c5f44bc43">

<img width="1060" alt="Screenshot 2023-08-15 at 11 15 58 PM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/82372fae-ab3c-4f66-94f9-2c4ea8a57b8e">

<img width="1063" alt="Screenshot 2023-08-16 at 12 28 24 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/0cf9249c-ad58-4ab2-8efe-34114d6be7ee">

<img width="1058" alt="Screenshot 2023-08-16 at 12 44 48 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/90e48be2-2508-43a5-85e1-15fffc82e583">

<img width="1058" alt="Screenshot 2023-08-16 at 12 46 36 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/3bac839b-f12b-47e9-9cf6-ae7bfffaaa5c">

<img width="1062" alt="Screenshot 2023-08-16 at 12 48 40 AM" src="https://github.com/mayank-2105/Audit-Log-Service/assets/72939306/8204678f-2ec0-4d11-b0d5-a248394f9b8b">
