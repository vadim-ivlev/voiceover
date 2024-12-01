#!/bin/bash
#$@

# go run ./cmd/voiceover -s 18 -e 42 -translate Russian  $@
# go run ./cmd/voiceover -s 18 -e 42 -translate Russian  texts/dahl.epub 
# go run ./cmd/voiceover -s 122 -e 126 -translate Russian texts/The_Mind_is_flat.epub
go run ./cmd/voiceover -s 119 -e 126 -translate Russian texts/The_Mind_is_flat.epub

