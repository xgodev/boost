package health

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/model/errors"
)

// Checker represents resty client checker.
type Checker struct {
	client  *resty.Client
	options *Options
}

// Check checks if resty client is OK.
func (c *Checker) Check(ctx context.Context) (err error) {

	request := c.client.R()

	var response *resty.Response

	response, err = request.Get(strings.Join([]string{c.options.Host, c.options.Endpoint}, ""))
	if err != nil {
		return errors.Internalf(err.Error())
	}

	if response.IsError() {
		return errors.New(strconv.Itoa(response.StatusCode()))
	}

	return err
}

// NewChecker creates the resty client checker.
func NewChecker(client *resty.Client, options *Options) *Checker {
	return &Checker{client: client, options: options}
}
