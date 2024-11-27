The following bash command convers a text into an mp3 file:

```bash
# Create an MP3 file from the text by sending a POST request to the Google Text-to-Speech API

curl -X POST -H "Content-Type: application/json" \
-H "X-Goog-User-Project: $GCLOUD_PROJECT" \
-H "Authorization: Bearer $GCLOUD_ACCESS_TOKEN" \
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
| jq -r .audioContent | base64 -d > google-speak.mp3

```

Write a go function that takes a string and file name as input and generates an mp3 file with the given text.

func GenerateGoogleSpeechMP3(speed float64, voice, text, fileName string) error