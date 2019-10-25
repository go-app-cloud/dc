package handler

import (
	"github.com/go-app-cloud/dc/cmd/matrix/db"
	"github.com/go-app-cloud/goapp"
	"github.com/go-xorm/xorm"
	"github.com/micro/go-micro/util/log"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
)

type Application struct {
}

func (p *Application) Handler(party goapp.Party, dbEngine *goapp.Engine) {
	/**
	  @api {post} /application/add.cgi add new application
	  @apiName addApplication
	  @apiGroup application

	  @apiParam {String} name application name
	  @apiParam {String} type application type
	  @apiParam {String} source application base from source
	  @apiParam {String} section section
	  @apiParam {String} description description

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0
	      }
	*/
	party.Post("/add.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}

		name := ctx.FormValue("name")
		_type := ctx.FormValue("type")
		sources := strings.Split(ctx.FormValue("source"), ";")
		section := ctx.FormValue("section")
		description := ctx.FormValue("description")

		_, err := dbEngine.Transaction(func(session *xorm.Session) (interface{}, error) {
			appId := strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
			ap := db.Application{
				Id:          appId,
				Name:        name,
				Description: description,
				Section:     section,
				SecretKey:   goapp.BuildPassword(128, goapp.Advance),
				Type:        _type,
			}
			models := make([]db.AppSource, 0)
			for _, v := range sources {
				model := db.AppSource{
					Id:       strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", ""),
					AppId:    appId,
					SourceId: v,
				}
				models = append(models, model)
			}
			if _, err := session.Insert(&models, &ap); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			res.Msg = err.Error()
			res.Code = goapp.DBError
			_, _ = ctx.JSON(res)
			return
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)

	})
	/**
	  @api {post} /application/modify.cgi modify application information
	  @apiName modifyApplication
	  @apiGroup application

	  @apiParam {String} name application name
	  @apiParam {String} name application name
	  @apiParam {String} type application type
	  @apiParam {String} source application base from source
	  @apiParam {String} section section
	  @apiParam {String} description description

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
		_type := ctx.FormValue("type")
		sources := strings.Split(ctx.FormValue("source"), ";")
		section := ctx.FormValue("section")
		description := ctx.FormValue("description")

		_, err := dbEngine.Transaction(func(session *xorm.Session) (interface{}, error) {
			entity := db.Application{
				Name:        name,
				Type:        _type,
				Section:     section,
				Description: description,
			}
			ok, err := session.ID(id).Get(&entity)
			if err != nil {
				return nil, err
			}
			if ok {
				// update application information
				if _, err := session.Cols("name", "description", "section", "type").Update(&entity); err != nil {
					return nil, err
				}
				// delete all app source item
				if _, err := session.Where("app_id = ?", id).Delete(&db.AppSource{}); err != nil {
					return nil, err
				}
				// insert app source item
				models := make([]db.AppSource, 0)
				for _, v := range sources {
					if len(v) < 32 {
						continue
					}
					model := db.AppSource{
						Id:       strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", ""),
						AppId:    id,
						SourceId: v,
					}
					models = append(models, model)
				}
				if _, err := session.Insert(models); err != nil {
					return nil, err
				}
			}
			return nil, nil
		})
		if err != nil {
			res.Msg = err.Error()
			res.Code = goapp.DBError
			_, _ = ctx.JSON(res)
			return
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)

	})
	/**
	 @api {get} /application/select.cgi select application item
	 @apiName SelectApplication
	 @apiGroup application

	 @apiParam {String} page page size
	 @apiParam {String} index page index
	 @apiSuccessExample {json} Success-Response:
	     HTTP/1.1 200 OK
	     {
	       "code":0,
		   "data":{
				"total": 60,
				"items":[{
					description: "电业局指..."
					id: "301689d9....d92d8ae86e8"
					name: "电业在线业务平台"
					owner: ""
					secret_key: "6c...7,h.)=ADFSoYl[0]hiJ6Kr$O+RjzgmL~eOb[Kfg[AK"
					section: "发展和改革局"
					type: "平台应用"
				}]
			}
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
			total, err := dbEngine.Count(new(db.Application))
			if err != nil {
				goto ErrorSelectCGI
			}
			applications := make([]db.Application, 0)
			if err := dbEngine.Limit(p, i*p).Find(&applications); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				goto ErrorSelectCGI
			}
			res.Code = goapp.Success
			res.Data = goapp.Map{
				"total": total,
				"items": applications,
			}
			_, _ = ctx.JSON(res)
			return
		}
	ErrorSelectCGI:
		_, _ = ctx.JSON(res)
	})
	/**
	 @api {get} /application/source.cgi select application by source item
	 @apiName SelectSourceApplication
	 @apiGroup application

	 @apiParam {String} id application id
	 @apiSuccessExample {json} Success-Response:
	     HTTP/1.1 200 OK
	     {
	       "code": 0,
	       "data":[{
	          "id": "...",
			  "name": "..."
	       }]
	     }
	*/
	party.Get("/source.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.URLParam("id")
		result, err := dbEngine.QueryString("select source.id, source.name from app_source left join source on app_source.source_id = source.id where app_source.app_id = ?", id)
		if err != nil {
			goto ErrorSelectSourceCGI
		}
		res.Code = 0
		res.Data = result
		ctx.JSON(res)
		return
	ErrorSelectSourceCGI:
		_, _ = ctx.JSON(res)
	})

	/**
	  @api {get} /application/source_delete.cgi delete application
	  @apiName deleteApplicationSource
	  @apiGroup application

	  @apiParam {String} id source id

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0
	      }
	*/
	party.Get("/source_delete.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.URLParam("id")
		app := ctx.URLParam("app")
		if _, err := dbEngine.Where("source_id = ? and app_id = ?", id, app).Delete(&db.AppSource{}); err != nil {
			res.Code = goapp.DBError
			res.Msg = err.Error()
			_, _ = ctx.JSON(res)
			return
		}
		res.Code = goapp.Success
		_, _ = ctx.JSON(res)
	})

	/**
	  @api {get} /application/delete.cgi delete application
	  @apiName deleteApplication
	  @apiGroup application

	  @apiParam {String} id application id

	  @apiSuccessExample {json} Success-Response:
	      HTTP/1.1 200 OK
	      {
	        "code": 0
	      }
	*/
	party.Get("/delete.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.URLParam("id")
		if _, err := dbEngine.Transaction(func(session *xorm.Session) (interface{}, error) {
			if _, err := session.ID(id).Delete(&db.Application{}); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				return nil, err
			}
			if _, err := session.In("app_id", id).Delete(&db.AppSource{}); err != nil {
				res.Code = goapp.DBError
				res.Msg = err.Error()
				return nil, err
			}
			return nil, nil
		}); err != nil {
			res.Code = goapp.DBError
			res.Msg = err.Error()
			_, _ = ctx.JSON(res)
			return
		}
		res.Code = goapp.Success
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

			}
	      }
	*/
	party.Post("/get.cgi", func(ctx goapp.Context) {
		res := goapp.Response{}
		id := ctx.FormValue("id")
		s := db.Application{}
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
}
