#!/usr/bin/env sh

protoc --go_out=../.. --go_opt=module=github.com/elliotmr/plug plug.proto