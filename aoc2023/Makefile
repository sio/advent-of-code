SPAGO?=spago
SPAGO_ARGS?=--platform browser

ifneq (,$(RELEASE))
SPAGO_ARGS+=--minify
endif

ifneq (,$(WATCH))
SPAGO_ARGS+=--watch
endif

.PHONY: build
build: $(shell find -type f -name '*.purs')
	$(SPAGO) bundle-app $(SPAGO_ARGS)

.PHONY: serve
serve: build
	python3 -m http.server --bind 127.0.0.1
