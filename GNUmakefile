RESOURSES?=user,group,app,auth_server,event_hook,group_rule,trusted_origin,user_type,template_sms,inline_hook,idp,policy,network_zone
LOGFILE= $(shell date +'%Y%m%d_%H%M%S').log

default: build

dep: # Download required dependencies
	go mod tidy
	# go mod vendor

build: cleanup
	@echo "==> Buiding binary"
	go build -v

cleanup:
	@echo "==> Cleanup"
	rm -rf ./generated
	rm -f ./terraformer

cleanup-log:
	@echo "Delete log files"
	@find . -type f -name '*.log*' -exec rm -f {} +

sample:
	@echo "Manual testing"
	@echo "Cleanup previous generated output"
	rm -rf ./generated
	@echo "Generated output"
	@echo $(LOGFILE)
	# Override: make sample RESOURCES=user
	./terraformer import okta --resources=$(RESOURSES) --verbose > $(LOGFILE)

test:
	go test ./...

vet:
	go vet

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test cleanup cleanup-log