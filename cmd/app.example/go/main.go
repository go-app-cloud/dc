package main

import (
	"encoding/xml"
	"fmt"
	"github.com/go-app-cloud/goapp"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Go struct {
	XMLName xml.Name `xml:"go"`
	Port    int      `xml:"port,attr"`
	Secret  string   `xml:"secret,attr"`
}

const (
	config = "go.conf.xml"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	devices sync.Map
)

func echo(ctx goapp.Context) {
	c, err := upgrade.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	devices.Store(c, c)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}
func main() {
	conf := Go{}
	if err := goapp.LoadXMLConfig(config, &conf); err != nil {
		log.Fatal(err)
	}
	app := goapp.Default()

	token := goapp.Token{
		Secret: conf.Secret,
	}

	app.HandleDir("/", "html")
	app.Any("/echo", echo)

	app.Use(func(ctx goapp.Context) {
		authorization := ctx.GetHeader("authorization")
		if authorization == "" || len(authorization) < 32 {
			ctx.StatusCode(400)
			return
		}
		_, err := token.Parse(authorization)
		if err != nil {
			ctx.StatusCode(400)
			return
		}
		ctx.Next()
	})
	app.Post("/rec.cgi", func(ctx goapp.Context) {
		var a interface{}
		ctx.ReadJSON(&a)
		devices.Range(func(key, value interface{}) bool {
			c := value.(*websocket.Conn)
			c.WriteJSON(a)
			log.Println(a)
			return true
		})
		ctx.JSON(goapp.Map{
			`success`: true,
		})
	})

	if err := app.Run(goapp.Addr(fmt.Sprintf(":%d", conf.Port)), nil); err != nil {
		log.Fatal(err)
	}
}
