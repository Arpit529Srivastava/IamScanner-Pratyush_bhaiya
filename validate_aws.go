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

type awsValidator struct{}

// this is specific to aws
func (a awsValidator) FindCredentials(content string) ([]Credentials, error) {
	res := []Credentials{}
	regex := regexp.MustCompile(AwsKeyPattern)

	matches := regex.FindAllString(string(content), -1)
	for _, match := range matches {
		if IsBase64Encoded(match) {
			match,_ = DecodeBase64(match)
			matchArr := regexp.MustCompile(`[^\S]+`).Split(match, 2)
			res = append(res, Credentials{
				Id:    matchArr[0],
				Token: matchArr[1],
			})
		} else {
			matchArr := regexp.MustCompile(`[^\S]+`).Split(match, 2)
			res = append(res, Credentials{
				Id:    matchArr[0],
				Token: matchArr[1],
			})
		}
	}
	return res, nil
}

// this function is specific to aws
func (a awsValidator) ValidateCredentials(c Credentials) bool {
	return ValidateAwsCredentials(c.Id, c.Token)
}

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
		return !strings.Contains(err.Error(), "InvalidClientTokenId")
	}

	// IAM keys are valid and the role has permission to list groups
	return true
}
