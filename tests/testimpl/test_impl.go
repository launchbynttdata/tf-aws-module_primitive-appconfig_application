package testimpl

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComposableComplete verifies the deployed AppConfig application.
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	verifyApplication(t, ctx)
}

// TestComposableCompleteReadOnly verifies the deployed AppConfig application using read-only AWS API calls.
func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	verifyApplication(t, ctx)
}

func verifyApplication(t *testing.T, ctx types.TestContext) {
	opts := ctx.TerratestTerraformOptions()
	region := terraform.Output(t, opts, "region")
	id := terraform.Output(t, opts, "id")
	name := terraform.Output(t, opts, "name")
	description := terraform.Output(t, opts, "description")

	require.NotEqual(t, "", id)
	assert.Equal(t, terraform.Output(t, opts, "expected_name"), name)
	assert.Equal(t, terraform.Output(t, opts, "expected_description"), description)

	client := appConfigClient(t, region)
	app, err := client.GetApplication(context.Background(), &appconfig.GetApplicationInput{ApplicationId: aws.String(id)})
	require.NoError(t, err)

	assert.Equal(t, id, aws.ToString(app.Id))
	assert.Equal(t, name, aws.ToString(app.Name))
	assert.Equal(t, description, aws.ToString(app.Description))
}

func appConfigClient(t *testing.T, region string) *appconfig.Client {
	t.Helper()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err)

	return appconfig.NewFromConfig(cfg)
}
