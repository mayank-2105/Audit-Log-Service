go mod init audit-log-service
go mod tidy
go get github.com/dgrijalva/jwt-go
go get "github.com/gorilla/mux"
go get "github.com/olivere/elastic/v7" 
go build .
./audit-log-service