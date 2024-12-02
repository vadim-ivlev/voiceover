package translator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

/*

curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
  "model": "gpt-4o",
  "messages": [
    {
      "role": "system",
      "content": [
        {
          "type": "text",
          "text": "1. Translate the text provided by the user into Russian.\nMake the text easier to read.\n\n2. Transform every sentence according to the following  key principles:\n\nClarity and Simplicity\n\n\t•\tUse clear and concise language.\n\t•\tAvoid jargon, overly complex sentences, and unnecessary words.\n\t•\tPrefer plain language over complicated or obscure vocabulary.\n\t•\tUse active voice whenever possible,.\n\t•\tAvoid vague terms and redundant phrases.\n\n\n3. Eliminate unnecessary adverbs and rearrange adjectives in accord with Russian norms.\n"
        }
      ]
    },
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "Hello World"
        }
      ]
    }
  ],
  "response_format": {
    "type": "text"
  },
  "temperature": 1,
  "max_tokens": 2048,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0
}'

*/

var TextTranslInstructions = `
1. Translate the text provided by the user into %s.
Make the text easier to read.

2. Transform every sentence according to the following  key principles:

Clarity and Simplicity

	•	Use clear and concise language.
	•	Avoid jargon, overly complex sentences, and unnecessary words.
	•	Prefer plain language over complicated or obscure vocabulary.
	•	Use active voice whenever possible,.
	•	Avoid vague terms and redundant phrases.


3. Eliminate unnecessary adverbs and rearrange adjectives in accord with %s norms.
`

var HtmlTranslInstructions = `
1. Translate the text provided by the user into %s.
Make the text easier to read.

2. Transform every sentence according to the following  key principles:

Clarity and Simplicity

	•	Use clear and concise language.
	•	Avoid jargon, overly complex sentences, and unnecessary words.
	•	Prefer plain language over complicated or obscure vocabulary.
	•	Use active voice whenever possible,.
	•	Avoid vague terms and redundant phrases.


3. Eliminate unnecessary adverbs and rearrange adjectives in accord with %s norms.

4. If user's text contains proper names only, just translate them. 

5. If there are any inner tags such as <span>, <a> or <small> , try to keep tags in place. Preserve both the formatting and the semantic flow. 

`

// TranslateText translates the text using OpenAI API
// and returns the translated text.
// Parameters:
// - apiURL: the URL of the OpenAI API
// - apiKey: the API key for the OpenAI API
// - language: the language to translate the text into
// - instructions: the instructions for the translation
// - text: the text to translate
// Returns:
// - translation: the translated text
func TranslateText(apiURL, apiKey, language, instructions, text string) (translation string, err error) {
	systemMessage := fmt.Sprintf(instructions, language, language)
	requestBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemMessage,
			},
			{
				"role":    "user",
				"content": text,
			},
		},
		"temperature": 1,
	}

	jsonData, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return "", err
	}

	// fmt.Println(string(jsonData))

	url := fmt.Sprintf("%s/v1/chat/completions", apiURL)
	// fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		fmt.Println(string(body))

		return "", errors.New("failed to get a valid response from OpenAI API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no translation found in the response")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid response format")
	}

	translatedText, ok := choice["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", errors.New("invalid response format")
	}

	return translatedText, nil
}
