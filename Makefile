run: 
	./scripts/run-blog.sh

dev: scripts/run-db.sh
	./scripts/run-db.sh

clean: 
	docker rm -f blog-db
