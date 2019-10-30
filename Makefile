.PHONY: clean build

clean: 
	rm -rf ./bin/signup/signup
	rm -rf ./bin/login/login
	rm -rf ./bin/authorizer/authorizer
	rm -rf ./bin/userinfo/userinfo
	rm -rf ./bin/sendemail/sendemail
	rm -rf ./bin/verifyemail/verifyemail
	rm -rf ./bin/createcaptcha/createcaptcha

build:
	GOOS=linux GOARCH=amd64 go build -o bin/signup/signup ./functions/signup
	GOOS=linux GOARCH=amd64 go build -o bin/login/login ./functions/login
	GOOS=linux GOARCH=amd64 go build -o bin/authorizer/authorizer ./functions/authorizer
	GOOS=linux GOARCH=amd64 go build -o bin/userinfo/userinfo ./functions/userinfo
	GOOS=linux GOARCH=amd64 go build -o bin/sendemail/sendemail ./functions/sendemail
	GOOS=linux GOARCH=amd64 go build -o bin/verifyemail/verifyemail ./functions/verifyemail
	GOOS=linux GOARCH=amd64 go build -o bin/createcaptcha/createcaptcha ./functions/createcaptcha
