
DOCKER ?= docker

run-couchdb:
	$(DOCKER) run -p 5984:5984 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=admin -d couchdb:3.1.1
