package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	Args:  cobra.ExactArgs(1),
	Long: `Basic prompt: gptcmd prompt "What is the meaning of life?"
	You can also pass a file as input with the --file flag: gptcmd prompt "What is the meaning of life?" --file input.txt
	You can tweak the model and temperature as well with the --model and --temperature flags: gptcmd prompt "What is the meaning of life?" --model davinci --temperature 0.5`,
	Run: func(cmd *cobra.Command, args []string) {
		requestURL := "https://api.openai.com/v1/chat/completions"

		prompt := args[0]
		data, _ := cmd.Flags().GetString("data")
		model, _ := cmd.Flags().GetString("model")
		temperature, _ := cmd.Flags().GetFloat64("temperature")
		fileName, _ := cmd.Flags().GetString("file")
		openAIApiKey := os.Getenv("OPENAI_API_KEY")

		if openAIApiKey == "" {
			log.Println("Please set the OPENAI_API_KEY environment variable")
			return
		}

		if data != "" {
			prompt = fmt.Sprintf("%s: %s", prompt, data)
		}

		if fileName != "" {
			file, err := os.ReadFile(fileName)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}
			content := string(file)
			prompt = fmt.Sprintf("%s: %s", prompt, content)
		}

		requestBody := RequestBody{
			Model:       model,
			Messages:    []Message{{Role: "user", Content: prompt}},
			Temperature: temperature,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
		}

		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(jsonBody))
		if err != nil {
			log.Println(err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+openAIApiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		var responseBody ResponseBody
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			log.Println(err)
		}

		if len(responseBody.Choices) == 0 {
			log.Println("No response:", responseBody.Usage)
			return
		}

		fmt.Println(responseBody.Choices[0].Message.Content)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringP("data", "d", "", "Optional data string to pass for a prompt")
	promptCmd.Flags().StringP("file", "f", "", "Optional data file to pass for a prompt")
	promptCmd.Flags().Float64P("temperature", "t", 0.7, "The temperature to use for the chatbot")
	promptCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "The model to use for the chatbot")
}
