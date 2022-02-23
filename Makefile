.DEFAULT_GOAL = build

tools/ko = tools/bin/ko
tools/bin/%: tools/src/%/go.mod tools/src/%/pin.go
	cd $(<D) && GOOS= GOARCH= go build -o $(abspath $@) $$(sed -En 's,^import "(.*)".*,\1,p' pin.go)

build: $(tools/ko)
	$(tools/ko) build --local .
.PHONY: build
