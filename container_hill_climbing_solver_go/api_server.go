package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joseesco24/sudoku_solver/container_hill_climbing_solver_go/model"

	"github.com/ansel1/merry"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.GET("/health_test", func(context echo.Context) error {

		return context.NoContent(http.StatusOK)

	})

	e.GET("/solver", func(context echo.Context) error {

		var authorization string = context.Request().Header.Get("Authorization")

		if authorization == "" {
			return merry.New("authorization not found").
				WithValue("authorization", authorization).
				WithUserMessage("authorization not found").
				WithHTTPCode(http.StatusUnauthorized)
		}

		if authorization != os.Getenv("ACCESS_KEY") {
			return merry.New("authorization not valid").
				WithValue("authorization", authorization).
				WithUserMessage("authorization not valid").
				WithHTTPCode(http.StatusUnauthorized)
		}

		var middleProxyRequestBody model.MiddleProxyRequest

		body, err := ioutil.ReadAll(context.Request().Body)
		if err != nil {
			return merry.Wrap(err).
				WithValue("middleProxyRequestBody", string(body)).
				WithUserMessage("error while reading request body").
				WithHTTPCode(http.StatusInternalServerError)
		}

		err = json.Unmarshal(body, &middleProxyRequestBody)
		if err != nil {
			return merry.Wrap(err).
				WithValue("middleProxyRequestBody", string(body)).
				WithUserMessage("error while parsing the request body to struct").
				WithHTTPCode(http.StatusInternalServerError)
		}

		if middleProxyRequestBody.Restarts == 0 {
			middleProxyRequestBody.Restarts = 10
		}

		if middleProxyRequestBody.Searchs == 0 {
			middleProxyRequestBody.Searchs = 10
		}

		return context.NoContent(http.StatusOK)

	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("ACCESS_PORT")))

}
