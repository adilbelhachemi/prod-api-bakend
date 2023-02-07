package storage

import (
	"errors"
	"fmt"
	"pratbacknd/internal/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var ErrorNotFound = errors.New("Item not found")

func (d *Dynamo) UpdateInventory(productId string, delta int) error {
	p, err := d.getProductById(productId)
	if err != nil {
		return fmt.Errorf("error - to retrieve product: %w", err)
	}

	newStock := int(p.Stock) + delta
	if newStock < 0 {
		return fmt.Errorf("error - stock should not be less than 0")
	}

	keyCondition := make(map[string]*dynamodb.AttributeValue)
	// PK
	keyCondition[PartitionKeyAttributeName] = &dynamodb.AttributeValue{S: aws.String(pkProduct)}
	// SK
	keyCondition[SortkeyAttributeName] = &dynamodb.AttributeValue{S: aws.String(productId)}

	// Condition experession
	condition := expression.Name("version").Equal(expression.Value(p.Version))

	// update
	update := expression.Set(expression.Name("stock"), expression.Value(newStock))
	update.Set(expression.Name("version"), expression.Value(p.Version+1))

	// build the expression with expression builder
	builder := expression.NewBuilder().WithCondition(condition).WithUpdate(update)
	expr, err := builder.Build()
	if err != nil {
		return fmt.Errorf("error - building the expression: %w", err)
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 &d.tableName,
		Key:                       keyCondition,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = d.client.UpdateItem(&input)
	if err != nil {
		return fmt.Errorf("error - run update item request: %w", err)
	}

	return nil
}

func (d *Dynamo) getProductById(productID string) (types.Product, error) {
	getItemInput := dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		Key: map[string]*dynamodb.AttributeValue{
			PartitionKeyAttributeName: {
				S: aws.String(pkProduct),
			},
			SortkeyAttributeName: {
				S: aws.String(productID),
			},
		},
		TableName: &d.tableName,
	}

	out, err := d.client.GetItem(&getItemInput)
	if err != nil {
		return types.Product{}, fmt.Errorf("error - geeting item: %w", err)
	}

	if len(out.Item) == 0 {
		return types.Product{}, ErrorNotFound
	}

	var p types.Product
	err = dynamodbattribute.UnmarshalMap(out.Item, &p)
	if err != nil {
		return p, fmt.Errorf("error - marchalling product: %s", err)
	}

	return p, nil
}
