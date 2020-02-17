# output directory, where all artifacts will be created and managed
OUTPUT_DIR ?= build/_output
# relative path to operator binary
OPERATOR = $(OUTPUT_DIR)/bin/operator
# golang cache directory path
GOCACHE ?= "$(shell echo ${PWD})/$(OUTPUT_DIR)/gocache"

default: build

.PHONY: build
build: $(OPERATOR)

$(OPERATOR): 
	$(Q)GOARCH=amd64 GOOS=linux go build -o $(OPERATOR) cmd/manager/main.go

.PHONY: local
local:
	- kubectl delete -f deploy/role.yaml
	- kubectl delete -f deploy/service_account.yaml
	- kubectl delete -f deploy/role_binding.yaml
	- kubectl delete -f deploy/operator.yaml

	kubectl apply -f deploy/role.yaml
	kubectl apply -f deploy/service_account.yaml
	kubectl apply -f deploy/role_binding.yaml
	kubectl apply -f deploy/operator.yaml

	operator-sdk run --local

clean:
	rm -rfv $(OUTPUT_DIR)

.PHONY: test
test: build
	$(Q)GOCACHE=$(GOCACHE) go test ./pkg/apis/... ./pkg/controller/...