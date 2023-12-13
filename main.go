package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Clones the repository into the given dir, just as a normal git clone does
func cloneRepository(repoUrl string) error {

	// Cloned the repository onto my desktop under Devops-Node folder
	_, err := git.PlainClone("C:/Users/aries/Desktop/Devops-Node", false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}

// converts byte slice to string slice
func byteSliceToStringSlice(byteSlice []byte) []string {
	var stringSlice []string
	for _, b := range byteSlice {
		stringSlice = append(stringSlice, string(b))
	}
	return stringSlice
}

// function to fetch the commit history of a repository
func getCommitHistory(dirName string) ([]string, error) {
	
	// creates a command to execute the Git command-line tool
	cmd := exec.Command("git", "-C", dirName,"rev-list", "HEAD")

	output, err := cmd.Output()
	commits := strings.Split(strings.Join(byteSliceToStringSlice(output), ""), "\n")
	return commits, err
}

func main() {

	// Invoke method for cloning the repository locally
	err := cloneRepository("https://github.com/abhishek-pingsafe/Devops-Node")
	if err != nil {
		fmt.Println("Error cloning the repository")
	}

	// gets the commits history for the specified directory
	// same dir where I had cloned the repository
	commits, err := getCommitHistory("C:/Users/aries/Desktop/Devops-Node")
	if err != nil {
		fmt.Println("Error getting commit history",err)
	}

	// Print the commit history
	for _, val := range commits {
		fmt.Println(val)
	}
}
