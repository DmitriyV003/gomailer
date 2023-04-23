package main

import "github.com/labstack/echo/v4"

func (c *Config) routes() *echo.Echo {
	e := echo.New()

	gr := e.Group("api/v1")
	_ = gr

	return e
}
