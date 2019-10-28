package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/go-app-cloud/dc/cmd/matrix/db"
	"github.com/go-app-cloud/dc/cmd/matrix/handler"
	"github.com/go-app-cloud/goapp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"log"
)

const (
	config             = "matrix.conf.xml"
	driver             = "mysql"
	_uriLogin          = "/login.cgi"
	_routerSource      = "/source"
	_routerApplication = "/application"
)

type matrix struct {
	XMLName xml.Name `xml:"matrix"`
	Port    int      `xml:"port,attr"`
	Secret  string   `xml:"secret,attr"`
	MySQL   struct {
		URI          string `xml:"uri,attr"`
		ShowSql      bool   `xml:"show_sql,attr"`
		MaxIdleConns int    `xml:"max_idle_conns,attr"`
		MaxOpenConns int    `xml:"max_open_conns,attr"`
	} `xml:"mysql"`
}

type Auth struct {
	Id       string
	Username string
}

func main() {
	app := goapp.Default()
	app.Use(func(ctx goapp.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Next()
	})
	conf := matrix{}
	if err := goapp.LoadXMLConfig(config, &conf); err != nil {
		log.Fatal(err)
	}

	engine, err := goapp.NewEngine(driver, conf.MySQL.URI, conf.MySQL.ShowSql, conf.MySQL.MaxIdleConns, conf.MySQL.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}
	if err = engine.Ping(); err != nil {
		log.Println(err)
	}
	static(app)
	/**
	@api {websocket} /source.cgi source connection
	@apiName SourceConnect
	@apiGroup message
	*/
	device := goapp.BuildSocket(func(req goapp.AuthRequest) error {
		source := db.Source{}
		ok, err := engine.Where("id = ?", req.AppId).Get(&source)
		if err != nil || !ok {
			log.Println(err)
			return errors.New("auth error")
		}
		if source.Secret != req.SecretKey {
			return errors.New("auth error")
		}
		return nil
	}, func(req goapp.AuthRequest, conn *websocket.Conn) {

	}, func(conn *websocket.Conn, message goapp.Message) error {
		switch message.Type {
		case 0:
			// heart
		case 1000:
			break
		}
		return nil
	}, func(appId string) {
	}, func(appId string, conn *websocket.Conn) error {
		// update app source secret
		result, err := engine.QueryString("select application.secret_key from app_source left join application on app_id = application.id left join source on source_id = source.id where source.id=?", appId)
		if err != nil {
			log.Println(err)
			return err
		}
		return conn.WriteJSON(goapp.Response{Code: 1001, Data: result})
	})

	app.Any("/source.cgi", device.SocketHandler)

	/**
	  @api {post} /source/login.cgi user login
	  @apiName UserLogin
	  @apiGroup user

	  @apiParam {String} username user name
	  @apiParam {String} pass password

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0,
	        "data": {
	         "token": "*****",
	        },
	      }
	*/
	app.Post(_uriLogin, func(ctx goapp.Context) {
		res := goapp.Response{}
		username := ctx.FormValue("username")
		pass := ctx.FormValue("pass")
		u := db.User{}
		ok, err := engine.Where("username = ?", username).Get(&u)
		if err != nil {
			res.Code = -1
			res.Msg = err.Error()
			goto ErrorLogin
		} else {
			if !ok {
				res.Code = -2
				goto ErrorLogin
			}

			password := goapp.MD5(fmt.Sprintf("%s%s%s", u.Salt, pass, u.Salt))
			log.Println(password)
			if password != u.Password {
				res.Code = -100
				goto ErrorLogin
			}
			token := goapp.Token{Secret: conf.Secret}
			t, err := token.Build(goapp.Claims{Data: Auth{Id: u.Id, Username: u.Username}})
			if err != nil {
				res.Code = -3
				res.Msg = err.Error()
				goto ErrorLogin
			}
			res.Code = 0
			res.Data = goapp.Map{
				"token": t,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorLogin:
		_, _ = ctx.JSON(res)
	})

	// 数据来源
	source := handler.Source{
		Device: &device.Devices,
	}
	source.Handler(app.Party(_routerSource), engine)

	// 数据应用
	application := handler.Application{
		Devices:  &device.Devices,
		DbEngine: engine,
	}
	application.Handler(app.Party(_routerApplication))

	if err := app.Run(goapp.Addr(fmt.Sprintf(":%d", conf.Port)), nil); err != nil {
		log.Fatal(err)
	}
}
