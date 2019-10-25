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

type go0 struct {
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
const (
	api     = "/api"
	receive = "/receive.cgi"
	search  = "/search.cgi"
)

var (
	tokens sync.Map
)

type StanderRequest struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {
	// 加载配置文件
	conf := go0{}
	if err := goapp.LoadXMLConfig(config, &conf); err != nil {
		log.Fatal(err)
	}
	app := goapp.Default()

	u, err := url.Parse(conf.Matrix.URI)
	if err != nil {
		log.Fatal(err)
	}
	// 初始化 数据库 连接池
	engine, err := goapp.NewEngine(driver, conf.MySQL.URI, conf.MySQL.ShowSql, conf.MySQL.MaxIdleConns, conf.MySQL.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}
	if err = engine.Ping(); err != nil {
		log.Println(err)
	}

	// 初始化 认证器
	tokens.Store(conf.Secret, goapp.Token{
		Secret: conf.Secret,
	})
	// 静态文件处理流程
	app.HandleDir("/", "html")

	// TCP 连接到资源发布中心
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

	// API 自定义单元
	api := handler{
		dbEngine: engine,
		party:    app.Party(api),
	}

	api.Handler()

	// 验证过滤
	app.Use(func(ctx goapp.Context) {
		authorization := ctx.GetHeader("authorization")
		if authorization == "" || len(authorization) < 32 {
			goto ErrorAuth
		} else {
			// 授权应用、资源传输单元认证
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
	app.Post(receive, api.ApiUpload)
	app.Post(search, api.ApiSearch)

	if err := app.Run(goapp.Addr(fmt.Sprintf(":%d", conf.Port)), nil); err != nil {
		log.Fatal(err)
	}
}

/**
 * TODO: API 上报、查询数据
 */
type handler struct {
	dbEngine *goapp.Engine
	party    goapp.Party
}

/**
 * TODO: 自定义
 */
func (p *handler) Handler() {
	p.party.Post("/", func(ctx goapp.Context) {

	})
}

/**
 * TODO: 数据上报
 */
func (p *handler) ApiUpload(ctx goapp.Context) {
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
}

/**
 * TODO：数据查询
 */
func (p *handler) ApiSearch(ctx goapp.Context) {

}
