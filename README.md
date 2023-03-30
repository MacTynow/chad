# GPTCMD

A chatgpt cli powered by cobra, completely ripped of from https://kadekillary.work/posts/1000x-eng/.

## Usage

```
A CLI for interacting with GPT

Usage:
  gptcmd [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  image       Generate an image
  prompt      Ask a question and get a response

Flags:
  -h, --help     help for gptcmd
  -t, --toggle   Help message for toggle

Use "gptcmd [command] --help" for more information about a command.
```

### Chat 

```
Usage:
  gptcmd prompt [flags]

Flags:
  -d, --data string         Optional data string to pass for a prompt
  -f, --file string         Optional data file to pass for a prompt
  -h, --help                help for prompt
  -m, --model string        The model to use for the chatbot (default "gpt-3.5-turbo")
  -t, --temperature float   The temperature to use for the chatbot (default 0.7)
```

### Image 

```
Usage:
  gptcmd image [flags]

Flags:
  -h, --help          help for image
  -n, --n int         Number of images to generate (default 1)
  -s, --size string   Size of the image (default "256x256")
```