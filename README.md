# GPTCMD

A chatgpt cli powered by cobra, completely ripped of from https://kadekillary.work/posts/1000x-eng/.

```
Basic prompt: gptcmd prompt "What is the meaning of life?"
	You can also pass a file as input with the --file flag: gptcmd prompt "What is the meaning of life?" --file input.txt
	You can tweak the model and temperature as well with the --model and --temperature flags: gptcmd prompt "What is the meaning of life?" --model davinci --temperature 0.5

Usage:
  gptcmd prompt [flags]

Flags:
  -d, --data string         Optional data string to pass for a prompt
  -f, --file string         Optional data file to pass for a prompt
  -h, --help                help for prompt
  -m, --model string        The model to use for the chatbot (default "gpt-3.5-turbo")
  -t, --temperature float   The temperature to use for the chatbot (default 0.7)
```