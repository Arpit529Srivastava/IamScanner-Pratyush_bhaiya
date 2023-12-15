package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type awsValidator struct{}


// this is specific to aws
func (a awsValidator) FindCredentials(content string, pattern string) ([]string, error) {
	regex := regexp.MustCompile(pattern)
	var matchArr []string
	matches := regex.FindAllString(string(content), -1)
	for _, match := range matches {
		matchArr = regexp.MustCompile(`[^\S]+`).Split(match, 2)
	}
	return matchArr, nil
}


// this function is specific to aws
func (a awsValidator) ValidateCredentials(key string, token string) (error, bool) {

	if ValidateAWSCredentials(key, token) != nil {
		return errors.New("invalid credentials"), false
	}

	return nil, true
}


// this function can be used to validate aws credentials from a pool
// of valid as well as invalid credentials
// this is specific to aws
func ValidateAWSCredentials(accessKey, secretToken string) error {
	// Create a new AWS session
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretToken, ""),
	})

	// Create a new iam service client using the session
	svc := iam.New(sess)

	// Basic API call to check the IAM keys' validity
	d, err := svc.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		// InvalidClientTokenId error occurs for invalid keys.
		// If keys are valid, if the role doesn't have permission
		// to list groups, it returns an AccessDenied error
		if strings.Contains(err.Error(), "InvalidClientTokenId") {
			return nil
		}
		return err
	}

	fmt.Print(d)

	// IAM keys are valid and the role has permission to list groups
	return nil
}
