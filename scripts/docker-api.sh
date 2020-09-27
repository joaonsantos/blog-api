set -x
docker rm -f backend
docker build -f docker/Dockerfile.db -t db .
docker run \
  -e PG_URL=postgres://postgres:postgres@db/postgres \
  --name backend --network backend -p 8000:8000 -d backend
