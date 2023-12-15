package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

// Clones the repository into the given dir, just as a normal git clone does
func cloneRepository(repoUrl string, dir string) error {

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func getAllBranches(dirName string) ([]string, error) {
	cmd := exec.Command("git", "-C", dirName, "branch", "-r", "--format", "%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var newBranch []string
	for _, val := range branches {
		parts := strings.Split(val, "/")
		if len(parts) > 0 {
			newBranch = append(newBranch, parts[len(parts)-1])
		}
	}
	return newBranch, nil
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [repository path]")
	}

	repoUrl := os.Args[1]

	dir, e := os.MkdirTemp("", "example")
	if e != nil {
		log.Fatal(e)
	}

	logsFolder := "logs"
	logsPath := fmt.Sprintf("%s/output.txt", logsFolder) // Path to the output file

	if _, e := os.Stat(logsFolder); os.IsNotExist(e) {
		os.Mkdir(logsFolder, 0755)
	}

	outputFile, outputErr := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if outputErr != nil {
		log.Fatal(outputErr)
	}

	defer outputFile.Close()

	customLog := log.New(&CustomLogger{Output: outputFile}, "", 0)

	err := cloneRepository(repoUrl, dir)
	if err != nil {
		fmt.Println("Error cloning the repository")
	}

	branches, err := getAllBranches(dir)

	rs := NewRepoScanner(dir, customLog)

	start := time.Now()
	rs.ScanBranches(branches)
	fmt.Println("TOOK ",time.Since(start))

	if err != nil {
		fmt.Println("Error getting all the branches")
	}

	// clean up temp dir
	defer os.RemoveAll(dir)

}
