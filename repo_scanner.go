package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type RepoScanner struct {
	dir       string
	customLog *log.Logger
}

func NewRepoScanner(dir string, customLog *log.Logger) *RepoScanner {
	return &RepoScanner{dir: dir, customLog: customLog}
}

func (rs *RepoScanner) ScanBranches(branches []string) {
	// var wg sync.WaitGroup

	for _, val := range branches {
		rs.switchAndScan(val)
	}

}

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

func (rs *RepoScanner) getFileContentFromCommit(commitHash, filePath string) (string, error) {
	cmd := exec.Command("git", "-C", rs.dir, "show", fmt.Sprintf("%s:%s", commitHash, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (rs *RepoScanner) scanFileWithRegex(fileContent string) ([]string, error) {
	r := regexp.MustCompile(`(?m)(?i)AKIA[0-9A-Z]{16}\s+\S{40}|AWS[0-9A-Z]{38}\s+?\S{40}`)

	matches := r.FindAllString(fileContent, -1)
	var matchArr []string
	for _, match := range matches {
		matchArr = regexp.MustCompile(`[^\S]+`).Split(match, 2)
	}
	return matchArr, nil
}

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

func (rs *RepoScanner) scanFileContent(branch, commit, file string, wg *sync.WaitGroup) {

	defer wg.Done()

	fileContent, err := rs.getFileContentFromCommit(commit, file)
	if err != nil {
		fmt.Println("Error getting file content for", file, "in commit", commit)
		return
	}

	matches, err := rs.scanFileWithRegex(fileContent)
	if err != nil {
		fmt.Println("Error scanning file", file, "in commit", commit)
		return
	}

	if len(matches) > 0 {

		log.Println("LOGGED ",branch)

		rs.customLog.Println("Branch: ", branch)
		rs.customLog.Println("\t File: ", file, "Commit: ", commit)
		rs.customLog.Println("\t Access Key: ", matches[0])
		rs.customLog.Println("\t Secret Token: ", matches[1])
	}
}
