#!/bin/bash

go run ./cmd/voiceover -s 90 -e 100 -ttsapi elevenlabs -voices ulNeoiyl3bUW7oQjWZE8 -translateto Russian -o texts/rus texts/dahl.txt

