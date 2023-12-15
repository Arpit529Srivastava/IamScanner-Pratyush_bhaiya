package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

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

func getCommitHistory(dirName string) ([]string, error) {

	// creates a command to execute the Git command-line tool
	cmd := exec.Command("git", "-C", dirName, "rev-list", "HEAD")

	output, err := cmd.Output()
	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	return commits, err
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

func switchBranch(dirName, branchName string) error {
	cmd := exec.Command("git", "-C", dirName, "checkout", branchName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error switching branch: %w", err)
	}
	return nil
}

func getFileContentFromCommit(dirName, commitHash, filePath string) (string, error) {
	cmd := exec.Command("git", "-C", dirName, "show", fmt.Sprintf("%s:%s", commitHash, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func scanFileWithRegex(fileContent string) ([]string, error) {
	r := regexp.MustCompile(`(?m)(?i)AKIA[0-9A-Z]{16}\s+\S{40}|AWS[0-9A-Z]{38}\s+?\S{40}`)

	matches := r.FindAllString(fileContent, -1)
	var matchArr []string
	for _, match := range matches {
		matchArr = regexp.MustCompile(`[^\S]+`).Split(match, 2)
	}
	return matchArr, nil
}

func listFilesInCommit(dirName string, commitHash string) ([]string, error) {
	cmd := exec.Command("git", "-C", dirName, "ls-tree", "-r", commitHash)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 3 {
			files = append(files, parts[3]) // the file path
		}
	}

	return files, nil
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

	if err != nil {
		fmt.Println("Error getting all the branches")
	}

	for _, val := range branches {
		switchBranch(dir, val)
		commits, _ := getCommitHistory(dir)
		for _, commit := range commits {
			files, _ := listFilesInCommit(dir, commit)
			for _, file := range files {
				fileContent, err := getFileContentFromCommit(dir, commit, file)
				if err != nil {
					fmt.Println("Error getting file content for", file, "in commit", commit)
					continue
				}

				matches, err := scanFileWithRegex(fileContent)
				if err != nil {
					fmt.Println("Error scanning file", file, "in commit", commit)
					continue
				}

				if len(matches) > 0 {
					customLog.Println("Branch: ", val)
					customLog.Println("\t File: ", file, "Commit: ", commit)
					customLog.Println("\t Access Key: ", matches[0])
					customLog.Println("\t Secret Token: ", matches[1])
				}
			}
		}

	}

	// clean up temp dir
	defer os.RemoveAll(dir)

}
