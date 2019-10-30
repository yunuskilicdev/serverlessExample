package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dchest/captcha"
	"github.com/jinzhu/gorm"
	"github.com/yunuskilicdev/serverlessExample/common"
	"github.com/yunuskilicdev/serverlessExample/common/errormessage"
	"github.com/yunuskilicdev/serverlessExample/common/model"
	"github.com/yunuskilicdev/serverlessExample/database"
	"github.com/yunuskilicdev/serverlessExample/database/entity"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var user entity.User
	var loginRequest LoginRequest
	jsonErr := json.Unmarshal([]byte(request.Body), &loginRequest)
	if jsonErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.JsonParseError)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, nil)
	}

	v := validator.New()
	validateErr := v.Struct(loginRequest)

	if validateErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.JsonParseError)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, nil)
	}

	postgresConnector := database.PostgresConnector{}
	dbConn, dbErr := postgresConnector.GetConnection()
	defer dbConn.Close()
	if dbErr != nil {
		fmt.Print(dbErr)
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.DatabaseError)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, dbConn)
	}

	var users []entity.User
	filter := &entity.User{}
	filter.Email = loginRequest.Email
	dbConn.Where(filter).Find(&users)

	if users == nil || len(users) == 0 {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.UserNameOrPasswordWrong)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, dbConn)
	}

	user = users[0]
	if user.LoginTry >= 5 && !validateCaptcha(loginRequest) {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.CaptchaNeeded)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, dbConn)
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if passwordErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.UserNameOrPasswordWrong)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, dbConn)
	}

	tokenSet, signErr := common.CreateTokens(user)

	if signErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.DatabaseError)
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       body.ConvertToJson(),
		}
		return createApiLoginFailResponse(response, user, dbConn)
	}

	user.LoginTry = 0
	dbConn.Save(&user)
	body := model.ResponseBody{}
	body.Message = errormessage.StatusText(errormessage.Ok)
	var loginResponse LoginResponse
	loginResponse.AccessToken = tokenSet.AccessToken
	loginResponse.RefreshToken = tokenSet.RefreshToken
	body.ResponseObject = loginResponse
	resp := events.APIGatewayProxyResponse{
		Body:       body.ConvertToJson(),
		StatusCode: http.StatusOK,
	}
	return resp, nil

}

func validateCaptcha(request LoginRequest) bool {
	if request.CaptchaId == "" || request.CaptchaResponse == "" {
		return false
	}
	store := common.GetStore()
	captcha.SetCustomStore(store)
	return captcha.VerifyString(request.CaptchaId, request.CaptchaResponse)
}

func createApiLoginFailResponse(response events.APIGatewayProxyResponse, user entity.User, dbConn *gorm.DB) (events.APIGatewayProxyResponse, error) {
	if user.ID > 0 {
		user.LoginTry = user.LoginTry + 1
		dbConn.Save(user)
		if user.LoginTry >= 5 {
			body := model.ResponseBody{}
			body.Message = errormessage.StatusText(errormessage.CaptchaNeeded)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       body.ConvertToJson(),
			}, nil
		} else {
			return response, nil
		}
	} else {
		return response, nil
	}
}

func main() {
	lambda.Start(handler)
}
