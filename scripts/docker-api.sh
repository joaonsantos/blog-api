set -x
docker rm -f backend
docker build -f docker/Dockerfile -t backend .
docker run \
  -e PG_URL=postgres://postgres:postgres@db/postgres \
  --name backend --network blog -p 8000:8000 -d backend
