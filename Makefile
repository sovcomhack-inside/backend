# path to actual config - the one that is copied to the docker container
CONFIG_PATH:=resources/config/config.yaml

# path to docker compose file
DCOMPOSE:=docker-compose.yaml

# path to external config which will copied to CONFIG_PATH
CONFIG_SOURCE_PATH=resources/config/config_default.yaml

# improve build time
DOCKER_BUILD_KIT:=COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1
DCOMPOSE_BUILD_ARGS:=--build-arg CONFIG_PATH=${CONFIG_PATH} --parallel

all: down build up

down:
	docker-compose -f ${DCOMPOSE} down --remove-orphans

build:
	cp ${CONFIG_SOURCE_PATH} ${CONFIG_PATH}
	${DOCKER_BUILD_KIT} docker-compose build ${DCOMPOSE_BUILD_ARGS}

up:
	docker-compose -f ${DCOMPOSE} up --remove-orphans

# Vendoring is useful for local debugging since you don't have to
# reinstall all packages again and again in docker
mod:
	go mod tidy && go mod vendor && go install ./...

acceptance:
	go test -v ./test/acceptance
