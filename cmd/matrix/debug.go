// +build !release

package main

import "github.com/go-app-cloud/goapp"

func static(app *goapp.Application) {
	app.HandleDir("/", "html")
}
