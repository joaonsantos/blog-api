docker rm -f blog-db || true
docker build -f docker/Dockerfile.db -t blog-db .
docker run -d -p 5432:5432 --name blog-db blog-db
