# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.15
GOLANG_DOCKER_CONTAINER := goesquerydsl-container

#help:
#       make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[\:alnum\:]]' -e '^$@$$'
build:
	sudo docker build --no-cache -t sonet .
ps:
	sudo docker ps --all
up:
	sudo docker compose -f docker-compose.yml up -d
down:
	sudo docker compose -f docker-compose.yml down
restart:
	git pull --no-edit
	sudo docker compose -f docker-compose.yml down
	sudo docker compose -f docker-compose.yml up -d
logs:
	sudo docker compose -f docker-compose.yml logs --tail=100 -f sonet
clean:
	sudo docker image rm -f $(sudo docker image ls -aq)
cluster-up:
	mkdir -p /tmp/data_pg{1,2,3}
	sudo docker compose -f ./Cluster/docker-compose.yml up -d
cluster-down:
	sudo docker compose -f ./Cluster/docker-compose.yml down
update:
	git pull --no-edit
	sudo docker build --no-cache -t sonet .
	sudo docker compose down
	sudo docker compose up -d
	sudo docker container prune -f
