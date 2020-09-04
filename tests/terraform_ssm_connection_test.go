package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAwsSsm(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../infrastructure/",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"region": awsRegion,
		},
	}
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceOutput := terraform.Output(t, terraformOptions, "instances_ip")

	// Creating an array of instances Ids from ASG
	replacer := strings.NewReplacer("[", "", "]", "", "\"", "", "\n", "", " ", "", "{", "", "}", "")
	instancesIds := replacer.Replace(instanceOutput)
	arrayInstances := strings.Split(instancesIds, ",")

	// Define Timeout For Calling SSM Agent on Instance
	timeout := 4 * time.Minute

	// Wait for SSM Instance Agent
	for i := 0; i < len(arrayInstances); i++ {
		aws.WaitForSsmInstance(t, awsRegion, arrayInstances[i], timeout)
	}

	// First Test: Check if Command gets Stdout
	for i := 0; i < len(arrayInstances); i++ {
		result, err := aws.CheckSsmCommandE(t, awsRegion, arrayInstances[i], "echo Hello, World", timeout)
		fmt.Printf("Checking instance %s", arrayInstances[i])
		if err != nil {
			fmt.Printf("Error Encountered in checking SSM Stdout: %s\n", err)
			return
		}
		assert.Equal(t, result.Stdout, "Hello, World\n")
		assert.Equal(t, result.Stderr, "")
		assert.Equal(t, int64(0), result.ExitCode)
	}

	// Second Test: Check if Command gets Stderr
	for i := 0; i < len(arrayInstances); i++ {
		result, err := aws.CheckSsmCommandE(t, awsRegion, arrayInstances[i], "cat /wrong/file", timeout)
		fmt.Printf("Checking instance %s", arrayInstances[i])
		if err != nil {
			fmt.Printf("Error Encountered in checking SSM Stderr: %s\n", err)
			return
		}
		assert.Equal(t, "Failed", err.Error())
		assert.Equal(t, "cat: /wrong/file: No such file or directory\nfailed to run commands: exit status 1", result.Stderr)
		assert.Equal(t, "", result.Stdout)
		assert.Equal(t, int64(1), result.ExitCode)
	}
}
