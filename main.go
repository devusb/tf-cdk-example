package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/s3bucket"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/s3bucketversioning"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

// Config represents what the developer writes
type Config struct {
	Project     string        `json:"project"`
	Environment string        `json:"environment"`
	Region      string        `json:"region"`
	Storage     StorageConfig `json:"storage"`
}

type StorageConfig struct {
	BucketName        string `json:"bucket_name"`
	EnableVersioning  bool   `json:"enable_versioning"`
}

func main() {
	// Step 1: Read the JSON config file
	fmt.Println("📄 Reading config.json...")
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Error reading config.json: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Parse the JSON into our struct
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Printf("Error parsing config.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Config loaded for project: %s (environment: %s)\n\n",
		config.Project, config.Environment)

	// Step 3: Create CDKTF app
	fmt.Println("🏗️  Creating infrastructure from config...")
	app := cdktf.NewApp(nil)

	// Step 4: Create a stack
	stackName := fmt.Sprintf("%s-%s-stack", config.Project, config.Environment)
	stack := cdktf.NewTerraformStack(app, jsii.String(stackName))

	// Step 5: Add AWS provider (from config)
	provider.NewAwsProvider(stack, jsii.String("aws"), &provider.AwsProviderConfig{
		Region: jsii.String(config.Region),
	})

	// Step 6: Create S3 bucket based on config
	fullBucketName := fmt.Sprintf("%s-%s-%s",
		config.Project, config.Environment, config.Storage.BucketName)

	bucket := s3bucket.NewS3Bucket(stack, jsii.String("bucket"), &s3bucket.S3BucketConfig{
		Bucket: jsii.String(fullBucketName),
		Tags: &map[string]*string{
			"Project":     jsii.String(config.Project),
			"Environment": jsii.String(config.Environment),
			"ManagedBy":   jsii.String("CDKTF-JSON-Platform"),
		},
	})

	// Step 7: Add versioning if requested
	if config.Storage.EnableVersioning {
		s3bucketversioning.NewS3BucketVersioningA(stack, jsii.String("versioning"),
			&s3bucketversioning.S3BucketVersioningAConfig{
				Bucket: bucket.Bucket(),
				VersioningConfiguration: &s3bucketversioning.S3BucketVersioningVersioningConfiguration{
					Status: jsii.String("Enabled"),
				},
			})
		fmt.Println("  ✓ S3 Bucket with versioning enabled")
	} else {
		fmt.Println("  ✓ S3 Bucket (no versioning)")
	}

	// Step 8: Add outputs
	cdktf.NewTerraformOutput(stack, jsii.String("bucket_name"), &cdktf.TerraformOutputConfig{
		Value:       bucket.Bucket(),
		Description: jsii.String("The name of the created S3 bucket"),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("bucket_arn"), &cdktf.TerraformOutputConfig{
		Value:       bucket.Arn(),
		Description: jsii.String("The ARN of the created S3 bucket"),
	})

	// Step 9: Synthesize to Terraform JSON
	fmt.Println("\n📝 Synthesizing to Terraform JSON...")
	app.Synth()

	fmt.Println("✓ Done!")
	fmt.Printf("\n📁 Generated Terraform in: cdktf.out/stacks/%s/\n", stackName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Review: cat cdktf.out/stacks/%s/cdk.tf.json\n", stackName)
	fmt.Println("  2. Deploy: cd cdktf.out/stacks/" + stackName + " && terraform init && terraform apply")
}
