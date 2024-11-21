#!/bin/bash

# This script will use Google's text-to-speech API to create an MP3 file
# https://console.cloud.google.com/vertex-ai/studio/speech/text-to-speech?project=auth-proxy-1572423489738

# You should receive a JSON response similar to the following
# {
#   audioContent: "uqj;lkfadvnoiqjlkdasjfladklfjkl2j3jq3n4kl...VVVagewred",
# }
# Copy the value of the audioContent field into a new file named synthesize-output-base64.txt using a command like
# $
# echo "<audio_content_value>" > synthesize-output-base64.txt
# Decode the contents of the synthesize-output-base64.txt file into a new file named synthesize-audio.mp3
# $
# base64 synthesize-output-base64.txt -d > synthesize-audio.mp3
# Open the created file on an audio device or in Chrome to listen
# $https://console.cloud.google.com/vertex-ai/studio/speech/text-to-speech?project=auth-proxy-1572423489738
# google-chrome ./my/path/to/synthesize-audio.mp3



# Create an MP3 file from the text by sending a POST request to the Google Text-to-Speech API
curl -X POST -H "Content-Type: application/json" \
-H "X-Goog-User-Project: $(gcloud config list --format='value(core.project)')" \
-H "Authorization: Bearer $(gcloud auth print-access-token)" \
--data '{
"input": {
  "text": "‘What sort of things?’ I asked him."
},
"voice": {
  "languageCode": "en-GB",
  "name": "en-GB-Journey-D"
},
"audioConfig": {
  "audioEncoding": "LINEAR16"
}
}' "https://texttospeech.googleapis.com/v1/text:synthesize" \
| jq -r .audioContent | base64 -d > out.mp3


# Play the audio file
ffplay out.mp3 -autoexit -nodisp -loglevel info


