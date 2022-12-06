test-service:
	go test -v ./internal/service
create-coverage:
	go test -v -covermode=atomic -coverpkg=./... ./... -coverprofile=coverage.out 
show-coverage:	
	go tool cover -html=coverage.out
up-container:
	docker-compose --env-file .env up -d
down-container:
	docker-compose down
redis-cli: # make sure you already install redis client
	redis-cli -h localhost -p 6379 -a "pendekinredis123456"
engine:
	go build app/main.go