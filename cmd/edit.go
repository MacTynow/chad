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

type EditRequestBody struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Instruction string  `json:"instruction"`
	Temperature float64 `json:"temperature"`
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Uses the edit api to improve an input",
	Args:  cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string

		requestURL := "https://api.openai.com/v1/edits"
		prompt := args[0]
		model, _ := cmd.Flags().GetString("model")
		temperature, _ := cmd.Flags().GetFloat64("temperature")
		data, _ := cmd.Flags().GetString("data")
		if data != "" {
			input = data
		}

		fileName, _ := cmd.Flags().GetString("file")
		if fileName != "" {
			file, err := os.ReadFile(fileName)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}
			input = string(file)
		}

		if input == "" {
			log.Println("No input provided")
			return
		}

		requestBody := EditRequestBody{
			Model:       model,
			Input:       input,
			Instruction: prompt,
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
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("data", "d", "", "Optional data string to pass for a prompt")
	editCmd.Flags().StringP("file", "f", "", "Optional data file to pass for a prompt")
	editCmd.Flags().Float64P("temperature", "t", 0.7, "The temperature to use for the chatbot")
	editCmd.Flags().StringP("model", "m", "code-davinci-edit-001", "The model to use for the chatbot")
}
