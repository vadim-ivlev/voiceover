#!/bin/bash

# Sam
go run ./cmd/voiceover -s 90 -e 100 -ttsapi elevenlabs -voices ulNeoiyl3bUW7oQjWZE8 -translate Russian -o texts/rus-sam texts/dahl.txt

# Valentino
# go run ./cmd/voiceover -s 90 -e 100 -ttsapi elevenlabs -voices HgJDD5cRFQsVhwzXouaI -translate Russian -o texts/rus texts/dahl.txt

# Adam Stone
# go run ./cmd/voiceover -s 90 -e 100 -ttsapi elevenlabs -voices NFG5qt843uXKj4pFvR7C -translate Russian -o texts/rus-adam texts/dahl.txt