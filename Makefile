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
	sudo docker compose -f ./Cluster/docker-compose.yml up -d pg1
	sleep 5
	sudo docker cp Cluster/Postgresql1.conf pg1:/var/lib/postgresql/data/postgresql.conf
	sudo docker cp Cluster/pg_hba.conf pg1:/var/lib/postgresql/data/pg_hba.conf
	sudo docker compose -f ./Cluster/docker-compose.yml restart pg1
	sudo docker compose -f ./Cluster/docker-compose.yml up -d cadvisor node-exporter sonet
cluster-setsync:
	sudo docker cp Cluster/Postgresql1_any.conf pg1:/var/lib/postgresql/data/postgresql.conf
	sudo docker compose -f ./Cluster/docker-compose.yml restart pg1
cluster-import:
	sudo docker cp db/user.csv.gz pg1:/var/lib/postgresql/data/user.csv.gz
	sudo docker exec pg1 psql -d cluster -U user -c "copy public.user from program 'gzip -d -c /var/lib/postgresql/data/user.csv.gz' (FORMAT CSV, HEADER);"
	sudo docker exec pg1 rm /var/lib/postgresql/data/user.csv.gz
cluster-upslave:
	sudo rm -rf /tmp/data_pg2
	sudo rm -rf /tmp/data_pg3
	sudo docker exec pg1 rm -rf /pgslave
	sudo docker exec pg1 mkdir -p /pgslave;
	sudo docker exec -e PGPASSWORD='pass' pg1 pg_basebackup -h pg1 -D /pgslave -U replicator -v -P --wal-method=stream
	sudo docker cp pg1:/pgslave /tmp/pgslave
	sudo cp -r /tmp/pgslave/ /tmp/data_pg2/
	sudo cp Cluster/Postgresql2.conf /tmp/data_pg2/postgresql.conf
	sudo cp Cluster/pg_hba.conf /tmp/data_pg2/pg_hba.conf
	sudo cp Cluster/standby.signal /tmp/data_pg2/standby.signal
	sudo cp -r /tmp/pgslave/ /tmp/data_pg3/
	sudo cp Cluster/Postgresql3.conf /tmp/data_pg3/postgresql.conf
	sudo cp Cluster/pg_hba.conf /tmp/data_pg3/pg_hba.conf
	sudo cp Cluster/standby.signal /tmp/data_pg3/standby.signal
	sudo docker compose -f ./Cluster/docker-compose.yml up -d pg2 pg3
cluster-upredis:
	sudo docker compose -f ./Cluster/docker-compose.yml up -d redis
cluster-downslave:
	sudo docker stop pg2 pg3
	sudo docker container prune -f
cluster-update:
	git pull --no-edit
	sudo docker build --no-cache -t sonet .
	sudo docker stop sonet || true
	sudo docker container prune -f
	sudo docker compose -f ./Cluster/docker-compose.yml up -d sonet
cluster-down:
	sudo docker compose -f ./Cluster/docker-compose.yml down
cluster-drop:
	sudo docker compose -f ./Cluster/docker-compose.yml down -v
	sudo rm -r /tmp/data_pg1 || true
	sudo rm -r /tmp/data_pg2 || true
	sudo rm -r /tmp/data_pg3 || true
	sudo rm -r /tmp/pgslave || true
	sudo rm -r /tmp/redis || true
redis:
	sudo docker exec -it redis sh
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
