
.PHONY: seed

seed:
	go build -o bin/seed cmd/migrate/seed/main.go
	./bin/seed

.PHONY: gen-docs

gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt