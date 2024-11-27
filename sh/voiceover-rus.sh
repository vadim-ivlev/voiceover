#!/bin/bash

go run ./cmd/voiceover -s 90 -e 100 -voices nova -translateto Russian -o texts/rus texts/dahl.txt

