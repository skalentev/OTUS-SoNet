# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.15
GOLANG_DOCKER_CONTAINER := goesquerydsl-container

#help:
#       make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[\:alnum\:]]' -e '^$@$$'
build:
        docker-compose -f docker-compose.yml build --no-cache -t sonet .
up:
        docker-compose -f docker-compose.yml up -d $(c)
down:
        docker-compose -f docker-compose.yml down $(c)
restart:
        docker-compose -f docker-compose.yml stop $(c)
        docker-compose -f docker-compose.yml up -d $(c)
logs:
        docker-compose -f docker-compose.yml logs --tail=100 -f $(c)
update:
		  git pull
          sudo docker build --no-cache -t sonet .
          sudo docker compose down
          sudo docker compose up -d
