#!/bin/bash

# Check if the OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Error: OPENAI_API_KEY is not set"
    exit 1
fi

# Check if there is an argument
if [ $# -eq 0 ]; then
    echo "Error: No argument provided"
    exit 1
fi

#https://serverfault.com/questions/72476/clean-way-to-write-complex-multi-line-string-to-a-variable



# Experiment with different voices (alloy, echo, fable, onyx, nova, and shimmer) 
cmd=$(cat <<EOF
curl https://api.openai.com/v1/audio/speech \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{ "voice": "echo", "speed": 1.0, "model": "tts-1", "input": "$@" }' \
    | ffplay -i pipe:0 -nodisp -autoexit -hide_banner -stats -loglevel error
EOF
)
    # --output openai-speak.mp3


echo $cmd
eval $cmd
#afplay openai-speak.mp3
# ffplay openai-speak.mp3 -autoexit -nodisp -loglevel info

