package storage

import (
	"fmt"
	"pratbacknd/internal/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	PartitionKeyAttributeName = "PK"
	SortkeyAttributeName      = "SK"
	pkProduct                 = "product"
	pkCart                    = "cart"
	pkCategory                = "category"
)

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

func (d *Dynamo) CreateProduct(p types.Product) error {
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

func (d *Dynamo) Products() ([]types.Product, error) {
	out, err := d.getElementByPkAndSk(pkProduct, "")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &products)
	if err != nil {
		return nil, fmt.Errorf("error - Unmarshalling results: %w", err)
	}
	return products, nil
}

func (d *Dynamo) CreateCategory(c types.Category) error {
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

func (d *Dynamo) Categories() ([]types.Category, error) {
	out, err := d.getElementByPkAndSk(pkCategory, "")
	if err != nil {
		return nil, err
	}

	categories := make([]types.Category, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &categories)
	if err != nil {
		return nil, fmt.Errorf("error - Unmarshalling results: %w", err)
	}
	return categories, nil
}

func (d *Dynamo) getElementByPkAndSk(pkAttributeValue, skAttributeValue string) (*dynamodb.QueryOutput, error) {
	keyCondition := expression.Key(PartitionKeyAttributeName).Equal(expression.Value(pkAttributeValue))

	if skAttributeValue != "" {
		sortKeyCondition := expression.Key(SortkeyAttributeName).Equal(expression.Value(skAttributeValue))
		keyCondition = keyCondition.And(sortKeyCondition)
	}

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

func (d *Dynamo) CreateCart(cart types.Cart, userId string) error {
	item, err := dynamodbattribute.MarshalMap(cart)
	if err != nil {
		return fmt.Errorf("error - marshal product: %w", err)
	}

	item[PartitionKeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(pkCart),
	}
	item[SortkeyAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(userId),
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error - Put Cart in db: %w", err)
	}
	return nil
}

func (d *Dynamo) GetCart(userID string) (types.Cart, error) {

	out, err := d.getElementByPkAndSk(pkCart, userID)
	if err != nil {
		return types.Cart{}, fmt.Errorf("error - retreiving Cart in db: %w", err)
	}

	if len(out.Items) == 0 {
		return types.Cart{}, fmt.Errorf("error - no cart found: %w", ErrorNotFound)
	}

	if len(out.Items) > 1 {
		return types.Cart{}, fmt.Errorf("error - more than one cart found: %w", err)
	}

	var c types.Cart
	err = dynamodbattribute.UnmarshalMap(out.Items[0], &c)
	if err != nil {
		return types.Cart{}, fmt.Errorf("error - Unmarshalling cart: %w", err)
	}

	return c, nil
}
