package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ShellToUse = "bash"
const AgentDir = "/tmp/otlp-agent"

var PostConfigSetupCommands = [...]string{
	fmt.Sprintf("%s/ocb --config=%s/builder-config.yaml", AgentDir, AgentDir),
	fmt.Sprintf("%s/collector/otelcorecol --config=%s/otelcol.yaml", AgentDir, AgentDir),
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) error {
	fileContent, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("Failed to read file content")
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return errors.New("Missing Content-Type header")
	}

	parts := strings.Split(contentType, "/")
	if len(parts) != 2 {
		return errors.New("Invalid Content-Type")
	}

	extension := parts[1]

	fileName := fmt.Sprintf("otelcol.%s", extension)
	filePath := filepath.Join(fmt.Sprintf("%s/", AgentDir), fileName)

	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to write file to disk")
	}

	fmt.Println("finished saving file")
	return nil
}

func execShellCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe: %v", err)
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating stderr pipe: %v", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting command: %v", err)
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Printf("[stdout]: %s", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stdout: %v", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("[stderr]: %s", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stderr: %v", err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution failed: %v", err)
		return err
	}

	log.Println("Command finished successfully.")
	return nil
}

func setUpConfigAndStartCollector(w http.ResponseWriter, r *http.Request) {
	if err := fileUploadHandler(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error: %v", err)
		return
	}

	go func() {
		for _, cmd := range PostConfigSetupCommands {
			fmt.Printf("Running command %s in shell\n", cmd)
			err := execShellCommand(cmd)
			if err != nil {
				fmt.Printf("command %s failed with error :: %s, %v\n", cmd, err.Error(), err.Error())
			}
		}
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Agent successfully configured. Please wait for agent boot up")
}

func main() {
	http.HandleFunc("/collector/config", setUpConfigAndStartCollector)

	fmt.Println("Server starting on port 4343...")
	err := http.ListenAndServe(":4343", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
