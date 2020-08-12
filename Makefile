SERVICE_NAME	= snmp-trap-generator
DIST_DIR 	= ./dist
BIN_DIR		= bin
WINDOWS32_BIN 	= $(DIST_DIR)/win-x86/$(SERVICE_NAME).exe
WINDOWS64_BIN 	= $(DIST_DIR)/win-x64/$(SERVICE_NAME).exe
LINUX32_BIN		= $(DIST_DIR)/linux-32/$(SERVICE_NAME)
LINUX64_BIN		= $(DIST_DIR)/linux-64/$(SERVICE_NAME)

ifeq ($(OS), Windows_NT)
TARGETD	= $(BIN_DIR)/$(SERVICE_NAME).exe
else
TARGETD	= $(BIN_DIR)/$(SERVICE_NAME)
endif

.PHONY: all windows linux dist clean dist

all:
	rm -rf $(BIN_DIR)
	go build -o $(TARGETD)

$(DIST_DIR)/:
	mkdir -p $@

$(WINDOWS32_BIN): $(DIST_DIR)/
	GOOS=windows GOARCH=386 go build -o $@ .

$(WINDOWS64_BIN): $(DIST_DIR)/
	GOOS=windows GOARCH=amd64 go build -o $@ .

$(LINUX32_BIN): $(DIST_DIR)/
	GOOS=windows GOARCH=386 go build -o $@ .

$(LINUX64_BIN): $(DIST_DIR)/
	GOOS=windows GOARCH=amd64 go build -o $@ .

windows: $(WINDOWS32_BIN) $(WINDOWS64_BIN)
	@echo "===> Built $^"

linux: $(LINUX32_BIN) $(LINUX64_BIN)
	@echo "===> Built $<"

dist: clean windows linux

clean:
	rm -rf $(DIST_DIR)
