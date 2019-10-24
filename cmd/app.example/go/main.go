package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-app-cloud/goapp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
	"time"
)

type Go struct {
	XMLName xml.Name `xml:"go"`
	Port    int      `xml:"port,attr"`
	Appid   string   `xml:"appid,attr"`
	Secret  string   `xml:"secret,attr"`
	Matrix  struct {
		URI string `xml:"uri,attr"`
	} `xml:"matrix"`
	MySQL struct {
		URI          string `xml:"uri,attr"`
		ShowSql      bool   `xml:"show_sql,attr"`
		MaxIdleConns int    `xml:"max_idle_conns,attr"`
		MaxOpenConns int    `xml:"max_open_conns,attr"`
	} `xml:"mysql"`
	Mongodb struct {
		URI string `xml:"uri,attr"`
	} `xml:"mongodb"`
}

const (
	config = "go.conf.xml"
	driver = "mysql"
)

var (
	tokens sync.Map
)

type StanderRequest struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {
	conf := Go{}
	if err := goapp.LoadXMLConfig(config, &conf); err != nil {
		log.Fatal(err)
	}
	app := goapp.Default()

	u, err := url.Parse(conf.Matrix.URI)
	if err != nil {
		log.Fatal(err)
	}
	tokens.Store(conf.Secret, goapp.Token{
		Secret: conf.Secret,
	})

	engine, err := goapp.NewEngine(driver, conf.MySQL.URI, conf.MySQL.ShowSql, conf.MySQL.MaxIdleConns, conf.MySQL.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}
	if err = engine.Ping(); err != nil {
		log.Println(err)
	}

	app.HandleDir("/", "html")

	goapp.BuildSocketClient(u, func(raw json.RawMessage) {
		res := goapp.Response{}
		if err := json.Unmarshal(raw, &res); err != nil {
			log.Println(err)
			return
		}
		switch res.Code {
		case 1001:
			d := res.Data
			for _, v := range d.([]interface{}) {
				vv := v.(map[string]interface{})
				secretKey := vv["secret_key"].(string)
				tokens.Store(secretKey, goapp.Token{
					Secret: secretKey,
				})
			}
			break
		}
	}, func(con *websocket.Conn) {
		con.WriteJSON(goapp.AuthRequest{AppId: conf.Appid, SecretKey: conf.Secret})
	}, func(socket *goapp.SocketClient, uri *url.URL) {
		<-time.After(time.Second * 5)
		socket.ReConnect()
	})

	app.Use(func(ctx goapp.Context) {
		authorization := ctx.GetHeader("authorization")
		if authorization == "" || len(authorization) < 32 {
			goto ErrorAuth
		} else {
			isAuth := false
			tokens.Range(func(key, value interface{}) bool {
				token := value.(goapp.Token)
				if _, err := token.Parse(authorization); err != nil {
					return true
				}
				isAuth = true
				return !isAuth
			})
			if isAuth {
				ctx.Next()
				return
			}
		}
	ErrorAuth:
		ctx.StatusCode(400)
		return
	})

	app.Post("/rec.cgi", func(ctx goapp.Context) {
		var req StanderRequest
		ctx.ReadJSON(&req)
		switch req.Code {
		case 0:
			log.Println(req.Data)
			break
		case 1:
			break

		}
		ctx.JSON(req)
	})

	if err := app.Run(goapp.Addr(fmt.Sprintf(":%d", conf.Port)), nil); err != nil {
		log.Fatal(err)
	}
}
