/*
Copyright © 2023 mactynow charles@mactynow.ovh
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type ChatRequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
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
		messages := []Message{}
		prompt := args[0]
		model, _ := cmd.Flags().GetString("model")
		temperature, _ := cmd.Flags().GetFloat64("temperature")
		data, _ := cmd.Flags().GetString("data")
		if data != "" {
			prompt = fmt.Sprintf("%s: %s", prompt, data)
		}

		fileName, _ := cmd.Flags().GetString("file")
		if fileName != "" {
			file, err := os.ReadFile(fileName)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}
			content := string(file)
			prompt = fmt.Sprintf("%s: %s", prompt, content)
		}

		promptMessage := Message{Role: "user", Content: prompt}

		newChat, _ := cmd.Flags().GetBool("new-chat")
		if newChat {
			if _, err := os.Stat(historyFilePath); err == nil {
				os.Remove(historyFilePath)
			}
		} else {
			messages = readHistory()
		}

		storeHistory(promptMessage)
		messages = append(messages, promptMessage)
		requestBody := ChatRequestBody{
			Model:       model,
			Messages:    messages,
			Temperature: temperature,
		}

		resp, err := sendRequesttoOpenAI(requestURL, requestBody)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(resp)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringP("data", "d", "", "Optional data string to pass for a prompt")
	promptCmd.Flags().StringP("file", "f", "", "Optional data file to pass for a prompt")
	promptCmd.Flags().Float64P("temperature", "t", 0.7, "The temperature to use for the chatbot")
	promptCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "The model to use for the chatbot")
	promptCmd.Flags().BoolP("new-chat", "n", false, "Start a new chat conversation")
}
