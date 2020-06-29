#!/usr/bin/env sh

protoc -I=. -I=$(go list -f '{{.Dir}}' github.com/elliotmr/plug/pkg/plugpb) --go_out=.. --go_opt=module=github.com/elliotmr/plug --plug_out=.. --plug_opt=module=github.com/elliotmr/plug example.proto
