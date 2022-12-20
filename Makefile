BUILD_TIME := $(shell date --rfc-3339=seconds)
COMMIT_ID := $(shell git rev-parse HEAD)

LDFLAGS = -X "github.com/TeemoKill/BBot/utils.BuildTime='"$(BUILD_TIME)"'" -X "github.com/TeemoKill/BBot/utils.CommitId='"$(COMMIT_ID)"'"

SRC := $(shell find . -type f -name '*.go') static/*
COV := .coverage.out
TARGET := BBot

$(COV): $(SRC)
	go test ./... -coverprofile=$(COV)


$(TARGET): $(SRC) go.mod go.sum
	go build -ldflags '$(LDFLAGS)' -o $(TARGET) github.com/TeemoKill/BBot/main

build: $(TARGET)

test: $(COV)

coverage: $(COV)
	go tool cover -func=$(COV) | grep -v 'pb.go'

report: $(COV)
	go tool cover -html=$(COV)

clean:
	- rm -rf $(TARGET) $(COV)
