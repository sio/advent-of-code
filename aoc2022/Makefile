DAY?=$(lastword $(sort $(wildcard day*)))
GO?=go

ifdef PROFILE
ARGS=-cpuprofile $(DIRECTORY)/cpu.prof -memprofile $(DIRECTORY)/mem.prof
endif

DIRECTORY=$(DAY)
ifeq (,$(wildcard $(DIRECTORY)))
DIRECTORY:=day$(DAY)
endif
ifeq (,$(wildcard $(DIRECTORY)))
DIRECTORY:=day0$(DAY)
endif
ifeq (,$(wildcard $(DIRECTORY)))
$(error Directory not found: $(DAY))
endif

.PHONY: solve
solve:  ## run solution for the latest puzzle
solve: input

DAYS=$(wildcard day*)
DAYS+=$(subst day0,day,$(wildcard day*))
DAYS:=$(sort $(DAYS))
.PHONY: $(DAYS)
$(DAYS):
	@$(MAKE) DAY=$(subst day,,$@)

.PHONY: input sample
input sample sample2: build/$(DIRECTORY)  ## execute solution against input file or sample
	$< $(DIRECTORY)/$@.txt $(ARGS)

.PHONY: build
build: build/$(DIRECTORY)  ## build solution from source code
build/$(DIRECTORY): $(shell find $(DIRECTORY) -type f)
	$(GO) build -o $@ ./$(DIRECTORY)

.PHONY: test
test:  ## run automated tests for current day
	$(GO) test ./$(DIRECTORY) -timeout 15s $(GOTEST_ARGS)

.PHONY: test-verbose
export GOTEST_ARGS
test-verbose: GOTEST_ARGS+=-v --count=1
test-verbose:  ## run tests with increased verbosity
	$(MAKE) test

.PHONY: bench
bench:  ## run benchmarks for current day
	cd $(DIRECTORY) && go test -bench=. -count=3 -benchmem -benchtime=2s -run='^#'

.PHONY: all
all:  ## build and execute all solutions
	$(foreach d,$(wildcard day??),$(MAKE) DAY=$(d) &&) exit

.PHONY: answer
answer:  ## show answers for current day
	grep --color=auto -R 'Your puzzle answer was' $(DIRECTORY)

.PHONY: answers
answers:  ## show answers for all days
answers:  DIRECTORY=day*/
answers:  answer

.PHONY: fmt
fmt:  ## format Go code
	$(GO) fmt ./...

STATICCHECK=$(shell $(GO) env GOPATH)/bin/staticcheck$(shell $(GO) env GOEXE)
$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

ERRCHECK=$(shell $(GO) env GOPATH)/bin/errcheck$(shell $(GO) env GOEXE)
$(ERRCHECK):
	$(GO) install github.com/kisielk/errcheck@latest

.PHONY: lint
lint: $(STATICCHECK) $(ERRCHECK)
lint:  ## run linter for current day's solution
	$(GO) fmt ./$(DIRECTORY)
	$(GO) vet ./$(DIRECTORY)
	$(STATICCHECK) ./$(DIRECTORY)
	$(ERRCHECK) ./$(DIRECTORY)

.PHONY: stats
stats:  ## show stats based on git history
	@for d in $(sort $(wildcard day??)); \
	do export d; git log --format=format:%as -- $$d ':!**README*' | sort -u | awk "$$STATS_AWK"; \
	done

define STATS_AWK
{
	days += 1;
	details = details ((length(details)==0) ? "" : ",") $$0;
}
END {
	edited = sprintf("%d day%s", days, (days>1) ? "s" : " ");
	printf("%s: edited for %s (%s)\n", ENVIRON["d"], edited, details);
}
endef
export STATS_AWK

.PHONY: help
help:  ## show help message
	@awk "$$HELP_AWK" $(lastword $(MAKEFILE_LIST)) | column -t -s '	'

define HELP_AWK
/  ##/ {
	if (stop) exit;

	match($$0, /^[^:]+:/);
	target = substr($$0, RSTART, RLENGTH-1);
	gsub(" ", ",", target);

	match($$0, /  ## /);
	help = substr($$0, RSTART + RLENGTH);

	printf("%s\t%s\n", target, help);
}
/HELP_AWK/ {
	stop = 1;
}
endef
export HELP_AWK
