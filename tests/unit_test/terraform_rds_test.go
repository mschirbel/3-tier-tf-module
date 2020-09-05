package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-network-example using Terratest.
func TestTerraformAwsRDS(t *testing.T) {
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
	RDSInstanceID := terraform.Output(t, terraformOptions, "rds_id")
	connectionStringParameter := terraform.Output(t, terraformOptions, "rds_connection_string_parameter")

	// Defines which SSM Parameter contains the connection string and get its value
	key, err := aws.GetParameterE(t, awsRegion, connectionStringParameter)
	if err != nil {
		fmt.Printf("Error Encountered in getting AWS Parameter: %s\n", err)
		return
	}
	// Convert the String output to Json
	connToJSON := []byte(key)
	var JSONMapConnString map[string]interface{}
	if err := json.Unmarshal(connToJSON, &JSONMapConnString); err != nil {
		fmt.Printf("Error Encountered in Unmarshal: %s\n", err)
		return
	}

	// Values expected in the Test Result
	expectedPort := int64(3306)
	expectedDatabaseName := fmt.Sprint(JSONMapConnString["DATABASE"])
	username := fmt.Sprint(JSONMapConnString["USER"])
	password := fmt.Sprint(JSONMapConnString["PASS"])

	// Define values to test
	address, err := aws.GetAddressOfRdsInstanceE(t, RDSInstanceID, awsRegion)
	if err != nil {
		fmt.Printf("Error Encountered in getting RDS Address: %s\n", err)
		return
	}
	port, err := aws.GetPortOfRdsInstanceE(t, RDSInstanceID, awsRegion)
	if err != nil {
		fmt.Printf("Error Encountered in getting RDS Port: %s\n", err)
		return
	}
	schemaExistsInRdsInstance, err := aws.GetWhetherSchemaExistsInRdsMySqlInstanceE(t, address, port, username, password, expectedDatabaseName)
	if err != nil {
		fmt.Printf("Error Encountered in getting RDS Schema: %s\n", err)
		return
	}

	// Verify that the address is not null
	assert.NotNil(t, address)
	// Verify that the DB instance is listening on the port mentioned
	assert.Equal(t, expectedPort, port)
	// Verify that the table/schema requested for creation is actually present in the database
	assert.True(t, schemaExistsInRdsInstance)
}
