/*
Copyright © 2023 mactynow charles@mactynow.ovh
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type RequestBody interface{}

type ResponseBody interface{}

type ImageResponseBody struct {
	Created int64 `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type ChatResponseBody struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type EditResponseBody struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

func sendRequesttoOpenAI(requestURL string, requestBody RequestBody) (string, error) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	if openAIApiKey == "" {
		log.Println("Please set the OPENAI_API_KEY environment variable")
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBody := handleResponseBody(requestURL)

	err = json.NewDecoder(resp.Body).Decode(responseBody)
	if err != nil {
		return "", err
	}

	return returnResponseString(responseBody), nil
}

func handleResponseBody(requestURL string) ResponseBody {
	switch {
	case strings.Contains(requestURL, "images"):
		return &ImageResponseBody{}
	case strings.Contains(requestURL, "chat"):
		return &ChatResponseBody{}
	case strings.Contains(requestURL, "edit"):
		return &EditResponseBody{}
	default:
		return ""
	}
}

func returnResponseString(responseBody ResponseBody) string {
	switch responseBody := responseBody.(type) {
	case *ImageResponseBody:
		return responseBody.Data[0].Url
	case *ChatResponseBody:
		choices := responseBody.Choices
		if len(choices) == 0 {
			return "no response"
		}
		return choices[0].Message.Content
	case *EditResponseBody:
		choices := responseBody.Choices
		if len(choices) == 0 {
			return "no response"
		}
		return choices[0].Text
	default:
		return "invalid response body"
	}
}
