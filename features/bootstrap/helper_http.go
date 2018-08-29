package bootstrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gregbiv/consus/pkg/api"
	"github.com/onsi/gomega"
)

func assertStatusEquals(response *http.Response, statusCode int) error {
	if response.StatusCode == statusCode {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Invalid response. Error when trying to retrieve the body. StatusCode: %d ", response.StatusCode)
	}

	return fmt.Errorf("Invalid response. StatusCode: %d Body: %s", response.StatusCode, body)
}

func assertNotFoundResponse(response *http.Response) error {
	if err := assertStatusEquals(response, http.StatusNotFound); err != nil {
		return err
	}
	return assertErrorResponse(response.Body, "InvalidUri", "", "The requested URI does not represent any resource on the server.")
}

func assertErrorResponse(body io.ReadCloser, code string, target string, message string) error {
	actualBuff := new(bytes.Buffer)
	actualBuff.ReadFrom(body)
	actual := actualBuff.String()

	expectedError := api.ErrResponse{
		Errors: api.Error{
			Code:    code,
			Target:  target,
			Message: message,
		},
		HTTPStatusCode: http.StatusBadRequest,
	}

	expectedJSON, err := json.Marshal(expectedError)

	if !gomega.Expect(actual).Should(gomega.MatchJSON(expectedJSON)) {
		return errors.New("Invalid response")
	}

	return err
}
