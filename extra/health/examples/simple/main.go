package main

import (
	"context"
	"log"

	"github.com/xgodev/boost/extra/health"
)

type Checker struct {
}

func (c *Checker) Check(ctx context.Context) error {
	return nil
}

func NewChecker() *Checker {
	return &Checker{}
}

func main() {

	hc := health.NewHealthChecker("simple", "simplechecker", NewChecker(), true, true)
	health.Add(hc)

	all := health.CheckAll(context.Background())

	log.Println(all)
}
