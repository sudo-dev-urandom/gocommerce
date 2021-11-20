package unittest

import (
	"fmt"
	"gocommerce/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	// Setup
	jsonFile, err := os.Open("../request/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	userPayload := string(byteValue)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(userPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	res := rec.Result()
	defer res.Body.Close()

	// fmt.Println("blablabla", rec.Body)

	// Assertions
	if assert.NoError(t, controllers.UserList(c)) {
		// assert.Equal(t, 201, 201)
		assert.Equal(t, strings.TrimSpace(userPayload), strings.TrimSpace(rec.Body.String()))
	}
}
