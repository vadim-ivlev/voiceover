The following bash command convers a text into an mp3 file:

```bash
# Create an MP3 file from the text by sending a POST request to Elevenlabs Text-to-Speech API

curl --request POST \
  --url https://api.elevenlabs.io/v1/text-to-speech/{voice_id} \
  --header 'Content-Type: application/json' \
  --data '{
  "text": "<string>",
  "model_id": "<string>",
  "language_code": "<string>",
  "voice_settings": {
    "stability": 123,
    "similarity_boost": 123,
    "style": 123,
    "use_speaker_boost": true
  },
  "pronunciation_dictionary_locators": [
    {
      "pronunciation_dictionary_id": "<string>",
      "version_id": "<string>"
    }
  ],
  "seed": 123,
  "previous_text": "<string>",
  "next_text": "<string>",
  "previous_request_ids": [
    "<string>"
  ],
  "next_request_ids": [
    "<string>"
  ],
  "use_pvc_as_ivc": true,
  "apply_text_normalization": "auto"
}'

```

Write a go function that takes a string and file name as input and generates an mp3 file with the given text.

func GenerateElevenlabsSpeechMP3(speed float64, voice, text, fileName string) error