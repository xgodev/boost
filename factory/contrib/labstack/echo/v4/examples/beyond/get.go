package main

import (
	"net/http"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/config"
)

func Get(c e.Context) (err error) {

	resp := Response{
		Message: "Hello Google!!",
	}

	err = config.Unmarshal(&resp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
