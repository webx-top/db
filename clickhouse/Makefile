SHELL := bash

DB_HOST ?= 127.0.0.1
DB_PORT	?= 9000

DB_USERNAME ?= upperio_tests
DB_PASSWORD ?= upperio_secret
DB_NAME     ?= upperio_tests

TEST_FLAGS ?=

export DB_HOST
export DB_NAME
export DB_PASSWORD
export DB_PORT
export DB_USERNAME

build:
	go build && go install

require-client:
	@if [ -z "$$(which clickhouse-client)" ]; then \
		echo 'Missing "clickhouse-client" command. Please install the clickhouse client and try again.' && \
		exit 1; \
	fi

generate:
	go generate && \
	go get -d -t -v ./...

reset-db: require-client
	SQL="" && \
	SQL+="DROP DATABASE IF EXISTS $(DB_NAME);" && \
	SQL+="CREATE DATABASE $(DB_NAME);" && \
	SQL+="GRANT ALL PRIVILEGES ON $(DB_NAME).* TO $(DB_USERNAME) IDENTIFIED BY '$(DB_PASSWORD)';" && \
	clickhouse-client --user=root --host="$(DB_HOST)" --port=$(DB_PORT) --database=test --multiquery <<< $$SQL

test: reset-db generate
	#go test -tags generated -v -race # race: limit on 8192 simultaneously alive goroutines is exceeded, dying
	go test -tags generated -v $(TEST_FLAGS)
