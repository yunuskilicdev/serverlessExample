package main

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dchest/captcha"
	"github.com/yunuskilicdev/serverlessNear/common"
	"github.com/yunuskilicdev/serverlessNear/common/errormessage"
	"github.com/yunuskilicdev/serverlessNear/common/model"
	"net/http"
	"time"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	store := common.NewStore(time.Duration(5 * time.Minute))
	captcha.SetCustomStore(store)

	captchaResponse := model.CaptchaResponse{}
	captchaId := captcha.New()

	var ImageBuffer bytes.Buffer
	captcha.WriteImage(&ImageBuffer, captchaId, 300, 90)

	captchaResponse.Id = captchaId
	captchaResponse.Image = base64.StdEncoding.EncodeToString(ImageBuffer.Bytes())

	body := model.ResponseBody{}
	body.Message = errormessage.StatusText(errormessage.Ok)
	body.ResponseObject = captchaResponse
	return events.APIGatewayProxyResponse{
		Body:       body.ConvertToJson(),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
