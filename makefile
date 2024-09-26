SERVICES = rust-actix

COMPOSE_FILE_STANDARD = docker-compose.yml

TEST_COMMAND = pytest tests/common_tests.py

LOCAL_DATABASE = ./database/docker-compose.yml

DATABASE_TABLES = users, followers, messages, latest


.PHONY: start-local-db
start-local-db:
	@echo "Starting the database container..."
	docker-compose -f $(LOCAL_DATABASE) up -d

.PHONY: stop-local-db
stop-local-db:
	@echo "Stopping and removing the database container..."
	@docker-compose -f $(LOCAL_DATABASE) stop
	@docker-compose -f $(LOCAL_DATABASE) rm -f local_database

.PHONY: clean-db
clean-db:
	@echo "Cleaning the database..."
	@export PGPASSWORD=pass
	@docker-compose -f $(LOCAL_DATABASE) exec local_database psql -U user -d waect -c "TRUNCATE TABLE $(DATABASE_TABLES) CASCADE;"

.PHONY: start-service
start-service:
	@echo "Starting service with Docker Compose file: $(SERVICE)..."
	docker-compose -f $(COMPOSE_FILE_STANDARD) up -d

.PHONY: stop-service
stop-service:
	@echo "Stopping and removing service with Docker Compose file: $(SERVICE)..."
	docker-compose -f $(COMPOSE_FILE_STANDARD) stop
	docker-compose -f $(COMPOSE_FILE_STANDARD) rm -f

.PHONY: test-service
test-service:
	@echo "Running tests for service..."
	$(TEST_COMMAND)

.PHONY: test-all test-service
test-all: start-local-db
	@for service in $(SERVICES); do \
		echo "=========================="; \
		echo "Testing service: $$service..."; \
		echo "=========================="; \
		if [ -d "$$service" ]; then \
			cd $$service && echo "Starting service: $$service..." && docker-compose up -d; \
			cd .. && $(MAKE) test-service; \
			cd $$service && echo "Stopping service: $$service..." && docker-compose stop; \
			cd $$service && echo "Cleaning up service: $$service..." && docker-compose rm -f; \
		else \
			echo "Directory $$service does not exist."; \
			$(MAKE) stop-local-db; \
			exit 1; \
		fi; \
	done
	$(MAKE) stop-local-db
	@echo "All services tested successfully!"



