DAY?=$(lastword $(sort $(wildcard day??)))
GO?=go

.PHONY: input sample
input sample:
	$(GO) build -o build/ ./$(DAY)
	./build/$(DAY) $(DAY)/$@.txt
