init:
	docker-compose build
	# create db
	docker-compose exec backend sqlite3 /root/sqlite3.db
	# migrate
	docker-compose exec migrate -path db/migration -database "sqlite3://root/sqlite3.db" -verbose up

start:
	docker-compose up
	
