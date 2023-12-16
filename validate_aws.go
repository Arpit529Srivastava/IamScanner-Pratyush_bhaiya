package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const awsKeyPattern = `(?m)(?i)AKIA[0-9A-Z]{16}\s+\S{40}|AWS[0-9A-Z]{38}\s+?\S{40}`

type awsValidator struct{}

// this is specific to aws
func (a awsValidator) FindCredentials(content string) ([]Credentials, error) {
	res := []Credentials{}
	regex := regexp.MustCompile(awsKeyPattern)

	matches := regex.FindAllString(string(content), -1)
	for _, match := range matches {
		matchArr := regexp.MustCompile(`[^\S]+`).Split(match, 2)
		res = append(res, Credentials{
			Id:    matchArr[0],
			Token: matchArr[1],
		})
	}
	return res, nil
}

// this function is specific to aws
func (a awsValidator) ValidateCredentials(c Credentials) bool {
	return ValidateAwsCredentials(c.Id, c.Token)
}

// this function can be used to validate aws credentials from a pool
// of valid as well as invalid credentials
// this is specific to aws
// func ValidateAwsCredentials(accessKeyID, secretAccessKey string) bool {
// 	// Create a new AWS session with the IAM keys
// 	sess, _ := session.NewSession(&aws.Config{
// 		Region:      aws.String("us-west-2"),
// 		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
// 	})
// 	// Unnecessary to handle error since session created is static
// 	// it doesn't send any request

// 	// Create a new iam service client using the session
// 	svc := iam.New(sess)

// 	// Basic API call to check the IAM keys' validity
// 	d, err := svc.ListGroups(&iam.ListGroupsInput{})
// 	if err != nil {
// 		// InvalidClientTokenId error occurs for invalid keys.
// 		// If keys are valid, if the role doesn't have permission
// 		// to list groups, it returns an AccessDenied error
// 		if strings.Contains(err.Error(), "InvalidClientTokenId") {
// 			return false
// 		}
// 		return true
// 	}

// 	fmt.Print(d)

// 	// IAM keys are valid and the role has permission to list groups
// 	return true
// }

func ValidateAwsCredentials(accessKey, secretKey string) bool {

	log.Println("HERE: ", accessKey)
	log.Println("HERE: ", secretKey)

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"), // Change the region as per your requirements
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	// Create an IAM service client
	svc := iam.New(sess)

	_, err := svc.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		// InvalidClientTokenId error occurs for invalid keys.
		// If keys are valid, if the role doesn't have permission
		// to list groups, it returns an AccessDenied error
		if strings.Contains(err.Error(), "InvalidClientTokenId") {
			return false
		}
		return true
	}

	// IAM keys are valid and the role has permission to list groups
	return true
}
