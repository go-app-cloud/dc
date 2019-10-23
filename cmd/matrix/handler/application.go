package handler

import (
	"github.com/go-app-cloud/dc/cmd/matrix/db"
	"github.com/go-app-cloud/goapp"
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
	  @apiParam {String} uri service uri

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
	 @api {get} /application/select.cgi select application item
	 @apiName SelectApplication
	 @apiGroup application

	 @apiSuccessExample {json} Success-Response:
	     HTTP/1.1 200 OK
	     {
	       "total": 60,
	       "items":[{
	          "api_doc": "..."
				"check": 1
				"description": "..."
				"id": "..."
				"name": "..."
				"owner": ""
				"section": "..."
				"service": "..."
				"type": "..."
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
}
