DAY?=$(lastword $(sort $(wildcard day??)))
GO?=go

.PHONY: input sample
input sample: build/$(DAY)
	$< $(DAY)/$@.txt

build/$(DAY): $(shell find $(DAY) -type f)
	$(GO) build -o $@ ./$(DAY)

.PHONY: all
all:
	$(foreach d,$(wildcard day??),$(MAKE) DAY=$(d) &&) exit

.PHONY: stats
stats:
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
