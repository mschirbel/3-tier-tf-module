package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// An example of how to test the Terraform module in examples/terraform-aws-network-example using Terratest.
func TestTerraformAwsRDS(t *testing.T) {
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
	RDSInstanceID := terraform.Output(t, terraformOptions, "rds_id")
	connectionStringParameter := terraform.Output(t, terraformOptions, "rds_connection_string_parameter")

	// Defines which SSM Parameter contains the connection string
	keyName, err := aws.GetParameterE(t, awsRegion, connectionStringParameter)
	if err != nil {
		fmt.Printf("Error Encountered: %s", err)
		return
	}
	fmt.Println(keyName)
	fmt.Println(RDSInstanceID)
}
