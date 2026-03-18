include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-cleanup:
	read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres && \
	  rm -rf out/pgdata && \
	  echo "Файлы окружения очищены"; \
	else \
	  echo "Очистка окружения отменена";
	fi
