# Definitions
.PHONY: help build ps up down restart logs clean cluster-clean cluster-upmaster cluster-upslave cluster-down cluster-drop cluster-import
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
	echo 'run  docker volume rm $(docker volume ls -q)'
cluster-upmaster:
#	mkdir -p /tmp/data-pg1
	sudo docker compose -f ./Cluster/docker-compose.yml create pg1
	sudo docker cp Cluster/Postgresql1.conf pg1:/var/lib/postgresql/data/postgresql.conf
	sudo docker cp Cluster/pg_hba.conf pg1:/var/lib/postgresql/data/pg_hba.conf
#	sudo cp db/user.csv.gz /tmp/data-pg1/user.csv.gz
	sudo docker compose -f ./Cluster/docker-compose.yml run -d pg1
	sudo docker compose -f ./Cluster/docker-compose.yml up -d cadvisor node-exporter sonet
cluster-import:
	sudo docker cp db/user.csv.gz pg1:/var/lib/postgresql/data/user.csv.gz
	sudo docker exec pg1 psql -d cluster -U user -c "copy public.user from program 'gzip -d -c /var/lib/postgresql/data/user.csv.gz' (FORMAT CSV, HEADER);"
	sudo docker exec pg1 rm /var/lib/postgresql/data/user.csv.gz
cluster-upslave:
	sudo mkdir -p /tmp/data-pg2
	sudo mkdir -p /tmp/data-pg3
	sudo docker exec pg1 mkdir -p /pgslave
	sudo docker exec -e PGPASSWORD='pass' pg1 pg_basebackup -h pg1 -D /pgslave -U replicator -v -P --wal-method=stream
	sudo docker cp pg1:/pgslave /tmp/pgslave
	sudo cp /tmp/pgslave /tmp/data-pg2
	sudo cp Cluster/Postgresql2.conf /tmp/data-pg2/postgresql.conf
	sudo cp Cluster/pg_hba.conf /tmp/data-pg2/pg_hba.conf
	sudo cp Cluster/standby.signal /tmp/data-pg2/standby.signal
	sudo cp /tmp/pgslave /tmp/data-pg3
	sudo cp Cluster/Postgresql3.conf /tmp/data-pg3/postgresql.conf
	sudo cp Cluster/pg_hba.conf /tmp/data-pg3/pg_hba.conf
	sudo cp Cluster/standby.signal /tmp/data-pg3/standby.signal
	sudo docker compose -f ./Cluster/docker-compose.yml up -d pg2 pg3
cluster-down:
	sudo docker compose -f ./Cluster/docker-compose.yml down
cluster-drop:
	sudo docker compose -f ./Cluster/docker-compose.yml down -v
	sudo rm -r /tmp/data-pg1
	sudo rm -r /tmp/data-pg2
	sudo rm -r /tmp/data-pg3
	sudo rm -r /tmp/pgslave
psql1:
	sudo docker exec -ti pg1 psql -d cluster -U user
psql2:
	sudo docker exec -ti pg2 psql -d cluster -U user
psql3:
	sudo docker exec -ti pg3 psql -d cluster -U user
pg1:
	sudo docker exec -ti pg1 bash
pg2:
	sudo docker exec -ti pg2 bash
pg3:
	sudo docker exec -ti pg3 bash

update:
	git pull --no-edit
	sudo docker build --no-cache -t sonet .
	sudo docker compose down
	sudo docker compose up -d
	sudo docker container prune -f
