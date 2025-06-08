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
