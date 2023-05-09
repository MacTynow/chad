/*
Copyright © 2023 mactynow charles@mactynow.ovh
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const historyFilePath = "/tmp/messages.json"

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
			Role    string `json:"role"`
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

		storeHistory(choices[0].Message)

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

func storeHistory(message Message) {
	file, err := os.OpenFile(historyFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := file.Write(jsonBytes); err != nil {
		log.Println(err)
		return
	}

	if _, err := file.WriteString("\n"); err != nil {
		log.Println(err)
		return
	}
}

func readHistory() []Message {
	messages := []Message{}
	file, err := os.OpenFile(historyFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var message Message
		if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
			log.Println("Error parsing JSON:", err)
			continue
		}

		messages = append(messages, message)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
		return nil
	}

	return messages
}
