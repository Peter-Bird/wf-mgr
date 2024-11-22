package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	mdl "github.com/Peter-Bird/models"
)

func ExecuteWorkflow(w http.ResponseWriter, workflow mdl.Workflow) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Flushing not supported by ResponseWriter\n")
		return
	}

	for i, step := range workflow.Steps {
		if err := executeStep(step); err != nil {
			log.Printf("Step failed: %v", err)
			return
		}

		fmt.Fprintf(w, "Step %d Executed\n", i)
		flusher.Flush()

	}
}

// executeStep executes an individual step of a workflow
func executeStep(step mdl.Step) error {

	red := "\033[31m"
	reset := "\033[0m"

	log.Printf("Executing step:\n%s%s %s%s\n\n", red, step.Method, step.Endpoint, reset)

	client := &http.Client{}
	var req *http.Request
	var err error

	if step.Method == http.MethodPost || step.Method == http.MethodPut {
		body, _ := json.Marshal(step.Parameters)
		req, err = http.NewRequest(step.Method, step.Endpoint, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(step.Method, step.Endpoint, nil)
	}

	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		log.Printf("Request failed: %v", err)
		return err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the HTML content
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return err
	}

	logHTMLContent(htmlContent)

	return nil
}

func logHTMLContent(htmlContent []byte) {
	// Check if the content is JSON
	var jsonContent interface{}
	if err := json.Unmarshal(htmlContent, &jsonContent); err == nil {
		// Pretty print the JSON content
		prettyJSON, err := json.MarshalIndent(jsonContent, "", "    ")
		if err != nil {
			log.Printf("Failed to format JSON: %v", err)
		} else {
			log.Printf("Response JSON:\n%s\n\n", string(prettyJSON))
		}
	} else {
		// Fallback: Print as a regular string
		log.Printf("Response HTML:\n%s\n\n", string(htmlContent))
	}
}
