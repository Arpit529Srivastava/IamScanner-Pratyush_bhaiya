package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type RepoScanner struct {
	dir       string
	customLog *log.Logger
}

// initializes a new RepoScanner instance.
func NewRepoScanner(dir string, customLog *log.Logger) *RepoScanner {
	return &RepoScanner{dir: dir, customLog: customLog}
}

// scans through a list of branches in the repository
func (rs *RepoScanner) ScanBranches(branches []string) {
	for _, val := range branches {
		rs.switchAndScan(val)
	}
}

// switches to a branch and scans its commit history and files
func (rs *RepoScanner) switchAndScan(val string) {

	wg := sync.WaitGroup{}

	rs.switchBranch(val)
	commits, _ := rs.getCommitHistory()
	for _, commit := range commits {
		files, _ := rs.listFilesInCommit(commit)
		for _, file := range files {
			wg.Add(1)
			go rs.scanFileContent(val, commit, file, &wg)
		}
		wg.Wait()
	}
}

// retrieves the commit history of the current branch
func (rs *RepoScanner) getCommitHistory() ([]string, error) {
	cmd := exec.Command("git", "-C", rs.dir, "rev-list", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	return commits, nil
}

func (rs *RepoScanner) switchBranch(branchName string) error {
	cmd := exec.Command("git", "-C", rs.dir, "checkout", branchName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error switching branch: %w", err)
	}
	return nil
}

// fetches the content of a file in a specific commit
func (rs *RepoScanner) getFileContentFromCommit(commitHash, filePath string) (string, error) {
	cmd := exec.Command("git", "-C", rs.dir, "show", fmt.Sprintf("%s:%s", commitHash, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// lists all files in a specific commit
func (rs *RepoScanner) listFilesInCommit(commitHash string) ([]string, error) {
	cmd := exec.Command("git", "-C", rs.dir, "ls-tree", "-r", commitHash)
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

// scans the content of a file in a commit for AWS credentials and logs any matches
func (rs *RepoScanner) scanFileContent(branch, commit, file string, wg *sync.WaitGroup) {

	defer wg.Done()

	fileContent, err := rs.getFileContentFromCommit(commit, file)
	if err != nil {
		fmt.Println("Error getting file content for", file, "in commit", commit)
		return
	}
	validator := awsValidator{}
	matches, err := validator.FindCredentials(fileContent)
	if err != nil {
		fmt.Println("Error scanning file", file, "in commit", commit)
		return
	}

	if len(matches) > 0 {
		rs.customLog.Println("Branch: ", branch)
		rs.customLog.Println("\t File: ", file, "Commit: ", commit)

		for _, val := range matches {
			rs.customLog.Println("\t Access Key: ", val.Id)
			rs.customLog.Println("\t Secret Token: ", val.Token)

		}
	}
}
