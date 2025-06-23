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

.PHONY: test-coverage-html
test-coverage-html:
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@go tool cover -html=coverage.out -o coverage.html
	@rm -f coverage.out
	@echo "Coverage report generated: coverage.html"
	@open coverage.html
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
