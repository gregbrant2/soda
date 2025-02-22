
.PHONY : soda-dev
soda-dev:
	cd cmd/soda; air

.PHONY : soda-test
soda-test:
	go test ./...
	
.PHONY : start-soda-db
start-soda-db:
	docker run -p 3306:3306 --name soda_system_mysql -e MYSQL_ROOT_PASSWORD=password -d mysql;
	sleep 20;
	docker exec soda_system_mysql mysql -u root --password=password -Bse 'CREATE DATABASE soda;'

.PHONY : run-soda-db
run-soda-db:
	make start-soda-db
	make migrate-up

.PHONY : stop-soda-db
stop-soda-db:
	docker stop soda_system_mysql

.PHONY : rebuild-soda-db
rebuild-soda-db:
	docker stop soda_system_mysql
	docker rm soda_system_mysql
	make start-soda-db
	make migrate-up

.PHONY : migration
migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY : migrate-up
migrate-up:
	migrate -database 'mysql://root:password@/soda' -path db/migrations up $(count)

.PHONY : migrate-up-one
migrate-up-one:
	migrate -database 'mysql://root:password@/soda' -path db/migrations up 1

.PHONY : migrate-down
migrate-down:
	migrate -database 'mysql://root:password@/soda' -path db/migrations down $(count)

.PHONY : migrate-down-one
migrate-down-one:
	migrate -database 'mysql://root:password@/soda' -path db/migrations down 1

.PHONY : migrate-test-last
migrate-test-last:
	make migrate-up-one
	make migrate-down-one
	make migrate-up-one