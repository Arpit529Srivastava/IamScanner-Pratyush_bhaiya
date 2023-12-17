package main

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// this is aws specific
type awsValidator struct{}

// this is specific to aws
// it will also check if the keys are base64 encoded & will thus proceed
// to decode those keys
func (a awsValidator) FindCredentials(content string) ([]Credentials, error) {
	res := []Credentials{}
	regex := regexp.MustCompile(AwsKeyPattern)

	matches := regex.FindAllString(string(content), -1)
	for _, match := range matches {
		if IsBase64Encoded(match) {
			match, _ = DecodeBase64(match) //decode base64
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

// Regex will give us valid as well as invalid credentials so the purpose
// of the below function is to segregate the valid tokens from the invalid ones
// NOTE: IT SEEMS THAT THE CREDENTIALS HAVE EXPIRED AS I GOT INVALID CREDENTIALS FOR ALL OF THE KEYS
func ValidateAwsCredentials(accessKey, secretKey string) bool {

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"), // Change the region as per your requirements
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	// Create STS service client
	svc := sts.New(sess)

	// Call the GetCallerIdentity API to validate the credentials
	_, e := svc.GetCallerIdentity(nil)

	if e != nil {
		fmt.Println("Error verifying credentials:", e)
		return false
	}

	// If no error occurred, credentials are valid
	fmt.Println("Credentials are valid")

	return true
}
