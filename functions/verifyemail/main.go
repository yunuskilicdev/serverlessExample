package main

import (
	"encoding/binary"
	"fmt"
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

	token := request.QueryStringParameters["token"]
	validateToken, err := common.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
	}
	claims := validateToken.Claims.(*model.CustomClaims)
	if validateToken.Valid && claims.Type == "Mail" {
		store := common.GetStore()
		value := store.Get(token, true)
		var userFilter entity.User
		u, _ := binary.Uvarint(value)
		userFilter.ID = uint(u)
		postgresConnector := database.PostgresConnector{}
		dbConn, dbErr := postgresConnector.GetConnection()
		defer dbConn.Close()
		if dbErr != nil {
			fmt.Print(dbErr)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "",
			}, nil
		}
		var user entity.User
		dbConn.Where(userFilter).Find(&user)
		user.EmailVerified = true
		dbConn.Save(&user)
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.Ok)
		body.ResponseObject = user
		return events.APIGatewayProxyResponse{
			Body:       body.ConvertToJson(),
			StatusCode: http.StatusOK,
		}, nil
	}

	body := model.ResponseBody{}
	body.Message = errormessage.StatusText(errormessage.TokenIsNotValid)
	body.ResponseObject = nil
	return events.APIGatewayProxyResponse{
		Body:       body.ConvertToJson(),
		StatusCode: http.StatusBadRequest,
	}, nil

}

func main() {
	lambda.Start(handler)
}
