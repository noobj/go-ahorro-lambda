package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
	"github.com/noobj/go-serverless-services/internal/types"
)

func GenerateRedirectResponse[T types.ApiResponse](urlForRedirect string) (T, error) {
	var res T

	switch t := any(&res).(type) {
	case *events.APIGatewayProxyResponse:
		t.StatusCode = 301
		t.IsBase64Encoded = false
		t.Headers = map[string]string{
			"Location": urlForRedirect,
		}
	case *events.APIGatewayV2HTTPResponse:
		t.StatusCode = 301
		t.IsBase64Encoded = false
		t.Headers = map[string]string{
			"Location": urlForRedirect,
		}
	}

	return res, nil
}

func GenerateApiResponse[T types.ApiResponse](resultForBody interface{}) (T, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(resultForBody)

	var res T

	// TODO: any better any to do this mess
	if err != nil {
		switch t := any(&res).(type) {
		case *events.APIGatewayProxyResponse:
			t.StatusCode = 404
		case *events.APIGatewayV2HTTPResponse:
			t.StatusCode = 404
		}

		return res, err
	}

	json.HTMLEscape(&buf, body)

	switch t := any(&res).(type) {
	case *events.APIGatewayProxyResponse:
		t.StatusCode = 200
		t.IsBase64Encoded = false
		t.Body = buf.String()
		t.Headers = map[string]string{
			"Content-Type": "application/json",
		}
	case *events.APIGatewayV2HTTPResponse:
		t.StatusCode = 200
		t.IsBase64Encoded = false
		t.Body = buf.String()
		t.Headers = map[string]string{
			"Content-Type": "application/json",
		}
	}

	return res, nil
}

func GenerateErrorResponse[T types.ApiResponse](statusCode int, messages ...string) (T, error) {
	messageResp := StatusCodeDefaultMsgMap[statusCode]
	if len(messages) != 0 {
		messageResp = strings.Join(messages, "")
	}

	var resType T
	var res any
	switch t := any(resType).(type) {
	case events.APIGatewayProxyResponse:
		t.Body = messageResp
		t.StatusCode = statusCode
		res = t
	case events.APIGatewayV2HTTPResponse:
		t.Body = messageResp
		t.StatusCode = statusCode
		res = t
	}

	return res.(T), nil
}

func SetCookie(cookie http.Cookie, reps *events.APIGatewayV2HTTPResponse) {
	reps.Cookies = append(reps.Cookies, cookie.String())
}

func ParseMultipartForm(contentType string, body io.Reader, isBase64encoded bool) (*multipart.Form, error) {

	if isBase64encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if !strings.HasPrefix(mediaType, "multipart/") || err != nil {
		if err != nil {
			log.Println(err)
		}

		return nil, err
	}
	mr := multipart.NewReader(body, params["boundary"])
	formData, err := mr.ReadForm(5000)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	return formData, nil
}

func ParseCookie(cookies []string) map[string]string {
	result := make(map[string]string)
	for _, cookie := range cookies {
		splitStrings := strings.SplitN(cookie, "=", 2)
		if len(splitStrings) != 2 {
			continue
		}

		result[splitStrings[0]] = splitStrings[1]
	}

	return result
}

func ExtractPayloadFromToken(key string, jwtToken string) (interface{}, error) {
	var claims types.MyCustomClaims
	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		log.Printf("jwt parse error: %v", err)
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims.Payload, nil
}

var StatusCodeDefaultMsgMap = map[int]string{
	401: "please login in",
	500: "internal error",
}

func SendGetRequest(req string) error {
	res, err := http.Get(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)

		return err
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		fmt.Printf("error making http request: %s\n", body)
		fmt.Println(err)
		res.Body.Close()

		return fmt.Errorf("request send to telegram failed")
	}

	return nil
}
