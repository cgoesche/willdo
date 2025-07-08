PACKAGE=willdo
PACKAGE_VERSION=$(shell grep "Version" app/app.go | cut -d"=" -f 2 | tr -d '" ')
.DEFAULT_GOAL := build
INSTALL_DIR ?= /usr/local/bin

MAN_DATABASE ?= /usr/share/man/man1
MAN_DIR = resources/man
MANPAGE_GZ = $(MAN_DIR)/$(PACKAGE).1.gz

.PHONY:fmt vet build manpage
fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

build: vet manpage
	go build -o ./build/$(PACKAGE)

manpage: 
	asciidoctor -b manpage -a release-version=$(PACKAGE_VERSION) -a adjtime_path=/etc/adjtime $(MANPAGE_GZ:.gz=.adoc)
	gzip --suffix=.gz -f $(MAN_DIR)/$(PACKAGE).1

install: build
	sudo cp ./build/$(PACKAGE) $(INSTALL_DIR)/
	sudo cp $(MANPAGE_GZ) $(MAN_DATABASE)

clean: 
	go clean
	rm ./build/*