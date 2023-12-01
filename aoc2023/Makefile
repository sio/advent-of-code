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

.PHONY: messy-setup
messy-setup:
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
	bash -c "\
		source ~/.nvm/nvm.sh && \
		nvm install node && \
		npm install -g purescript spago esbuild"
