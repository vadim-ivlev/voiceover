#!/bin/bash

# Sam
go run ./cmd/voiceover -start 90 -end 100 -ttsapi elevenlabs -voices ulNeoiyl3bUW7oQjWZE8 -translate Russian -output texts/rus-sam texts/dahl.txt

# Valentino
# go run ./cmd/voiceover -start 90 -end 100 -ttsapi elevenlabs -voices HgJDD5cRFQsVhwzXouaI -translate Russian -output texts/rus texts/dahl.txt

# Adam Stone
# go run ./cmd/voiceover -start 90 -end 100 -ttsapi elevenlabs -voices NFG5qt843uXKj4pFvR7C -translate Russian -output texts/rus-adam texts/dahl.txt