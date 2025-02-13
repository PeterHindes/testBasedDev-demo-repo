package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
)

func runTests(dir string) (string, error) {
	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Change to target directory
	if err := os.Chdir(dir); err != nil {
		return "", fmt.Errorf("failed to change directory: %v", err)
	}

	// Run go test and capture output
	cmd := exec.Command("go", "test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Don't return error as we want to send test failures to AI
		return string(output), nil
	}

	// Ensure we change back to original directory
	defer os.Chdir(currentDir)

	return string(output), nil
}

func runProgram(dir string) (string, error) {
	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Change to target directory
	if err := os.Chdir(dir); err != nil {
		return "", fmt.Errorf("failed to change directory: %v", err)
	}

	// Run the program and capture output
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run program: %v", err)
	}

	// Ensure we change back to original directory
	defer os.Chdir(currentDir)

	return string(output), nil
}

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("api key environment variable is not set")
		return
	}

	// Get target directory from args or use current directory
	targetDir := os.Getenv("HOLDING_PATH")

	// Start the program
	fmt.Println("Starting program...")
	// Create a channel to signal goroutine completion
	done := make(chan bool)
	go func() {
		defer close(done)
		output, err := runProgram(targetDir+"/verified")
		if err != nil {
			fmt.Printf("Error running program: %v\n", err)
			return
		}
		// return output
		var _ = output
	}()


	// Run tests and get output
	testOutput, err := runTests(targetDir+"/tests")
	if err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		return
	}

	// Kill the program after tests are done
	<-done

	// Create configuration with custom base URL
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"

	// Create a new client with config
	client := openai.NewClientWithConfig(config)

	// Create message content with test output
	messageContent := fmt.Sprintf("Here are the test results from my Go project:\n\n%s\n\nPlease analyze these test results and provide insights.", testOutput)

	fmt.Println("Sending message to AI...")
	// then print it in teal
	fmt.Println("\033[36m", messageContent, "\033[0m")

	// Create a completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gemini-2.0-flash",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that analyzes test results and provides insights.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: messageContent,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Error creating chat completion: %v\n", err)
		return
	}

	// Print the response
	fmt.Println(resp.Choices[0].Message.Content)

	// Save the response to a file
	_, err = os.Create("response.txt")
	if err != nil {
		fmt.Printf("Error creating response file: %v\n", err)
		return
	}
}
