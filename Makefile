GOBIN := $(shell go env GOPATH)/bin
MOCKGEN := $(GOBIN)/mockgen
SWAG := $(GOBIN)/swag

.PHONY: test mocks swag

test:
	go test ./...

mocks: $(MOCKGEN)
	$(MOCKGEN) -source=store/service.go -destination=mocks/mock_store.go -package=mocks
	$(MOCKGEN) -source=internal/services/in_app_store/service.go -destination=mocks/mock_in_app_store_service.go -package=mocks

swag: $(SWAG)
	$(SWAG) init -g cmd/adora-coding-assessment/main.go

$(MOCKGEN):
	go install go.uber.org/mock/mockgen@latest

$(SWAG):
	go install github.com/swaggo/swag/cmd/swag@latest
