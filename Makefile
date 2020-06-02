build: go.sum
ifeq ($(OS),Windows_NT)
	@go build -mod=readonly $(BUILD_FLAGS) -o build/gifTyper.exe main.go
else
	@go build -mod=readonly $(BUILD_FLAGS) -o build/gifTyper main.go
endif

install: go.sum
	@go build -mod=readonly $(BUILD_FLAGS) -o $${GOBIN-$${GOPATH-$$HOME/go}/bin}/gifTyper main.go