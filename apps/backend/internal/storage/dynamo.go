package storage

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var Client *dynamodb.Client

func Initialize(ctx context.Context) {
	awsEndpoint := "http://localstack:4566"
	awsRegion := "us-east-1"

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "test",
				SecretAccessKey: "test",
				SessionToken:    "",
			},
		}),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	Client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
		o.Region = awsRegion
	})
}

func GetUpdateExpression(item any, attrNames map[string]string, attrVals map[string]types.AttributeValue) (string, error) {
	var updateParts []string
	v := reflect.ValueOf(item).Elem()
	tType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := tType.Field(i)
		tag := field.Tag.Get("updateable")
		if tag == "true" {
			ddbName := field.Tag.Get("dynamodbav")
			if ddbName == "" {
				ddbName = field.Tag.Get("json")
			}

			attrNames[fmt.Sprintf("#%s", ddbName)] = ddbName
			attrVal, err := attributevalue.Marshal(v.Field(i).Interface())
			if err != nil {
				return "", err
			}
			attrVals[fmt.Sprintf(":%s", ddbName)] = attrVal
			updateParts = append(updateParts, fmt.Sprintf("#%s = :%s", ddbName, ddbName))
		}
	}

	updateExpression := "SET " + strings.Join(updateParts, ", ")
	return updateExpression, nil
}

func GetFilterExpression(queryParams map[string][]string, attrNames map[string]string, attrVals map[string]types.AttributeValue) string {
	var filterParts []string
	for key, values := range queryParams {
		for i, value := range values {
			attrName := fmt.Sprintf("#%s", key)
			attrValue := fmt.Sprintf(":%s%d", key, i)
			attrNames[attrName] = key
			attrVals[attrValue] = &types.AttributeValueMemberS{Value: value}
			filterParts = append(filterParts, fmt.Sprintf("%s = %s", attrName, attrValue))
		}
	}
	filterExpression := strings.Join(filterParts, " AND ")
	return filterExpression
}
