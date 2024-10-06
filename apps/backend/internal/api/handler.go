package api

import (
	"net/http"

	"github.com/Taiterbase/vtrips/apps/backend/internal/storage"
	"github.com/Taiterbase/vtrips/apps/backend/internal/storage/index"
	"github.com/Taiterbase/vtrips/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo"
)

func CreateTrip(c echo.Context) error {
	var (
		trip     = models.NewTrip()
		clientID = c.QueryParam("client_id")
	)

	err := c.Bind(&trip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}
	if err = trip.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	item, err := attributevalue.MarshalMap(trip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	item["PK"] = &attributevalue.AttributeValueMemberS{Value: index.MakePK(clientID)}
	item["SK"] = &attributevalue.AttributeValueMemberS{Value: index.MakeSK(trip.GetID())}
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

	out, err := storage.Client.TransactWriteItems(c.Request().Context(), tx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// now add any write actions for manually maintained indexes
	// then add any transactions on top of this create if needed
	// then execute the write transaction

	return nil
}

func GetTrip(c echo.Context) error {
	return nil
}

func UpdateTrip(c echo.Context) error {
	return nil
}

func DeleteTrip(c echo.Context) error {
	return nil
}

func ListTrip(c echo.Context) error {
	return nil
}

func UpdateTrips(c echo.Context) error {
	return nil
}

func DeleteTrips(c echo.Context) error {
	return nil
}
