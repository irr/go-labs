package main

/*

https://epsagon.com/blog/getting-started-with-aws-lambda-and-go/

npm -g install serverless

go mod init lambda-time
go get github.com/aws/aws-lambda-go

GOOS=linux GOARCH=amd64  go build -o bin/lambda-time .
sls deploy

*/

import (
    "context"
    "encoding/json"
    "time"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type response struct {
    UTC time.Time `json:"utc"`
}

func timezone(ip string) *time.Location {
    resp, err := httpClient.Get("https://ipapi.co/" + ip + "/timezone/")
    if err != nil {
        return nil
    }
    defer resp.Body.Close()
    tz, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil
    }
    loc, err := time.LoadLocation(string(tz))
    if err != nil {
        return nil
    }
    return loc
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    now := time.Now()
    resp := &response{
        UTC: now.UTC(),
    }
    body, err := json.Marshal(resp)
    if err != nil {
        return events.APIGatewayProxyResponse{}, err
    }
    return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
    lambda.Start(handleRequest)
}
