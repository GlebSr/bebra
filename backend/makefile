SHELL := /bin/bash
SUBDIRS := $(shell find . -mindepth 2 -type f -name 'makefile' -exec dirname {} \; 2>/dev/null | sort | uniq)
export MIGRATIONS_DIR="$(shell pwd)/migrations"

GENERATE_SQL_SH ?= $(realpath utils/generate-sql.sh)
export GENERATE_SQL_SH

build: cmd/server.go
	go build cmd/server.go

generate:
	@if [ -n "$(SUBDIRS)" ]; then \
		for dir in $(SUBDIRS); do \
			echo "Running 'make generate' in $$dir"; \
			$(MAKE) -C $$dir generate || exit 1; \
		done; \
	else \
		echo "No subdirectories with Makefile found"; \
	fi

generate_config:
	@set -a && . .env && set +a && \
	echo "Running 'make generate' in config/" && \
	$(MAKE) -C config/ generate || exit 1

clean:
	@if [ -n "$(SUBDIRS)" ]; then \
		for dir in $(SUBDIRS); do \
			echo "Running 'make clean' in $$dir"; \
			$(MAKE) -C $$dir clean || exit 1; \
		done; \
	else \
		echo "No subdirectories with Makefile found"; \
	fi