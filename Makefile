.PHONY: start
start:
	docker-compose up --build -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: status
status:
	docker-compose ps

.PHONY: restart
restart: stop start

.PHONY: test-coverage
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@rm -f coverage.out

.PHONY: test
test:
	go test ./...
