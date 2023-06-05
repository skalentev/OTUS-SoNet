# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.18
GOLANG_DOCKER_CONTAINER := goesquerydsl-container

help:
	echo "use make <command>, commands: build, up, down, restart, update, clean, cluster-up, cluster-down, cluster-clean"
build:
	git pull --no-edit
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
	sudo docker logs --tail=100 -f sonet
clean:
	sudo docker image rm $(sudo docker image ls -aq)
cluster-clean:
	sudo docker volume rm $(sudo docker volume ls -q)
cluster-up:
	mkdir -p /tmp/data_pg1
	mkdir -p /tmp/data_pg2
	mkdir -p /tmp/data_pg3
	sudo docker compose -f ./Cluster/docker-compose.yml up -d
cluster-down:
	sudo docker compose -f ./Cluster/docker-compose.yml down
update:
	git pull --no-edit
	sudo docker build --no-cache -t sonet .
	sudo docker compose down
	sudo docker compose up -d
	sudo docker container prune -f
