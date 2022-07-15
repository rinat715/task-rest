build:
	docker-compose build

init:
	docker-compose up -d
	# migrate
	docker-compose exec backend /root/go_rest --migratePath "file://./migration" 

start:
	docker-compose up
	
