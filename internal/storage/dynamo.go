package storage

import (
	"fmt"
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const PartitionKeyAttributeName = "PK"
const SortkeyAttributeName = "SK"
const pkProduct = "product"
const pkCategory = "category"

type Dynamo struct {
	tableName  string
	awsSession *session.Session
	client     *dynamodb.DynamoDB
}

func NewDynamo(tableName string) (*Dynamo, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("error - creating aws session: %w", err)
	}
	dynamodbClient := dynamodb.New(awsSession)
	return &Dynamo{
		tableName:  tableName,
		awsSession: awsSession,
		client:     dynamodbClient,
	}, nil
}

func (d *Dynamo) CreateProduct(p product.Product) error {
	item, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return fmt.Errorf("error - marshal product: %w", err)
	}

	item[PartitionKeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(pkProduct),
	}
	item[SortkeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(p.ID),
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error - Put item in db: %w", err)
	}
	return nil
}

func (d *Dynamo) Products() ([]product.Product, error) {
	out, err := d.getElementByPk(pkProduct)
	if err != nil {
		return nil, err
	}

	products := make([]product.Product, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &products)
	if err != nil {
		return nil, fmt.Errorf("error - Unmarshalling results: %w", err)
	}
	return products, nil
}

func (d *Dynamo) CreateCategory(c category.Category) error {
	item, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return fmt.Errorf("error - marshal category: %w", err)
	}

	item[PartitionKeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(pkCategory),
	}
	item[SortkeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(c.ID),
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error - Put item in db: %w", err)
	}
	return nil
}

func (d *Dynamo) Categories() ([]category.Category, error) {
	out, err := d.getElementByPk(pkCategory)
	if err != nil {
		return nil, err
	}

	categories := make([]category.Category, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &categories)
	if err != nil {
		return nil, fmt.Errorf("error - Unmarshalling results: %w", err)
	}
	return categories, nil
}

func (d *Dynamo) getElementByPk(pkAttributeValue string) (*dynamodb.QueryOutput, error) {
	keyCondition := expression.Key(PartitionKeyAttributeName).Equal(expression.Value(pkAttributeValue))
	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("error - building expression: %w", err)
	}

	input := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 &d.tableName,
	}

	out, err := d.client.Query(&input)
	if err != nil {
		return nil, fmt.Errorf("error - building expression: %w", err)
	}

	return out, nil
}

func (d *Dynamo) UpdateProduct(input UpdateProductInput) error {
	p, err := d.getProductById(input.ProductId)
	if err != nil {
		return fmt.Errorf("error - to retrieve product: %w", err)
	}

	// key condition
	keyCondition := make(map[string]*dynamodb.AttributeValue)
	// PK
	keyCondition[PartitionKeyAttributeName] = &dynamodb.AttributeValue{S: aws.String(pkProduct)}
	// SK
	keyCondition[SortkeyAttributeName] = &dynamodb.AttributeValue{S: aws.String(input.ProductId)}

	// condition expression
	condition := expression.Name("version").Equal(expression.Value(p.Version))

	// update expression
	update := expression.Set(expression.Name("name"), expression.Value(input.Name))
	update.Set(expression.Name("version"), expression.Value(p.Version+1))
	update.Set(expression.Name("image"), expression.Value(input.Image))
	update.Set(expression.Name("shortDescription"), expression.Value(input.ShortDescription))
	update.Set(expression.Name("priceVatExcluded"), expression.Value(input.PriceVATExcluded))
	update.Set(expression.Name("vat"), expression.Value(input.VAT))
	update.Set(expression.Name("totalPrice"), expression.Value(input.TotalPrice))

	// build the expression with expression builder
	builder := expression.NewBuilder().WithCondition(condition).WithUpdate(update)
	expr, err := builder.Build()
	if err != nil {
		return fmt.Errorf("error - building the expression: %w", err)
	}

	// request UpdateItem
	item := dynamodb.UpdateItemInput{
		TableName:                 &d.tableName,
		Key:                       keyCondition,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = d.client.UpdateItem(&item)
	if err != nil {
		return fmt.Errorf("error - run update item request: %w", err)
	}

	return nil
}
