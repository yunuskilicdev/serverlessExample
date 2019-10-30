package main

import (
	"encoding/binary"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yunuskilicdev/serverlessExample/common"
	"github.com/yunuskilicdev/serverlessExample/common/errormessage"
	"github.com/yunuskilicdev/serverlessExample/common/model"
	"github.com/yunuskilicdev/serverlessExample/database"
	"github.com/yunuskilicdev/serverlessExample/database/entity"
	"net/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	token := request.Headers["Authorization"]
	userId := common.GetStore().Get(token, false)

	postgresConnector := database.PostgresConnector{}
	dbConn, dbErr := postgresConnector.GetConnection()
	defer dbConn.Close()
	if dbErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.DatabaseError)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       body.ConvertToJson(),
		}
		return response, nil
	}

	var userFilter entity.User
	u, _ := binary.Uvarint(userId)
	userFilter.ID = uint(u)
	var user entity.User
	dbConn.Where(userFilter).Find(&user)

	body := model.ResponseBody{}
	body.Message = errormessage.StatusText(errormessage.Ok)
	body.ResponseObject = user
	return events.APIGatewayProxyResponse{
		Body:       body.ConvertToJson(),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
