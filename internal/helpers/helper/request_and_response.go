package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func GenerateApiResponse(resultForBody interface{}) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(resultForBody)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"set-cookie":   "xxx=123",
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func SetCookie(cookie http.Cookie, reps *events.APIGatewayProxyResponse) {
	reps.MultiValueHeaders["set-cookie"] = append(reps.MultiValueHeaders["set-cookie"], cookie.String())
}

func ParseMultipartForm(contentType string, body io.Reader) (*multipart.Form, error) {
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
