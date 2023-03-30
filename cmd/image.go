/*
Copyright © 2023 mactynow charles@mactynow.ovh
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type ImageRequestBody struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type Data struct {
	Url string `json:"url"`
}

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image",
	Args:  cobra.ExactArgs(1),
	Long:  `Generate an image: gptcmd image "A picture of a cat"`,
	Run: func(cmd *cobra.Command, args []string) {
		requestURL := "https://api.openai.com/v1/images/generations"
		prompt := args[0]
		n, _ := cmd.Flags().GetInt("n")
		size, _ := cmd.Flags().GetString("size")

		requestBody := ImageRequestBody{
			Prompt: prompt,
			N:      n,
			Size:   size,
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
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().IntP("n", "n", 1, "Number of images to generate")
	imageCmd.Flags().StringP("size", "s", "256x256", "Size of the image")
}
