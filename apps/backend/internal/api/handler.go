package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/internal/storage"
	"github.com/Taiterbase/vtrips/apps/backend/internal/storage/index"
	"github.com/Taiterbase/vtrips/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func CreateTrip(c echo.Context) error {
	var (
		trip     = models.NewTrip()
		clientID = c.QueryParam("client_id")
	)

	trip.SetClientID(clientID)
	c.Logger().Debugj(log.JSON{"trip": trip})
	err := c.Bind(&trip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}
	c.Logger().Debugj(log.JSON{"trip": trip})
	if err = trip.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	item, err := attributevalue.MarshalMap(&trip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	item["pk"] = &types.AttributeValueMemberS{Value: index.MakePK(trip.GetClientID())}
	item["sk"] = &types.AttributeValueMemberS{Value: index.MakeSK(trip.GetID())}
	tx := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String("vtrips"),
					Item:      item,
				},
			},
		},
	}

	for _, action := range index.GetTripWriteActions() {
		action.Add(tx, trip)
	}

	_, err = storage.Client.TransactWriteItems(c.Request().Context(), tx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := make([]models.TripBase, 0)
	for _, item := range tx.TransactItems {
		var (
			// assumes all the transactions are on trips
			tripOut models.TripBase
			err     error
		)
		if item.Put != nil {
			err = attributevalue.UnmarshalMap(item.Put.Item, &tripOut)
		} else if item.Update != nil {
			err = attributevalue.UnmarshalMap(item.Update.Key, &tripOut)
		} else if item.Delete != nil {
			err = attributevalue.UnmarshalMap(item.Delete.Key, &tripOut)
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		response = append(response, tripOut)
	}

	return c.JSON(http.StatusOK, response)
}

func GetTrip(c echo.Context) error {
	var (
		clientID = c.QueryParam("client_id")
		tripID   = c.Param("trip_id")
	)
	trip, err := getTrip(c.Request().Context(), clientID, tripID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, trip)
}

func getTrip(ctx context.Context, clientID, tripID string) (*models.TripBase, error) {
	item, err := storage.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("vtrips"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: index.MakePK(clientID)},
			"sk": &types.AttributeValueMemberS{Value: index.MakeSK(tripID)},
		},
	})
	if err != nil {
		return nil, err
	}
	var trip models.TripBase
	err = attributevalue.UnmarshalMap(item.Item, &trip)
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func UpdateTrip(c echo.Context) error {
	var (
		tripID    = c.Param("trip_id")
		clientID  = c.QueryParam("client_id")
		attrNames = map[string]string{}
		attrVals  = map[string]types.AttributeValue{}
	)
	trip, err := getTrip(c.Request().Context(), clientID, tripID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err = c.Bind(&trip); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	trip.SetUpdatedAt(time.Now().UnixMilli())
	updateExpr, err := storage.GetUpdateExpression(trip, attrNames, attrVals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	tx := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Update: &types.Update{
					TableName: aws.String("vtrips"),
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: index.MakePK(trip.GetClientID())},
						"sk": &types.AttributeValueMemberS{Value: index.MakeSK(trip.GetID())},
					},
					UpdateExpression:                    aws.String(updateExpr),
					ConditionExpression:                 aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
					ExpressionAttributeNames:            attrNames,
					ExpressionAttributeValues:           attrVals,
					ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
		},
	}

	for _, action := range index.GetTripWriteActions() {
		action.Update(tx, trip, nil)
	}

	_, err = storage.Client.TransactWriteItems(c.Request().Context(), tx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trip)
}

func DeleteTrip(c echo.Context) error {
	var (
		clientID = c.QueryParam("client_id")
		tripID   = c.Param("trip_id")
	)

	tx := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					TableName: aws.String("vtrips"),
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: index.MakePK(clientID)},
						"sk": &types.AttributeValueMemberS{Value: index.MakeSK(tripID)},
					},
					ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
				},
			},
		},
	}
	_, err := storage.Client.TransactWriteItems(c.Request().Context(), tx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tripID)
}

func ListTrips(c echo.Context) error {
	var (
		clientID  = c.QueryParam("client_id")
		attrNames = map[string]string{
			"#pk": "pk",
			"#sk": "sk",
		}
		attrVals = map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: index.MakePK(clientID)},
			":sk": &types.AttributeValueMemberS{Value: "trip_id$"},
		}
		query = &dynamodb.QueryInput{
			TableName:                 aws.String("vtrips"),
			ExpressionAttributeNames:  attrNames,
			ExpressionAttributeValues: attrVals,
			Limit:                     aws.Int32(10),
			ScanIndexForward:          aws.Bool(false),
			KeyConditionExpression:    aws.String("#pk = :pk AND #sk < :sk"),
		}
		trips       []models.TripBase
		queryParams = c.Request().URL.Query()
	)
	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		query.Limit = aws.Int32(int32(limit))
	}
	if page := c.QueryParam("page"); page != "" {
		attrVals[":sk"] = &types.AttributeValueMemberS{Value: page}
	}

	filterExpr := storage.GetFilterExpression(queryParams, attrNames, attrVals)
	if filterExpr != "" {
		query.FilterExpression = aws.String(filterExpr)
	}

	items, err := storage.Client.Query(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Logger().Debugj(log.JSON{"query": query})
	attributevalue.UnmarshalListOfMaps(items.Items, &trips)
	lastPage := c.QueryParam("page")
	var nextPage string
	err = attributevalue.Unmarshal(items.LastEvaluatedKey["sk"], &nextPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"trips":         trips,
		"count":         items.Count,
		"scanned_count": items.ScannedCount,
		"last_page":     lastPage,
		"next_page":     nextPage,
	})
}

func UpdateTrips(c echo.Context) error {
	return nil
}

func DeleteTrips(c echo.Context) error {
	return nil
}
