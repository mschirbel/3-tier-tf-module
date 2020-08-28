package test

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"
	"encoding/json"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAlb(t *testing.T) {
	t.Parallel()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueId()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	mysqlVersion := fmt.Sprintf("5.7.19")

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	expectedResult := map[string]interface{}{
		"database_version": mysqlVersion,
		"region": awsRegion,
		"unique_id": uniqueID,
	}

	dataI, _ := json.Marshal(expectedResult)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../infrastructure/",

		// Variables to pass to our Terraform code using -var options
		Vars: expectedResult,
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceURL := terraform.Output(t, terraformOptions, "alb-dns")

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 10 * time.Second

	// Verify that we get back a 200 OK with the expected expectedResult
	http_helper.HttpGetWithRetry(t, instanceURL, nil, 200, string(dataI), maxRetries, timeBetweenRetries)
}