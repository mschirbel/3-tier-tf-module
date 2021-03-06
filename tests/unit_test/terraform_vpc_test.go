package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-network-example using Terratest.
func TestTerraformAwsNetwork(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../terraform/",

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
	publicSubnetIDs := terraform.Output(t, terraformOptions, "public_subnets_id")
	privateSubnetIDs := terraform.Output(t, terraformOptions, "private_subnets_id")
	vpcID := terraform.Output(t, terraformOptions, "main_vpc_id")

	// Declare a Replacer for the output string
	replacer := strings.NewReplacer("[", "", "]", "", "\"", "", "\n", "", " ", "")

	// Create an Array for the Subnets
	subnetPublID := replacer.Replace(publicSubnetIDs)
	arrayPublSubnets := strings.Split(subnetPublID, ",")

	subnetPrivID := replacer.Replace(privateSubnetIDs)
	arrayPrivSubnets := strings.Split(subnetPrivID, ",")

	// Test 1: Verify if the number of subnets is the same as declared
	require.Equal(t, 4, len(aws.GetSubnetsForVpc(t, vpcID, awsRegion)))

	// Test 2: Verify if the network that is supposed to be private is really private
	for i := 0; i < len(arrayPrivSubnets)-1; i++ {
		isPrivate, err := aws.IsPublicSubnetE(t, arrayPrivSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.False(t, isPrivate)
	}

	// Test 3: Verify if the network that is supposed to be public is really public
	for i := 0; i < len(arrayPublSubnets)-1; i++ {
		isPublic, err := aws.IsPublicSubnetE(t, arrayPublSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.True(t, isPublic)
	}

}
