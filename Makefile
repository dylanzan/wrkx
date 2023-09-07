.PHONY: install

all: help

install_wrk:
	@echo "Install wrk"
	bash -x ./install_wrk.sh

help:
	@echo "make install_wrk : 安装wrk"
