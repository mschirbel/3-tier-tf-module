package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformModules(t *testing.T) {
	t.Parallel()

	// The folder where we have our Terraform code
	workingDir := "../terraform/"

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueId()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// At the end of the test, undeploy the web app using Terraform
	defer test_structure.RunTestStage(t, "cleanup", func() {
		destroyTerraform(t, workingDir)
	})

	// Stage to Deploy Terraform
	test_structure.RunTestStage(t, "deploy", func() {
		awsRegion := test_structure.LoadString(t, workingDir, "awsRegion")
		deployUsingTerraform(t, awsRegion, workingDir, uniqueID)
	})

	// Validate that subnets are public and private
	test_structure.RunTestStage(t, "validate", func() {
		validateVpc(t, workingDir, awsRegion)
	})

	// Validate that RDS is available
	test_structure.RunTestStage(t, "validate", func() {
		validateRds(t, workingDir, awsRegion)
	})

	// Validate that SSM is responding
	test_structure.RunTestStage(t, "validate", func() {
		validateSsm(t, workingDir, awsRegion)
	})

	// Validate that the web app deployed and is responding to HTTP requests
	test_structure.RunTestStage(t, "validate", func() {
		validateAlb(t, workingDir, awsRegion, uniqueID)
	})

}

// Undeploy the app using Terraform
func destroyTerraform(t *testing.T, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	terraform.Destroy(t, terraformOptions)
}

// Deploy the app using Terraform
func deployUsingTerraform(t *testing.T, awsRegion string, workingDir string, uniqueID string) {

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: workingDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"region":    awsRegion,
			"unique_id": uniqueID,
		},
	}

	// Save the Terraform Options struct, instance name, and instance text so future test stages can use it
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)
}

func validateVpc(t *testing.T, workingDir string, awsRegion string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Defines the MySQL Version in the test
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
		_, err := aws.IsPublicSubnetE(t, arrayPrivSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.False(t, aws.IsPublicSubnet(t, arrayPrivSubnets[i], awsRegion))
	}

	// Test 3: Verify if the network that is supposed to be public is really public
	for i := 0; i < len(arrayPublSubnets)-1; i++ {
		_, err := aws.IsPublicSubnetE(t, arrayPublSubnets[i], awsRegion)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
		assert.True(t, aws.IsPublicSubnet(t, arrayPublSubnets[i], awsRegion))
	}
}

func validateRds(t *testing.T, workingDir string, awsRegion string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

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

func validateSsm(t *testing.T, workingDir string, awsRegion string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

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

func validateAlb(t *testing.T, workingDir string, awsRegion string, uniqueID string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Defines the MySQL Version in the test
	mysqlVersion := fmt.Sprintf("5.7.19")

	// Expected Result from the HTTP Request in the ALB
	httpJSON := map[string]interface{}{
		"database_version": fmt.Sprintf("%s-log", mysqlVersion),
		"region":           awsRegion,
		"unique_id":        uniqueID,
	}
	// Format the Expected Result to JSON
	expectedResult, _ := json.Marshal(httpJSON)

	// Run `terraform output` to get the value of an output variable
	albDNS := terraform.Output(t, terraformOptions, "alb-dns")

	// Formats the alb dns to match HTTP request valid url
	URL := fmt.Sprintf("http://%s:80", albDNS)

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 10 * time.Second

	// Verify that we get back a 200 OK with the expected expectedResult
	http_helper.HttpGetWithRetry(t, URL, nil, 200, string(expectedResult), maxRetries, timeBetweenRetries)
}
