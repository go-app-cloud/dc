package handler

import (
	"github.com/go-app-cloud/dc/cmd/matrix/db"
	"github.com/go-app-cloud/goapp"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
)

type Source struct {
}

// 登入/路由/数据源
type LoginResponse struct {
	URI   string `json:"uri"`
	Token string `json:"token"`
}

// 刷新/凭证
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (p *Source) Handler(party goapp.Party, dbEngine *goapp.Engine) {
	/**
	  @api {post} /source/add.cgi add source
	  @apiName addSource
	  @apiGroup source
	  @apiVersion 0.0.1

	  @apiParam {String} name source name
	  @apiParam {String} uri data service uri
	  @apiParam {String} check network check type
	  @apiParam {String} section system of section
	  @apiParam {String} type data sync type
	  @apiParam {String} description description forsource
	  @apiParam {String} api_doc api document url

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0
	      }
	*/
	party.Post("/add.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		name := ctx.FormValue("name")
		uri := ctx.FormValue("uri")
		section := ctx.FormValue("section")
		check, err := strconv.Atoi(ctx.FormValue("check"))
		_type := ctx.FormValue("type")
		description := ctx.FormValue("description")
		apiDoc := ctx.FormValue("api_doc")
		if err != nil {
			res.Code = goapp.CheckTypeError
			res.Msg = err.Error()
			goto ErrorAddCGI
		} else {
			s := db.Source{
				Id:          strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", ""),
				Name:        name,
				Service:     uri,
				Description: description,
				Section:     section,
				ApiDoc:      apiDoc,
				Type:        _type,
				Check:       check,
				Secret:      goapp.BuildPassword(64, goapp.Advance),
			}
			if _, err = dbEngine.Insert(s); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				goto ErrorAddCGI
			}
			res.Code = goapp.Success
			_, _ = ctx.JSON(res)
			return
		}
	ErrorAddCGI:
		_, _ = ctx.JSON(res)
	})
	/**
	 @api {get} /source/select.cgi select source items
	 @apiName selectSource
	 @apiGroup source
	 @apiVersion 0.0.1

	 @apiParam {Number} index index of page
	 @apiParam {Number} page  page size

	 @apiSuccessExample {json} Success-Response:
	   HTTP/1.1 200 OK
	    {
		 "total": 60,
		 "items":[{
			 "api_doc": "...",
			 "check": 1,
			 "description": "...",
			 "id": "...",
			 "name": "...",
			 "owner": "",
			 "section": "...",
			 "service": "...",
			 "type": "...",
			 "create_at":0
			}]
		}
	*/
	party.Get("/select.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}

		page := ctx.URLParam("page")
		p, err := strconv.Atoi(page)
		if err != nil {
			goto ErrorSelectCGI
		} else {
			index := ctx.URLParam("index")
			i, err := strconv.Atoi(index)
			if err != nil {
				goto ErrorSelectCGI
			}
			if i > 0 {
				i--
			}
			total, err := dbEngine.Count(new(Source))
			if err != nil {
				goto ErrorSelectCGI
			}
			sources := make([]db.Source, 0)
			if err := dbEngine.Desc("create_at").Limit(p, i*p).Find(&sources); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				goto ErrorSelectCGI
			}
			res.Code = goapp.Success
			res.Data = goapp.Map{
				"total": total,
				"items": sources,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorSelectCGI:
		_, _ = ctx.JSON(res)
	})
	/**
	  @api {get} /source/delete.cgi delete source
	  @apiName deleteSource
	  @apiGroup source

	  @apiParam {String} id source id

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0
	      }
	*/
	party.Get("/delete.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.URLParam("id")
		if _, err := dbEngine.Id(id).Delete(&db.Source{}); err != nil {
			res.Code = goapp.DBError
			res.Msg = err.Error()
			goto ErrorDeleteCGI
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)
		return
	ErrorDeleteCGI:
		_, _ = ctx.JSON(res)
	})

	/**
	  @api {post} /source/login.cgi route source
	  @apiName routeSource
	  @apiGroup source

	  @apiParam {String} id app id
	  @apiParam {String} secret secret key

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0,
	        "data": {
	         "uri": "*****",
	         "token": "*****",
	        },
	      }
	*/
	party.Post("/login.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		appId := ctx.FormValue("id")
		secretKey := ctx.FormValue("secret")
		s := db.Source{}
		ok, err := dbEngine.Where("id = ? and secret = ?", appId, secretKey).Get(&s)
		if err != nil {
			res.Code = -1
			res.Msg = err.Error()
			goto ErrorLogin
		} else {
			if !ok {
				res.Code = -2
				goto ErrorLogin
			}
			token := goapp.Token{Secret: s.Secret}
			t, err := token.Build(goapp.Claims{Data: goapp.App{AppId: s.Id, SecretKey: s.Secret}})
			if err != nil {
				res.Code = -3
				res.Msg = err.Error()
				goto ErrorLogin
			}
			res.Data = LoginResponse{
				URI:   s.Service,
				Token: t,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorLogin:
		_, _ = ctx.JSON(res)
	})
	/**
	  @api {post} /source/refresh.cgi refresh source token
	  @apiName refreshSourceToken
	  @apiGroup source

	  @apiParam {String} id app id
	  @apiParam {String} secret secret key

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0,
	        "data": {
	         "token": "*****",
	        },
	      }
	*/
	party.Post("/refresh.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		appId := ctx.FormValue("id")
		secretKey := ctx.FormValue("secret")
		s := db.Source{}
		ok, err := dbEngine.Where("id = ? and secret = ?", appId, secretKey).Get(&s)
		if err != nil {
			res.Code = -1
			res.Msg = err.Error()
			goto ErrorLogin
		} else {
			if !ok {
				res.Code = -2
				goto ErrorLogin
			}
			token := goapp.Token{Secret: s.Secret}
			t, err := token.Build(goapp.Claims{Data: goapp.App{AppId: s.Id, SecretKey: s.Secret}})
			if err != nil {
				res.Code = -3
				res.Msg = err.Error()
				goto ErrorLogin
			}
			res.Data = RefreshTokenResponse{
				Token: t,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorLogin:
		_, _ = ctx.JSON(res)
	})
}
