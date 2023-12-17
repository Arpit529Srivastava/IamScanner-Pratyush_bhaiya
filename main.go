package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Clones the repository in a given dir, just as a normal git clone does
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

// get all the branches of the specified repository url
func getAllBranches(dirName string) ([]string, error) {

	// command to get the remote branches
	cmd := exec.Command("git", "-C", dirName, "branch", "-r", "--format", "%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// returns []string containing branch names
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var newBranch []string

	// branches were in the form remote/origin/leak
	// hence extracted just the last part in a new array
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
		log.Fatal("Usage: go run . [repository path]")
	}

	repoUrl := os.Args[1]

	// creating a temp directory
	dir, e := os.MkdirTemp("", "example")
	if e != nil {
		log.Fatal(e)
	}

	// creating a folder by the name of logs to store the output
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

	// using a custom logger to log the output
	customLog := log.New(&CustomLogger{Output: outputFile}, "", 0)

	// cloning the repository, this is essentially the first step
	// of this entire project
	err := cloneRepository(repoUrl, dir)
	if err != nil {
		fmt.Println("Error cloning the repository")
	}

	// get all the branches and pass them to ScanBranches() function
	branches, err := getAllBranches(dir)
	if err != nil {
		fmt.Println("Error getting all the branches")
	}

	rs := NewRepoScanner(dir, customLog)
	rs.ScanBranches(branches)


	// clean up temp dir
	defer os.RemoveAll(dir)

}
