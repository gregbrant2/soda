
.PHONY : run
run:
	cd src;
	pwd;
	go run soda.go

.PHONY : start-soda-db
start-soda-db:
	docker run -p 3306:3306 --name soda_system_mysql -e MYSQL_ROOT_PASSWORD=password -d mysql;
	sleep 15;
	docker exec soda_system_mysql mysql -u root --password=password -Bse 'CREATE DATABASE soda;'

.PHONY : run-soda-db
run-soda-db:
	make start-soda-db
	make migrate-up

.PHONY : migration
migration:
	migrate create -ext sql -dir src/migrations -seq $(name)

.PHONY : migrate-up
migrate-up:
	migrate -database 'mysql://root:password@/soda' -path src/migrations up $(count)

.PHONY : migrate-up-one
migrate-up-one:
	migrate -database 'mysql://root:password@/soda' -path src/migrations up 1

.PHONY : migrate-down
migrate-down:
	migrate -database 'mysql://root:password@/soda' -path src/migrations down $(count)

.PHONY : migrate-down-one
migrate-down-one:
	migrate -database 'mysql://root:password@/soda' -path src/migrations down 1

.PHONY : migrate-test-last
migrate-test-last:
	make migrate-up-one
	make migrate-down-one
	make migrate-up-one