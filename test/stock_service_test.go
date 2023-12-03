package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WkBqVN/BackendTest/controller"
	"github.com/WkBqVN/BackendTest/controller/api"
	"github.com/WkBqVN/BackendTest/model"
	"github.com/WkBqVN/BackendTest/repository"
	"github.com/WkBqVN/BackendTest/service"
	"github.com/WkBqVN/BackendTest/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestData struct {
	Data []TestApi `json:"data"`
}

type TestApi struct {
	Url     string    `json:"url"`
	Handler string    `json:"handler"`
	Method  string    `json:"method"`
	Cases   []CaseApi `json:"cases"`
	Params  string    `json:"params"`
	Api     string    `json:"api"`
}

type CaseApi struct {
	Body     *model.Stock `json:"body"`
	Title    string       `json:"title"`
	Response string       `json:"response"`
	Query    string       `json:"query"`
	Status   int          `json:"status"`
}

var (
	repo   = repository.GetInstance()
	router *gin.Engine
)

func init() {
	repo.Config.Host = "localhost"
	repo.Config.Dbname = "postgres"
	repo.Config.Port = 5432
	repo.Config.Password = "1"
	repo.Config.Schema = "testdata"
	repo.Config.User = "postgres"
	router = gin.Default()
}

func setRouteTest(router *gin.Engine, method string, url string, handler gin.HandlerFunc) {
	switch method {
	case "GET":
		router.GET(url, handler)
	case "POST":
		router.POST(url, handler)
	case "PATCH":
		router.PATCH(url, handler)
	case "DELETE":
		router.DELETE(url, handler)
	}
}

func TestGeStocks(t *testing.T) {
	testData := &TestData{}
	// get data same at real api
	if err := utils.ConvertJsonToStruct(".\\testdata.json", testData); err != nil {
		log.Fatal(utils.ErrCannotFormatData)
	}
	c := controller.GetInstance()
	if err := c.StockService.Init("..\\" + service.BasePath + service.FileName); err != nil {
	}

	db, _ := gorm.Open(postgres.Open(repo.ConnectGenerator("PG")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: repo.Config.Schema + ".",
		},
	})
	c.StockService.DB = db
	// test api on each case
	// reflect not call package public method so need set manual
	for _, testApi := range testData.Data {
		executeMethod := reflect.ValueOf(&c.StockService).MethodByName(testApi.Handler)
		if !executeMethod.IsValid() {
		}
		if testApi.Api == "GetStocksApi" {
			setRouteTest(router, testApi.Method, testApi.Url, api.GetStocksApi(executeMethod))
		}
		if testApi.Api == "GetStockByIDApi" {
			setRouteTest(router, testApi.Method, testApi.Url+"/"+testApi.Params, api.GetStockByIDApi(executeMethod))
		}
		if testApi.Api == "UpdatePriceStockById" {
			setRouteTest(router, testApi.Method, testApi.Url+"/"+testApi.Params, api.UpdatePriceStockApi(executeMethod))
		}
		RunTest(t, testApi)
	}
}

func RunTest(t *testing.T, testApi TestApi) {
	for _, testcase := range testApi.Cases {
		rr := httptest.NewRecorder()
		var req *http.Request
		t.Run(fmt.Sprintf("%s", testcase.Title), func(t *testing.T) {
			url := testApi.Url
			if testApi.Params == "" {
				url += "?" + testcase.Query
			} else {
				url += "/" + testcase.Query
			}
			if testcase.Body != nil {
				body, _ := json.Marshal(testcase.Body)
				req, _ = http.NewRequest(testApi.Method, url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(testApi.Method, url, nil)
			}

			c, e := gin.CreateTestContext(rr)
			req, _ = http.NewRequest(testApi.Method, url, nil)
			c.Request = req
			e.GET()
			router.ServeHTTP(rr, req)
			responseData, _ := io.ReadAll(rr.Body)
			assert.Equal(t, testcase.Response, string(responseData))
			assert.Equal(t, testcase.Status, rr.Code)
		})
	}
}
