run: main.go
	go run main.go

dev: scripts/run-db.sh
	./scripts/run-db.sh

clean: 
	docker rm -f blog-db
