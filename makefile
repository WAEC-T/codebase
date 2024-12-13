ALL_SERVICES = python-flask c-sharp-razor go-gorilla ruby-sinatra rust-actix javascript-express go-gin

COMPOSE_FILE_STANDARD = compose-test.yml
TEST_COMMAND = pytest tests/test_flash_messages.py tests/test_api_endpoints.py
LOCAL_DATABASE = ./database/docker-compose.yml
DATABASE_TABLES = users, followers, messages, latest
DELAY_TEST_EXECUTION_SECONDS = 6

WHITE = \033[0;37m
CYAN = \033[0;36m
PINK = \033[0;35m
BLUE = \033[0;34m
RED = \033[0;31m
YELLOW = \033[1;33m
GREEN = \033[0;32m
RESET = \033[0m


.PHONY: start-local-db
start-local-db:
	@echo "$(BLUE)Starting the database container...$(RESET)"
	@docker-compose -f $(LOCAL_DATABASE) up -d > /dev/null 2>&1

.PHONY: stop-local-db
stop-local-db:
	@echo "\n$(BLUE)Stopping and removing the database container...$(RESET)"
	@docker-compose -f $(LOCAL_DATABASE) stop > /dev/null 2>&1
	@docker-compose -f $(LOCAL_DATABASE) rm -f database > /dev/null 2>&1

.PHONY: clean-db
clean-db:
	@echo "$(PINK)Cleaning the database...$(RESET)"
	@export PGPASSWORD=pass
	@docker-compose -f $(LOCAL_DATABASE) exec database psql -U user -d waect -c "TRUNCATE TABLE $(DATABASE_TABLES) > /dev/null 2>&1;"

.PHONY: start-service
start-service:
	@echo "$(CYAN)Spinning service and running tests...$(RESET) \n"
	@docker-compose -f ./$(SERVICE)/$(COMPOSE_FILE_STANDARD) up -d > /dev/null 2>&1

.PHONY: stop-service
stop-service:
	@echo "\n$(CYAN)Stopping and removing $(SERVICE)..$(RESET)"
	@docker-compose -f ./$(SERVICE)/$(COMPOSE_FILE_STANDARD) stop > /dev/null 2>&1
	@docker-compose -f ./$(SERVICE)/$(COMPOSE_FILE_STANDARD) rm -f > /dev/null 2>&1

.PHONY: test-single-service
test-single-service:
	@echo "\n$(BLUE)=====================================$(RESET)"
	@echo "$(PINK)Testing service: $(YELLOW)$(SERVICE)...$(RESET)"
	@echo "$(BLUE)=====================================$(RESET) \n"
	@if [ -d "$(SERVICE)" ] && [ -f "$(SERVICE)/$(COMPOSE_FILE_STANDARD)" ]; then \
		$(MAKE) -s start-service SERVICE=$(SERVICE) && sleep $(DELAY_TEST_EXECUTION_SECONDS); \
		$(TEST_COMMAND) || { echo "$(RED)Tests failed for $(SERVICE). Exiting.$(RESET)"; exit 1; }; \
		$(MAKE) -s stop-service SERVICE=$(SERVICE); \
	else \
		if [ ! -d "$(SERVICE)" ]; then \
			echo "$(WHITE)Skipping - directory $(SERVICE) does not exist.$(RESET)"; \
		else \
			echo "$(WHITE)Skipping - docker compose file $(COMPOSE_FILE_STANDARD) does not exist in $(SERVICE).$(RESET)"; \
		fi; \
	fi

.PHONY: test-all
test-all: start-local-db
	@set -e; \
	for service in $(ALL_SERVICES); do \
		$(MAKE) -s test-single-service SERVICE=$$service; \
	done
	@$(MAKE) -s stop-local-db
	@echo "$(GREEN)All services tested!$(RESET)"

.PHONY: test-service
test-service: start-local-db
	@set -e; \
	services=$$(echo "$(MAKECMDGOALS)" | tr ' ' '\n' | grep -v '^test-service$$'); \
	for service in $$services; do \
		if echo "$(ALL_SERVICES)" | grep -wq "$$service"; then \
			$(MAKE) -s test-single-service SERVICE=$$service; \
		else \
			echo "$(RED)Service $$service is not defined as a base service and is not valid.$(RESET)"; \
		fi; \
	done; \
	echo "$(GREEN)[$$services] tested$(RESET)" | tr '\n' ', ';
	@$(MAKE) -s stop-local-db

# Don't remove these weird thing here because it prevents the services to be treated as targets by make ;)
%:
	@: