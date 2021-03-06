MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`
BIN_PATH = "./bin"
BIN_NAME = "mystack"

setup: setup-hooks
	@go get -u github.com/golang/dep...
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/wadey/gocovmerge
	@dep ensure -update

setup-hooks:
	@cd .git/hooks && ln -sf ../../hooks/pre-commit.sh pre-commit

build:
	@go build -o ./bin/mystack main.go

build-all-platforms:
	@mkdir -p ${BIN_PATH}
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ${BIN_PATH}/${BIN_NAME}-linux-i386
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-linux-x86_64
	@echo "Building for darwin-i386..."
	@env GOOS=darwin GOARCH=386 go build -o ${BIN_PATH}/${BIN_NAME}-darwin-i386
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-darwin-x86_64
	@echo "Building for win-x86_64..."
	@env GOOS=windows GOARCH=amd64 go build -o ${BIN_PATH}/${BIN_NAME}-win-x86_64

unit: clear-coverage-profiles unit-run gather-unit-profiles

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit-run:
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

run:
	@go run main.go

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

integration int: clear-coverage-profiles integration-run gather-integration-profiles

integration-run:
	@ginkgo -tags integration -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-integration-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'

merge-profiles:
	@mkdir -p _build
	@gocovmerge _build/*.out > _build/coverage-all.out

test-coverage-func coverage-func: merge-profiles
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

test: unit integration test-coverage-func
