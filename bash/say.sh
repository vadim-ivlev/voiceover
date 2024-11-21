#!/bin/bash

#https://serverfault.com/questions/72476/clean-way-to-write-complex-multi-line-string-to-a-variable



# params=$(cat <<EOF
# {
#   "model": "tts-1",
#   "input": "$1",
#   "voice": "echo", 
#   "speed": 1.0
# }
# EOF
# )
  # "language": "en",
# echo "-d '$params'"

# Experiment with different voices (alloy, echo, fable, onyx, nova, and shimmer) 
cmd=$(cat <<EOF
curl https://api.openai.com/v1/audio/speech \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{ "voice": "echo", "speed": 1.0, "model": "tts-1", "input": "$@" }' \
    | ffplay -i pipe:0 -nodisp -autoexit -hide_banner -stats -loglevel error
EOF
)
    # -d '$params' \
    # --output output.mp3

echo $cmd
eval $cmd
#afplay output.mp3

