PLATFORMS := linux/arm linux/arm64 linux/amd64 windows/arm windows/arm64 windows/amd64 darwin/arm64 darwin/amd64
BIN_DIR := bin
SRC_DIR := src

compile:
	echo "========= Compiling for every OS and Platform ========="
	mkdir -p bin
	$(foreach platform,	$(PLATFORMS), \
		$(eval GOOS=$(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH=$(word 2,$(subst /, ,$(platform)))) \
		go build -C ./$(SRC_DIR) -o ../$(BIN_DIR)/deraph-$(GOOS)-$(GOARCH)$(if $(findstring windows,$(GOOS)),.exe) main.go;)
