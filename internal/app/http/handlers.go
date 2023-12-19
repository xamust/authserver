package httpapp

import (
	"github.com/labstack/echo/v4"
	"github.com/xamust/authserver/www"
	"html/template"
	"io/fs"
	"net/http"
)

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func index(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.HTML), "html/index.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func login(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.HTML), "html/login.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func style(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.Style), "styles/*")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func scripts(c echo.Context) error {
	tmpl, err := template.ParseFS(fs.FS(www.Scripts), "scripts/*")
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}
