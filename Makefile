PACKAGE=willdo
.DEFAULT_GOAL := build
INSTALL_DIR   ::= /usr/local/bin

.PHONY:fmt vet build
fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o ./build/$(PACKAGE)

install: build
	sudo cp ./build/$(PACKAGE) $(INSTALL_DIR)/

clean: 
	go clean
	rm ./build/*