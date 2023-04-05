# Makefile

YAML_FILES := $(wildcard *.yaml)
PDF_FILES := $(patsubst %.yaml,%.pdf,$(YAML_FILES))

all: $(PDF_FILES)

%.pdf: %.yaml
	../wordsearch -d ../dict-en.txt -i $<

clean:
	rm -f $(PDF_FILES)

.PHONY: all clean

