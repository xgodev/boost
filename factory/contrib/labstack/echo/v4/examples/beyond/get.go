package main

import (
	"github.com/xgodev/boost/wrapper/config"
	"net/http"

	e "github.com/labstack/echo/v4"
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
