package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "286a69d3-dc09-45e3-b4a1-1b7dc9a02f90"


func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "buru0003vmTest",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	vmNicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Check if the NIC is available and connected to the specified VM
	virtualMachineNics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, vmNicName, virtualMachineNics[0])

	// Check if the VM is running the correct version of Ubuntu
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	vmOsSKU := "22_04-lts-gen2"
	vmOsPublisher := "Canonical"
	assert.Equal(t, vmOsSKU, vmImage.SKU)
	assert.Equal(t, vmOsPublisher, vmImage.Publisher)

}
