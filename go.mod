module github.com/ryeguard/ddbcalc

go 1.22.0

require (
	github.com/aws/aws-sdk-go v1.50.25
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.29.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.25.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.0 // indirect
	// This is an interesting dependency
	// "Smithy is an open-source Interface Definition Language (IDL) for web services created by AWS."
	github.com/aws/smithy-go v1.20.0 // indirect
	// Also interesting: 
	// " ... It will take a JSON document and transform it into another JSON document through a JMESPath expression."
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
