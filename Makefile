# Set the name of the binary to be built
BINARY_NAME=gptcmd

# Set the installation directory
INSTALL_DIR=~/go/bin

# Set the Go compiler and flags
GO=go
GOFLAGS=-ldflags="-s -w"

# Define the build target
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME)

# Define the install target
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)

# Define the clean target
clean:
	rm -f $(BINARY_NAME)