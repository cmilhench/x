default: help
.PHONY: default

name ?= $(notdir $(realpath $(dir $(realpath $(MAKEFILE_LIST)/..))))
host ?= postgresql://postgres:password@localhost

help: # This
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

# -----------------------------------------------------------------------------

pg-start: ## Start a development postgresql database server
	$(call check-dependency,docker)
	@docker run --rm --detach \
					--volume "${PWD}/.data":/var/lib/postgresql/data \
					--name pg-$(name) \
					--env POSTGRES_PASSWORD=password \
					--publish 5432:5432 \
	  			postgres
.PHONY: pg-start

pg-stop: ## Stop a development postgresql database server
	$(call check-dependency,docker)
	@docker stop pg-$(name)
.PHONY: pg-stop

$(wildcard *.sql):
	$(call check-dependency,docker)
	@docker run --rm -it \
					--volume "${PWD}":/tmp/schema \
					--net=host \
					postgres psql $(host) -f /tmp/schema/$@
.PHONY: $(wildcard *.sql)

# -----------------------------------------------------------------------------

es-start: ## Start a development elasticsearch database server
	$(call check-dependency,docker)
	@-mkdir ${PWD}/../.data/es
	@docker run --rm --detach --interactive --tty \
					--network es-net \
					--memory 4GB \
					--mount type=bind,src=${PWD}/../.data/es,dst=/usr/share/elasticsearch/data \
					--name es-$(name) \
					--publish 9200:9200 \
					-e "discovery.type=single-node" \
					-e "ES_JAVA_OPTS=-Xms4g -Xmx4g" \
					-e "xpack.security.enabled=false" \
					--env ELASTIC_USERNAME=elastic \
          --env ELASTIC_PASSWORD=DkIedPPSCb \
					docker.elastic.co/elasticsearch/elasticsearch:8.11.1
.PHONY: es-start

es-cert:
	$(call check-dependency,docker)
	@docker cp es-$(name):/usr/share/elasticsearch/config/certs/http_ca.crt ${PWD}/../.data/es.crt 1>/dev/null
	@echo curl --cacert ${PWD}/../.data/es.crt -u elastic:DkIedPPSCb https://localhost:9200
.PHONY: es-cert

es-token:
	$(call check-dependency,docker)
	@docker exec --interactive --tty es-$(name) /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token --scope kibana
.PHONY: es-token

es-stop: ## Stop a development elasticsearch database server
	$(call check-dependency,docker)
	@docker stop es-$(name)
.PHONY: es-stop

# -----------------------------------------------------------------------------

kib-start: ## Start a development kibana server
	@docker run --rm --detach --interactive --tty \
					--network es-net \
					--name kib-$(name) \
					--publish 5601:5601 \
					docker.elastic.co/kibana/kibana:8.11.1
.PHONY: kib-start

kib-code:
	@docker exec --interactive --tty kib-$(name) /usr/share/kibana/bin/kibana-verification-code
.PHONY: kib-code

kib-stop: ## Stop a development kibana server
	@docker stop kib-$(name)
.PHONY: kib-stop

# -----------------------------------------------------------------------------

define check-dependency
	$(if $(shell command -v $(1)),,$(error Make sure $(1) is installed))
endef
