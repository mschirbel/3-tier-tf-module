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

	// Give the VPC and the subnets some CIDRs for testing
	vpcCidr := "10.0.0.0/16"
	privateSubnetCidr := "[\"10.0.1.0/24\", \"10.0.2.0/24\"]"
	publicSubnetCidr := "[\"10.0.101.0/24\", \"10.0.102.0/24\"]"

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../infrastructure/",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"main_vpc_cidr":        vpcCidr,
			"private_subnets_cidr": privateSubnetCidr,
			"public_subnets_cidr":  publicSubnetCidr,
			"region":               awsRegion,
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
		fmt.Printf("Realizando o teste na %s\n", arrayPrivSubnets[i])
		_, err := aws.IsPublicSubnetE(t, arrayPrivSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.False(t, aws.IsPublicSubnet(t, arrayPrivSubnets[i], awsRegion))
	}

	// Test 3: Verify if the network that is supposed to be public is really public
	for i := 0; i < len(arrayPublSubnets)-1; i++ {
		fmt.Printf("Realizando o teste na %s\n", arrayPublSubnets[i])
		_, err := aws.IsPublicSubnetE(t, arrayPublSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.True(t, aws.IsPublicSubnet(t, arrayPublSubnets[i], awsRegion))
	}

}
