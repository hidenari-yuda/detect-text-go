OS=linux
ARCH=amd64
ROOT=$(GOPATH)/src/github.com/hidenari-yuda/umerun-resume

setup:
	go install -v github.com/google/wire/cmd/wire@v0.5.0
	go install -v github.com/rubenv/sql-migrate/sql-migrate@v1.1.2
	go install -v github.com/golang/mock/mockgen@v1.6.0
	go install -v github.com/cosmtrek/air@v1.40.4
	make mock

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app

wire:
	wire ./infrastructure/di/wire.go

run:
	export APP_ENV=local
	air -c .conf/.air.toml

up:
	docker-compose -f docker-compose.yml up --build -d

down:
	docker-compose -f docker-compose.yml down --volumes

enc-envfile:
	cat .env | base64

migrate-new:
	sql-migrate new -env=local -config=.conf/dbconfig.yml ${FILE}

migrate-up:
	sql-migrate up -config=.conf/dbconfig.yml -env=local

migrate-down:
	sql-migrate down -config=.conf/dbconfig.yml -env=local


## Test Command
test:
	make down-test-db
	make up-test-db
	make test-all
	make down-test-db

test-all:
	make mock
	make repository-test
	# make entity-test

repository-test:
	DB_NAME=umerun_resume_test go test -v ./tests/repository/*_test.go

#entity-test:
	#go test -v ./tests/entity/*_test.go

## Mock
mock: repository-mock interactor-mock driver-mock # usecase-mock

repository-mock:
	mockgen -source ./usecase/repository.go -destination ./mock/mock_usecase/mock_repository.go

interactor-mock: 
	basename -a ./usecase/interactor/*.go | sed 's/.go//' | xargs -IFILE mockgen -source ./usecase/interactor/FILE.go -destination ./mock/mock_interactor/mock_FILE.go
	rm ./mock/mock_interactor/mock_wire_set.go

driver-mock: 
	mockgen -source ./usecase/driver.go -destination ./mock/mock_usecase/mock_driver.go

usecase-mock:
	mockgen -source ./usecase/usecase.go -destination ./mock/mock_usecase/mock_usecase.go