package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	allowOrigins = []string{
		"http://gostella.tk",
		"https://gostella.tk",
		"http://sub.gostella.tk",
		"https://sub.gostella.tk",
	}
)

type RawCookie struct {
	Name   string `json:"name" validate:"required"`
	Value  string `json:"value" validate:"required"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
	MaxAge int    `json:"max_age" validate:"required"`
	Secure bool   `json:"secure,omitempty"`
}

type BakeRequest struct {
	RawCookies []RawCookie `json:"raw_cookies" validate:"required"`
}

func bakeHandler(c echo.Context) error {
	var bakeRequest BakeRequest
	if err := c.Bind(&bakeRequest); err != nil {
		return err
	}

	var bakedCookies []*http.Cookie

	now := time.Now()
	for _, rawCookie := range bakeRequest.RawCookies {
		expires := now.Add(time.Duration(rawCookie.MaxAge) * time.Second)
		cookie := &http.Cookie{
			Name:    rawCookie.Name,
			Value:   rawCookie.Value,
			Expires: expires,
		}
		if len(rawCookie.Path) > 0 {
			cookie.Path = rawCookie.Path
		}
		if len(rawCookie.Domain) > 0 {
			cookie.Domain = rawCookie.Domain
		}
		if rawCookie.Secure {
			cookie.Secure = true
		}

		c.SetCookie(cookie)
		bakedCookies = append(bakedCookies, cookie)
	}

	return c.JSON(http.StatusOK, bakedCookies)
}

func addResponseHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Robots-Tag", "noindex")
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		// HSTSMaxAge:            31536000,
		// HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "",
	}))
	e.Use(addResponseHeader)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodPost},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		MaxAge:           60 * 60,
	}))
	e.Use(middleware.Gzip())

	e.POST("/", bakeHandler)
	e.POST("/bake/", bakeHandler)

	address := ":5000" // ElasticBeanstalk default port
	e.Logger.Fatal(e.Start(address))
}
