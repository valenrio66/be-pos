run:
	go run cmd/api/main.go

.PHONY: migration

migration:
	@if [ -z "$(name)" ]; then \
  		echo "Error: Migration name cannot be left blank!"; \
		echo "How to use: make migration name=migration_file_name"; \
		exit 1; \
	fi
	dbmate new $(name)