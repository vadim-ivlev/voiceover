# voiceover


Based on the workspace code, this appears to be a Go application that converts text files to voiceovers using OpenAI's text-to-speech API. Here are the key capabilities:

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



To use the application:

```sh
# Set your OpenAI API key
export OPENAI_API_KEY=<your-key>

# Run the application
go run main.go -i input.txt -o output.mp3
```

The application will:
1. Read the input text file
2. Convert each line to speech using OpenAI's API
3. Combine all audio files into a single MP3
4. Generate output files (.mp3, .txt, .log.json)

The configuration can be customized through environment variables or command line flags defined in config.go

