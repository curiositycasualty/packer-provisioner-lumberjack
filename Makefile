SHELL = /bin/bash

PACKER_LOG_PATH := log/packer.log
PACKER_LOG ?= 1

.PHONIES: 

all: build install run

build:
	go build -o packer-provisioner-lumberjack

install:
	mkdir -pv ~/.packer.d/plugins/
	mv -v packer-provisioner-lumberjack ~/.packer.d/plugins/

run:
	rm -fv $(PACKER_LOG_PATH)
	PACKER_LOG=$(PACKER_LOG) PACKER_LOG_PATH='$(PACKER_LOG_PATH)' packer build test/template.json
	cat $(PACKER_LOG_PATH)
