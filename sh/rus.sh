#!/bin/bash
#$@

# go run ./cmd/voiceover -start 18 -end 42 -translate Russian  $@
# go run ./cmd/voiceover -start 18 -end 42 -translate Russian  texts/dahl.epub 
# go run ./cmd/voiceover -start 122 -end 126 -translate Russian texts/The_Mind_is_flat.epub
go run ./cmd/voiceover -start 119 -end 126 -translate Russian texts/The_Mind_is_flat.epub

