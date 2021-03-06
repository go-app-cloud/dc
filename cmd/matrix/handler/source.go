package handler

import (
	"github.com/go-app-cloud/dc/cmd/matrix/db"
	"github.com/go-app-cloud/goapp"
	"github.com/micro/go-micro/util/log"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"sync"
)

//Source source struct
type Source struct {
	Device *sync.Map
}

// 登入/路由/数据源
type loginSourceResponse struct {
	URI   string `json:"uri"`
	Token string `json:"token"`
}

// 刷新/凭证
type refreshTokenResponse struct {
	Token string `json:"token"`
}

//Handler source handler
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
	  @apiParam {String} search data service search uri
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
		service := ctx.FormValue("service")
		section := ctx.FormValue("section")
		_type := ctx.FormValue("type")
		search := ctx.FormValue("search")
		description := ctx.FormValue("description")
		apiDoc := ctx.FormValue("api_doc")
		s := db.Source{
			Id:          strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", ""),
			Name:        name,
			Service:     service,
			Description: description,
			Section:     section,
			ApiDoc:      apiDoc,
			Type:        _type,
			Search:      search,
			Secret:      goapp.BuildPassword(64, goapp.Advance),
		}
		if _, err := dbEngine.Insert(s); err != nil {
			res.Code = goapp.DBError
			res.Msg = err.Error()
			goto ErrorAddCGI
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)
		return
	ErrorAddCGI:
		_, _ = ctx.JSON(res)
	})
	/**
	  @api {post} /source/modify.cgi modify source
	  @apiName modifySource
	  @apiGroup source
	  @apiVersion 0.0.1

	  @apiParam {String} name source name
	  @apiParam {String} uri data service uri
	  @apiParam {String} search data service search uri
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
	party.Post("/modify.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.FormValue("id")
		name := ctx.FormValue("name")
		service := ctx.FormValue("service")
		section := ctx.FormValue("section")
		_type := ctx.FormValue("type")
		search := ctx.FormValue("search")
		description := ctx.FormValue("description")
		apiDoc := ctx.FormValue("api_doc")
		s := db.Source{
			Name:        name,
			Service:     service,
			Description: description,
			Section:     section,
			ApiDoc:      apiDoc,
			Type:        _type,
			Search:      search,
		}
		if _, err := dbEngine.Where("id = ?", id).Cols("name", "service", "description", "section", "search", "api_doc", "type").Update(&s); err != nil {
			res.Code = goapp.DBError
			res.Msg = err.Error()
			goto ErrorModifyCGI
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)
		return
	ErrorModifyCGI:
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
			"code":0,
		    "data":{
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
				 "search":"..."
				 "type": "...",
				 "create_at":0
				}]
			}
		}
	*/
	party.Get("/select.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}

		page := ctx.URLParam("page")
		pl, err := strconv.Atoi(page)
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
			if err := dbEngine.Desc("create_at").Limit(pl, i*pl).Find(&sources); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				goto ErrorSelectCGI
			}
			res.Code = goapp.Success

			sed := make([]db.SourceEx, 0)
			for _, v := range sources {
				c := 0
				_, ok := p.Device.Load(v.Id)
				if ok {
					c = 1
				}
				sed = append(sed, db.SourceEx{
					Id:          v.Id,
					Name:        v.Name,
					Description: v.Description,
					Secret:      v.Secret,
					Owner:       v.Owner,
					Service:     v.Service,
					Search:      v.Search,
					ApiDoc:      v.ApiDoc,
					Type:        v.Type,
					Section:     v.Section,
					CreateAt:    v.CreateAt,
					Check:       c,
				})

			}
			res.Data = goapp.Map{
				"total": total,
				"items": sed,
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
		if _, err := dbEngine.Where("id = ?", id).Delete(&db.Source{}); err != nil {
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
	  @api {get} /source/get.cgi get source by id
	  @apiName getSource
	  @apiGroup source

	  @apiParam {String} id source id

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0,
			"data":{
	  			"id": "e3be251c7bf945078b8ec6eb40d3a8d9",
				"name": "新生入学报名.数据服务",
				"description": "新生入学报名.数据服务、监护人详细信息、班级详细信息、学校录取通知书等。",
				"section": "大数据中心",
				"owner": "",
				"service": "...",
				"search":"...."
				"api_doc": "...",
				"type": "数据库同步",
				"check": 0,
				"secret": "Yntt3lG1Om)%t.YGUm0NYAvgUPE(9XS$.bt]TTKhC3Vg2+g%6oIc@6HS*Nvg3,6.",
				"create_at": 0
			}
	      }
	*/
	party.Post("/get.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.FormValue("id")
		s := db.Source{}
		b, err := dbEngine.Where("id = ?", id).Get(&s)
		if err != nil || !b {
			res.Code = goapp.DBError
			log.Error(err)
			goto ErrorDeleteCGI
		}
		res.Code = goapp.Success
		res.Data = s
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
			res.Data = loginSourceResponse{
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
			res.Data = refreshTokenResponse{
				Token: t,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorLogin:
		_, _ = ctx.JSON(res)
	})
}
