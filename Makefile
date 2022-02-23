.DEFAULT_GOAL = build
.PHONY: FORCE
.DELETE_ON_ERROR:
.SECONDARY:

tools/ko = tools/bin/ko
tools/bin/%: tools/src/%/go.mod tools/src/%/pin.go
	cd $(<D) && GOOS= GOARCH= go build -o $(abspath $@) $$(sed -En 's,^import "(.*)".*,\1,p' pin.go)

build: $(tools/ko)
	$(tools/ko) build --local .
.PHONY: build

push: $(tools/ko)
	$(tools/ko) build .
.PHONY: push

speedtest.yaml: speedtest.in.yaml $(tools/ko)
	$(tools/ko) resolve --filename=$< > $@

apply: speedtest.in.yaml $(tools/ko)
	$(tools/ko) apply --filename=$<
.PHONY: apply
