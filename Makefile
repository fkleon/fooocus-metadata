GOFILES := $(shell find . -type f -name '*.go')
OUTFOLDER := ./out

.PHONY: all
all: $(OUTFOLDER)/print-metadata

$(OUTFOLDER)/print-metadata: ./cmd/file/main.go $(GOFILES)
	@go build -o $@ $<
	@chmod +x $@

template.png: template.txt
	@magick -size 240x85 -gravity center pango:@$< $@