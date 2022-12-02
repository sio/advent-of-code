DAY?=$(lastword $(sort $(wildcard day??)))
GO?=go

.PHONY: input sample
input sample: build/$(DAY)  ## execute solution against input file or sample
	$< $(DAY)/$@.txt

build/$(DAY): $(shell find $(DAY) -type f)  ## build solution from source code
	$(GO) build -o $@ ./$(DAY)

.PHONY: all
all:  ## build and execute all solutions
	$(foreach d,$(wildcard day??),$(MAKE) DAY=$(d) &&) exit

.PHONY: stats
stats:  ## show stats based on git history
	@for d in $(sort $(wildcard day??)); \
	do export d; git log --format=format:%as -- $$d | sort -u | awk "$$STATS_AWK"; \
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
