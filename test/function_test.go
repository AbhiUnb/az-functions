package test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAzureFunctionApp(t *testing.T) {
	t.Parallel()

	terraformDir := os.Getenv("TF_DIR")
	if terraformDir == "" {
		terraformDir = "../"
	}

	testName := os.Getenv("TEST_NAME")
	if testName == "" {
		testName = "Terratest"
	}

	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		functionName = "httpexample" // enforce lowercase to match deployment
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
	})

	t.Log("Initializing Terraform...")
	terraform.Init(t, terraformOptions)

	// Optional: If you want to reapply resources uncomment below
	// terraform.InitAndApply(t, terraformOptions)

	// Destroy at the end if needed
	// defer terraform.Destroy(t, terraformOptions)

	var functionAppURL string = os.Getenv("FUNCTION_APP_URL")
	if functionAppURL == "" {
		t.Log("FUNCTION_APP_URL env var not set. Trying Terraform output...")
		functionAppURL = terraform.Output(t, terraformOptions, "function_app_url")
		if functionAppURL == "" {
			t.Log("FUNCTION_APP_URL not found via env var or Terraform output. Skipping tests.")
			t.Skip("Set FUNCTION_APP_URL environment variable or ensure Terraform output exists.")
		}
	}

	// 1. Validate Function App existence
	t.Run("ValidateFunctionAppExistence", func(t *testing.T) {
		t.Log("Starting ValidateFunctionAppExistence test...")
		assert.NotEmpty(t, functionAppURL, "Function App URL should not be empty")
		t.Log("✅ Passed: ValidateFunctionAppExistence")
	})

	// 2. Validate hosting plan and SKU
	t.Run("ValidateHostingPlanSKU", func(t *testing.T) {
		t.Log("Starting ValidateHostingPlanSKU test...")
		sku := os.Getenv("FUNCTION_APP_SKU")
		if sku == "" {
			sku = terraform.Output(t, terraformOptions, "function_app_sku")
			if sku == "" {
				t.Skip("Skipping ValidateHostingPlanSKU as FUNCTION_APP_SKU env var or Terraform output is not set")
			}
		}
		assert.Contains(t, []string{"Y1", "EP1", "P1v2"}, sku, "Function App SKU should be Consumption, Premium, or Dedicated")
		t.Log("✅ Passed: ValidateHostingPlanSKU")
	})

	// 3. Validate deployment source
	t.Run("ValidateDeploymentSource", func(t *testing.T) {
		t.Log("Starting ValidateDeploymentSource test...")
		deploymentSource := os.Getenv("FUNCTION_APP_DEPLOYMENT_SOURCE")
		if deploymentSource == "" {
			deploymentSource = terraform.Output(t, terraformOptions, "function_app_deployment_source")
			if deploymentSource == "" {
				t.Skip("Skipping ValidateDeploymentSource as FUNCTION_APP_DEPLOYMENT_SOURCE env var or Terraform output is not set")
			}
		}
		assert.NotEmpty(t, deploymentSource, "Deployment source should be configured for CI/CD")
		t.Log("✅ Passed: ValidateDeploymentSource")
	})

	// 4. Validate App Settings
	t.Run("ValidateAppSettings", func(t *testing.T) {
		t.Log("Starting ValidateAppSettings test...")
		runtime := os.Getenv("FUNCTION_APP_RUNTIME")
		timeout := os.Getenv("FUNCTION_APP_TIMEOUT")
		if runtime == "" || timeout == "" {
			runtime = terraform.Output(t, terraformOptions, "function_app_runtime")
			timeout = terraform.Output(t, terraformOptions, "function_app_timeout")
			if runtime == "" || timeout == "" {
				t.Skip("Skipping ValidateAppSettings as FUNCTION_APP_RUNTIME or FUNCTION_APP_TIMEOUT env vars or Terraform outputs are not set")
			}
		}
		assert.Contains(t, runtime, "node", "Function App runtime should be Node.js")
		assert.NotEmpty(t, timeout, "Function App timeout setting should exist")
		t.Log("✅ Passed: ValidateAppSettings")
	})

	// 5. Validate trigger behavior (HTTP call)
	t.Run("ValidateTriggerBehavior", func(t *testing.T) {
		t.Log("Starting ValidateTriggerBehavior test...")
		testURL := functionAppURL
		if !strings.HasPrefix(testURL, "http") {
			testURL = "https://" + testURL
		}
		testURL += "/api/" + functionName + "?name=" + testName

		resp, err := http.Get(testURL)
		assert.NoError(t, err)
		if err != nil {
			t.FailNow()
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.Contains(string(body), "Hello, "+testName+"!"))
		t.Log("✅ Passed: ValidateTriggerBehavior")
	})

	// 6. Validate Application Insights integration
	t.Run("ValidateAppInsightsIntegration", func(t *testing.T) {
		t.Log("Starting ValidateAppInsightsIntegration test...")
		appInsightsKey := os.Getenv("FUNCTION_APP_APPINSIGHTS_KEY")
		if appInsightsKey == "" {
			appInsightsKey = terraform.Output(t, terraformOptions, "function_app_appinsights_key")
			if appInsightsKey == "" {
				t.Skip("Skipping ValidateAppInsightsIntegration as FUNCTION_APP_APPINSIGHTS_KEY env var or Terraform output is not set")
			}
		}
		assert.NotEmpty(t, appInsightsKey, "Application Insights instrumentation key should exist")
		t.Log("✅ Passed: ValidateAppInsightsIntegration")
	})
}

// TestNodeAppDeployment validates that a Node.js app is deployed in the Azure Function
func TestNodeAppDeployment(t *testing.T) {
	t.Parallel()

	terraformDir := os.Getenv("TF_DIR")
	if terraformDir == "" {
		terraformDir = "../"
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
	})

	t.Log("Initializing Terraform...")
	terraform.Init(t, terraformOptions)

	functionAppURL := os.Getenv("FUNCTION_APP_URL")
	if functionAppURL == "" {
		t.Log("FUNCTION_APP_URL env var not set. Trying Terraform output...")
		functionAppURL = terraform.Output(t, terraformOptions, "function_app_url")
		if functionAppURL == "" {
			t.Fatal("FUNCTION_APP_URL not found via env var or Terraform output")
		}
	}

	// Ensure URL has scheme
	if !strings.HasPrefix(functionAppURL, "http") {
		functionAppURL = "https://" + functionAppURL
	}

	fullURL := functionAppURL + "/api/httpexample?name=TestDeployment"
	t.Logf("Testing Node.js app deployment at %s", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	assert.NoError(t, err, "Failed to create HTTP request")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to send HTTP request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read HTTP response body")

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
	assert.Contains(t, string(body), "Hello, TestDeployment!", "Node.js function should respond correctly")

	t.Log("✅ Passed: TestNodeAppDeployment")
}
