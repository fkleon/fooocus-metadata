GOFILES := $(shell find . -type f -name '*.go')
GOFLAGS := -ldflags="-s -w" -trimpath
OUTFOLDER := ./out

.PHONY: cmd
cmd: $(OUTFOLDER)/read-metadata $(OUTFOLDER)/write-metadata

$(OUTFOLDER)/read-metadata: ./cmd/extract/main.go $(GOFILES)
	@go build ${GOFLAGS} -o $@ $<
	@chmod +x $@

$(OUTFOLDER)/write-metadata: ./cmd/embed/main.go $(GOFILES)
	@go build ${GOFLAGS} -o $@ $<
	@chmod +x $@

.PHONY: test
test:
	@go test ./...

types/template.png: types/template.txt
	@magick -size 240x85 -gravity center pango:@$< $@
