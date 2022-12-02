DAY?=$(lastword $(sort $(wildcard day??)))
GO?=go

.PHONY: input sample
input sample: build/$(DAY)
	$< $(DAY)/$@.txt

build/$(DAY): $(shell find $(DAY) -type f)
	$(GO) build -o $@ ./$(DAY)
