# IAMScanner
The goal of this assignment is to create a GoLang script that scans a Git repository
(GitHub/GitLab/Bitbucket) for any embedded, valid AWS IAM keys. This includes both keys
in the latest code and those in the repository's historical commits in all branches

##  Approach
To conquer this task , I had divided the project into the following steps:

-   Step 1: Clone the repository locally
-   Step 2: Traverse through all branches and, for each branch, iterate over its respective commits and their associated files
-   Step 3: Once you have started scanning the files, start looking for contents that will match the regex
-   Step 4: If there is a match output it to a folder. In my case it was under logs

## Usage
- In the terminal run the command given below:
        
        go run . <path_to_repository>

    In this case:

        go run . https://github.com/abhishek-pingsafe/Devops-Node

- You would see the progress of the repository being cloned in your terminal and right afterwards you would see the desired output in your `logs` folder under the name of `output.txt`
- Currently the outputs of evey repository will be logged here. Further enhancements may include creating a seperate directory for each repository

## Result 
<img width="1600" alt="image" src="https://github.com/PratyushSingh07/IssueHive/assets/90026952/d59d59de-5d71-4c02-b397-ced141ec2994">

## Enhancements Achieved

1. **Faster Execution with Parallelism**: I have implemented goroutines to scan the files content of a commit parallelly.The same can be found in `switchAndScan()` function. Specifically, the `scanFileContent` method is called within a goroutine using the go keyword

2. **Base64 Decoding**: I have added `IsBase64Encoded()` function to first check if the keys are base64 encoded. If yes then `DecodeBase64()` function is used to decode them. Implementation of both these methods can be found in `utility.go` file

3. **Extension for Other Secret Detection**: I have extend support to validate other cloud credentials & for this purpose I am using `CredentialValidator` interface that has fields for finding and then validating the credentials. This assignment is very specific to aws and you can see its usage in `validate_aws.go` file. To extend this to other cloud providers,it would require us to create another file lets say `validate_gcp.go` and then use the interface mentioned earlier to write another defination of `FindCredentials()` & `ValidateCredentials()` field. This would allows us to validate gcp credentials as well without affecting the aws defination
