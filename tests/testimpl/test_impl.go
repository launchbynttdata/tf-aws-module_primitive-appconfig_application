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

// TestComposableComplete verifies the deployed AppConfig application and exercises a reversible tag write.
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	client, arn := verifyApplication(t, ctx)
	exerciseTagWrite(t, client, arn)
}

// TestComposableCompleteReadOnly verifies the deployed AppConfig application using read-only AWS API calls.
func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	verifyApplication(t, ctx)
}

func verifyApplication(t *testing.T, ctx types.TestContext) (*appconfig.Client, string) {
	opts := ctx.TerratestTerraformOptions()
	region := terraform.Output(t, opts, "region")
	id := terraform.Output(t, opts, "id")
	arn := terraform.Output(t, opts, "arn")
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

	return client, arn
}

func appConfigClient(t *testing.T, region string) *appconfig.Client {
	t.Helper()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err)

	return appconfig.NewFromConfig(cfg)
}

func exerciseTagWrite(t *testing.T, client *appconfig.Client, resourceARN string) {
	t.Helper()

	const tagKey = "codex-functional-test"
	_, err := client.TagResource(context.Background(), &appconfig.TagResourceInput{
		ResourceArn: aws.String(resourceARN),
		Tags:        map[string]string{tagKey: "true"},
	})
	require.NoError(t, err)

	_, err = client.UntagResource(context.Background(), &appconfig.UntagResourceInput{
		ResourceArn: aws.String(resourceARN),
		TagKeys:     []string{tagKey},
	})
	require.NoError(t, err)
}
