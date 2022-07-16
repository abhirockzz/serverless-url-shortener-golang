package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const shortCodeDynamoDBAttributeName = "shortcode"
const envVarName = "TABLE_NAME"
const pathParameter = "/{shortcode}"

const createShortURLFunctionDirectory = "../create-shortcode"
const accessURLFunctionDirectory = "../access-url"
const updateShortURLStatusFunctionDirectory = "../update-status"
const deleteShortURLFunctionDirectory = "../delete-shortcode"

type ServerlessURLShortenerStackProps struct {
	awscdk.StackProps
}

func NewServerlessURLShortenerStack(scope constructs.Construct, id string, props *ServerlessURLShortenerStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	dynamoDBTable := awsdynamodb.NewTable(stack, jsii.String("url-shortener-dynamodb-table"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String(shortCodeDynamoDBAttributeName),
				Type: awsdynamodb.AttributeType_STRING}})

	dynamoDBTable.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	urlShortenerAPI := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("url-shortner-http-api"), nil)

	funcEnvVar := &map[string]*string{envVarName: dynamoDBTable.TableName()}

	createURLFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("create-url-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Environment: funcEnvVar,
			Entry:       jsii.String(createShortURLFunctionDirectory)})

	dynamoDBTable.GrantWriteData(createURLFunction)

	createFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("create-function-integration"), createURLFunction, nil)

	urlShortenerAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: createFunctionIntg})

	acccessURLFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("access-url-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Environment: funcEnvVar,
			Entry:       jsii.String(accessURLFunctionDirectory)})

	dynamoDBTable.GrantReadData(acccessURLFunction)

	accessFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("access-function-integration"), acccessURLFunction, nil)

	urlShortenerAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String(pathParameter),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: accessFunctionIntg})

	updateURLStatusFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("update-url-status-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Environment: funcEnvVar,
			Entry:       jsii.String(updateShortURLStatusFunctionDirectory)})

	dynamoDBTable.GrantWriteData(updateURLStatusFunction)

	updateFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("update-function-integration"), updateURLStatusFunction, nil)

	urlShortenerAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String(pathParameter),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_PUT},
		Integration: updateFunctionIntg})

	deleteURLStatusFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("delete-url-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Environment: funcEnvVar,
			Entry:       jsii.String(deleteShortURLFunctionDirectory)})

	dynamoDBTable.GrantWriteData(deleteURLStatusFunction)

	deleteFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("update-function-integration"), deleteURLStatusFunction, nil)

	urlShortenerAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String(pathParameter),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_DELETE},
		Integration: deleteFunctionIntg})

	awscdk.NewCfnOutput(stack, jsii.String("output"), &awscdk.CfnOutputProps{Value: urlShortenerAPI.Url(), Description: jsii.String("API Gateway endpoint")})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewServerlessURLShortenerStack(app, "ServerlessURLShortenerStack", &ServerlessURLShortenerStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
