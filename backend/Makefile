.PHONY: test start mysql.start mysql.load mysql.stop

test:
	go test -v ...

start:
	go build -o app
	./app

mysql.start:
	docker run --rm -d \
		-e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} \
		-e MYSQL_USER=${DB_USER} \
		-e MYSQL_PASSWORD=${DB_PASSWORD} \
		-e MYSQL_DATABASE=${DB_NAME} \
		-e TZ="Asia/Tokyo" \
		-h ${DB_HOST} \
		-p ${DB_PORT}:3306 \
		--name ${DB_HOST} mysql:5.7

mysql.load:
	mysql -h 127.0.0.1 --port ${DB_PORT} \
		-u${DB_USER} -p${DB_PASSWORD} \
		${DB_NAME} < mysql/init.sql

mysql.stop:
	docker stop ${DB_HOST}
