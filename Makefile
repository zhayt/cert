build:
	docker build -t cert-app .

up:
	docker-compose up --build

down:
	docker-compose stop

run:
	./run.sh