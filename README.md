# voiceover


The application will:
1. Read the input text file
2. Convert each line to speech using TTS API
3. Combine all audio files into a single MP3
4. Generate output files (.mp3, .txt, .log.json)


Here are the key capabilities:


1. Convert text to speech:
- Uses OpenAI's TTS API to generate MP3 files from text
- Supports multiple voices (alloy, echo, fable, onyx, nova, shimmer)
- Can control speech speed
- Main functionality in generate.go

2. Text processing:
- Can split text files into lines
- Handles text file reading/writing
- Text processing code in text.go


3. Pipeline operations:
- Processes text in parallel using worker pipelines
- Combines multiple audio files into single output
- Pipeline code in operations.go


## Configuration

The configuration can be customized through environment variables and command line flags defined in config.go

Environment variables can be defined in OS, voiceover.env, .env files listed in order of precedence.

```bash
  # OpenAI or FreeAI API
  OPENAI_API_URL="url for OpenAI API or FreeAI"
  API_KEY="API key for OpenAI or FreeAI"

  #
  GCLOUD_PROJECT="Google Cloud project ID"
  GCLOUD_ACCESS_TOKEN="Google Cloud access token"
```


To use the application:

```sh
# Set your OpenAI API key
export OPENAI_API_KEY=<your-key>

# Run the application
go run main.go -i input.txt -output output.mp3
```


The configuration can be customized through environment variables or command line flags defined in config.go

```json
{
    "input":
        {
            "text":"The Boy Who Talked with Animals."
        },
        "voice":
        {
            "languageCode":"en-GB",
            "name":"en-GB-Journey-D"
        },
        "audioConfig":
        {
            "audioEncoding":"MP3"
        }
    }

{
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
}    
```
