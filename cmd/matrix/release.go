// +build release

package main

import "github.com/go-app-cloud/goapp"

func static(app *goapp.Application) {
	app.HandleDir("/", "./html",
		goapp.DirOptions{
			Asset:      Asset,
			AssetInfo:  AssetInfo,
			AssetNames: AssetNames,
			ShowList:   true,
		})
}
