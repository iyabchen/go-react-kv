DOCKER_FILE=Dockerfile 
IMAGE_NAME=demo
IMAGE_VER?=1.0
ifdef DOCKER_BUILD_NO_CACHE
	DOCKER_BUILD_OPT=--no-cache
endif

PKG=github.com/iyabchen/go-react-kv/server
SRV_PKG=$(PKG)/cmd
SRV_DIR=$(PWD)/server
SRV_BIN=srv
SRV_BIN_DIR=$(SRV_DIR)/bin

TESTARGS=-short  # -v 

 
all:
	@echo "make <target>"
	@echo ""
	@echo "targets:"
	@echo "  build           - build binaries in the build/ directory"
	@echo "  test            - runs tests"
	@echo "  docker          - build a dockerimage of the server service"
	@echo "  clean           - remove intermediate files"


docker: 
	docker build $(DOCKER_BUILD_OPT) \
		-t $(IMAGE_NAME):$(IMAGE_VER) \
		-f $(DOCKER_FILE) \
		.
 
.deps: 
	-go get -u github.com/golang/dep/cmd/dep
	cd $(SRV_DIR) && dep ensure -v
	@touch $(SRV_DIR)/.deps

build: .deps test
	go build -o $(SRV_BIN) $(SRV_PKG)
	strip -s $(SRV_BIN)
	@mkdir -p $(SRV_BIN_DIR)
	mv $(SRV_BIN) $(SRV_BIN_DIR)/

clean:
	@rm -f $(SRV_DIR)/test.out
	@rm -f $(SRV_DIR)/coverage_tmp.out
	@rm -f $(SRV_DIR)/coverage.xml
	@rm -rf $(SRV_DIR)/Gopkg.lock
	@rm -rf $(SRV_DIR)/.deps
	@rm -rf $(SRV_DIR)/vendor/
	@rm -rf $(GOPATH)/pkg/dep/sources	
	@rm -f $(SRV_BIN)  
	@rm -rf $(SRV_BIN_DIR)
 
test: .deps  
	-go get -u github.com/axw/gocov/...
	-go get -u github.com/AlekSi/gocov-xml
	@rm -f $(SRV_DIR)/test.out
	@rm -f $(SRV_DIR)/coverage_tmp.out
	@rm -f $(SRV_DIR)/coverage.xml
	gocov test $(TESTARGS) $(SRV_DIR)/... 2>&1 1>$(SRV_DIR)/coverage_tmp.out | tee -a $(SRV_DIR)/test.out
	@grep "FAIL" $(SRV_DIR)/test.out > /dev/null; \
	if [ $$? -eq 0 ] ; then \
		exit 1; \
	fi
	gocov-xml > $(SRV_DIR)/coverage.xml < $(SRV_DIR)/coverage_tmp.out && rm -f $(SRV_DIR)/coverage_tmp.out


.PHONY: all docker build clean test 
