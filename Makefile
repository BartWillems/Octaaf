# octaaf Makefile
#
# This program is free software; you can redistribute
# it and/or modify it under the terms of the GNU
# General Public License …

SHELL = /bin/sh

srcdir = .

NAME 		= octaaf
DESCRIPTION = A Go Telegram bot
ARCH 		= x86_64

all: build

build:
	docker run --rm -v "$$PWD":/go/src/octaaf -w /go/src/octaaf golang:1.11 go build -ldflags "-s -w"

TMPDIR := $(shell mktemp -d)
TARGET := $(TMPDIR)/opt/octaaf
CONFIG := $(TARGET)/config
SYSTEM := $(TMPDIR)/usr/lib/systemd/system/

package:
	mkdir -p $(CONFIG)
	mkdir -p $(SYSTEM)

	cp ./octaaf $(TARGET)/
	cp ./octaaf.service $(SYSTEM)/octaaf.service
	cp ./config/settings.toml.dist $(CONFIG)/settings.toml
	cp ./config/database.yml $(CONFIG)/database.yml
	cp -r ./assets $(TARGET)/
	cp -r ./migrations $(TARGET)/
	
	fpm -s dir -t rpm \
		--name "$(NAME)" \
		--description "$(DESCRIPTION)" \
		--version "$(VERSION)" \
		--architecture "$(ARCH)" \
		--iteration $(BUILD_NUMBER) \
		--url 'https://gitlab.com/BartWillems/octaaf' \
		--license 'MIT' \
		--depends 'postgresql' \
		--depends 'redis' \
		--force \
		--config-files /usr/lib/systemd/system/octaaf.service \
		--config-files /opt/octaaf/config/settings.toml \
		--chdir $(TMPDIR) \
		.; \
	
	rm -R $(TMPDIR)

clean:
	rm -f octaaf*.rpm

.PHONY: clean
