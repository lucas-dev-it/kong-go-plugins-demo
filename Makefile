prepare-and-run: prepare run

run:
	@echo '#### RUNNING ALL CONTAINERS ####'
	docker start kong-database
	sleep 5s
	docker start kong
	docker start login-api-demo
	docker start konga
	docker ps

prepare: setup-containers run-kong-db run-kong-db-migration run-kong run-konga run-login-demo stop-all

setup-containers:
	@echo '#### SETUP ALL CONTAINERS ####'
	docker build --no-cache -t lucas-dev-it/kong-go-plugins-demo _demo/login-api-demo
	docker pull pantsel/konga
	docker build -t kong-go-plugins-demo/custom-plugins .

run-kong-db:
	docker run -d --name kong-database \
		--network=kong-net \
		-p 5432:5432 \
		-e "POSTGRES_USER=kong" \
		-e "POSTGRES_DB=kong" \
		-e "POSTGRES_PASSWORD=kong" \
		postgres:9.6

run-kong-db-migration:
	sleep 5s
	docker run --rm --name kong \
		--network=kong-net \
		-e "KONG_DATABASE=postgres" \
		-e "KONG_PG_HOST=kong-database" \
		-e "KONG_PG_USER=kong" \
		-e "KONG_PG_PASSWORD=kong" \
		-e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
		kong-go-plugins-demo/custom-plugins:latest kong migrations bootstrap

run-kong:
	docker run -d --name kong \
		--network=kong-net \
		-e "KONG_DATABASE=postgres" \
		-e "KONG_PG_HOST=kong-database" \
		-e "KONG_PG_USER=kong" \
		-e "KONG_PG_PASSWORD=kong" \
		-e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
		-e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
		-e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
		-e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
		-e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
		-e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
		-e "KONG_GO_PLUGINS_DIR=/tmp/go-plugins" \
		-e "KONG_PLUGINS=bundled,example" \
		-e "KONG_LOG_LEVEL=debug" \
		-p 8000:8000 \
		-p 8443:8443 \
		-p 127.0.0.1:8001:8001 \
		-p 127.0.0.1:8444:8444 \
		kong-go-plugins-demo/custom-plugins:latest

run-konga:
	docker run -d -p 1337:1337 \
		--network kong-net \
		--name konga \
		pantsel/konga

run-login-demo:
	docker run -d -p 3333:3333 --name login-api-demo --network=kong-net lucas-dev-it/kong-go-plugins-demo:latest

stop-all:
	docker stop kong-database kong konga login-api-demo

setup-endpoints:
	chmod +x setup_endpoints.sh
	./setup_endpoints.sh
