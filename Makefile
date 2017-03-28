MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`
MYSTACK_LOGIN_TEST='test'

setup: setup-hooks
	@go get -u github.com/golang/dep...
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/wadey/gocovmerge
	@dep ensure -update

setup-hooks:
	@cd .git/hooks && ln -sf ../../hooks/pre-commit.sh pre-commit

build:
	@go build -o ./bin/kubecos-cli main.go

export-test-var:
	@export MYSTACK_LOGIN_TEST

unexport-test-var:
	@unset MYSTACK_LOGIN_TEST

unit: export-test-var clear-coverage-profiles unit-run gather-unit-profiles

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit-run:
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

integration int: export-test-var clear-coverage-profiles integration-run gather-integration-profiles

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
