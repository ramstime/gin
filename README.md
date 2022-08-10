# Golang Gin Web Framework example
Gin Web Framework example with Redis DB and unit tests

# run/compile
go mod tidy

go mod vendor

go run main.go

go build -o bin/webserver main.go

# run redis DB
docker pull redis

docker run --name redis-test-instance -p 6379:6379 -d redis

# test webserver
curl -X GET   http://localhost:8080/user/ja*

curl -X POST   http://localhost:8080/admin   -H 'authorization: Basic Zm9vOmJhcg=='   -H 'content-type: application/json'   -d '{"name": "rams", "age":35, "books": ["solid","galaxy"] }'

curl -X GET   http://localhost:8080/user/rams
