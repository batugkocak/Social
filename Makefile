
.PHONY: seed

seed:
	go build -o bin/seed cmd/migrate/seed/main.go
	./bin/seed