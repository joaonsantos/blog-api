set -x

docker rm -f db && docker network rm blog
docker network create blog
docker build -f docker/Dockerfile.db -t db .
docker run \
  -e POSTGRES_DB=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  --name db --network blog -d db
