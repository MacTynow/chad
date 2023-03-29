/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens      int `json:"prompt_tokens"`
	CompletionsTokens int `json:"completions_tokens"`
	TotalTokens       int `json:"total_tokens"`
}

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Ask a question and get a response",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		requestURL := "https://api.openai.com/v1/chat/completions"

		prompt, _ := cmd.Flags().GetString("prompt")
		data, _ := cmd.Flags().GetString("data")
		model, _ := cmd.Flags().GetString("model")
		temperature, _ := cmd.Flags().GetFloat64("temperature")

		if data != "" {
			prompt = fmt.Sprintf("%s: %s", prompt, data)
		}

		requestBody := RequestBody{
			Model:       model,
			Messages:    []Message{{Role: "user", Content: prompt}},
			Temperature: temperature,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println(err)
		}

		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(jsonBody))
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()

		var responseBody ResponseBody
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(responseBody.Choices[0].Message.Content)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringP("prompt", "p", "Say hello", "The prompt to use for the chatbot")
	promptCmd.Flags().StringP("data", "d", "", "Optional data to pass for a prompt")
	promptCmd.Flags().Float64P("temperature", "t", 0.7, "The temperature to use for the chatbot")
	promptCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "The model to use for the chatbot")
}
