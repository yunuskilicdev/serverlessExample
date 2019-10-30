package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/yunuskilicdev/serverlessExample/common"
	"github.com/yunuskilicdev/serverlessExample/common/errormessage"
	"github.com/yunuskilicdev/serverlessExample/common/model"
	"github.com/yunuskilicdev/serverlessExample/database"
	"github.com/yunuskilicdev/serverlessExample/database/entity"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signupRequest SignupRequest
	jsonErr := json.Unmarshal([]byte(request.Body), &signupRequest)
	if jsonErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.JsonParseError)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body.ConvertToJson(),
		}, nil
	}
	v := validator.New()
	validateErr := v.Struct(signupRequest)

	if validateErr != nil {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.JsonParseError)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body.ConvertToJson(),
		}, nil
	}
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
	dbConn.AutoMigrate(&entity.User{})

	var users []entity.User
	filter := &entity.User{}
	filter.Email = signupRequest.Email
	dbConn.Where(filter).Find(&users)

	if users != nil && len(users) > 0 {
		body := model.ResponseBody{}
		body.Message = errormessage.StatusText(errormessage.UserAlreadyExist)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body.ConvertToJson(),
		}, nil
	}

	newUser := &entity.User{}
	newUser.Email = signupRequest.Email
	newUser.Password = hashAndSalt(signupRequest.Password)
	dbConn.Create(&newUser)

	expireAt := time.Now().Add(1 * time.Hour)
	token, jsonErr := common.CreateToken(*newUser, "Mail", expireAt)

	var mailRequest model.SendVerificationMailRequest
	mailRequest.UserId = newUser.ID
	mailRequest.Token = token
	mailRequest.Email = newUser.Email
	emailJsonData, _ := json.Marshal(mailRequest)
	s := string(emailJsonData)
	u := string(os.Getenv("email_queue_url"))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		fmt.Println(err)
	}
	sqsClient := sqs.New(sess)
	sqsClient.ServiceName = os.Getenv("email_queue")
	input := sqs.SendMessageInput{
		MessageBody: &s,
		QueueUrl:    &u,
	}
	_, jsonErr = sqsClient.SendMessage(&input)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}

	body := model.ResponseBody{}
	body.Message = errormessage.StatusText(errormessage.Ok)
	body.ResponseObject = newUser
	return events.APIGatewayProxyResponse{
		Body:       body.ConvertToJson(),
		StatusCode: http.StatusOK,
	}, nil
}

func hashAndSalt(pwd string) (hashed string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func main() {
	lambda.Start(handler)
}
