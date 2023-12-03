package controller

import (
	"errors"
	"github.com/WkBqVN/BackendTest/controller/api"
	_ "github.com/WkBqVN/BackendTest/docs"
	"github.com/WkBqVN/BackendTest/model"
	"github.com/WkBqVN/BackendTest/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"reflect"
	"sync"
)

type Controller struct {
	ControllerRouter *gin.Engine
	StockService     service.StockService
}

const BasePath = "." + afero.FilePathSeparator + "controller" + afero.FilePathSeparator + "route" + afero.FilePathSeparator
const FileName = "stock.json"

var once sync.Once
var controller *Controller

func GetInstance() *Controller {
	once.Do(func() {
		controller = &Controller{}
	})
	return controller
}

func (controller *Controller) InitController() error {
	controller.ControllerRouter = gin.Default()
	err := controller.initRoute()
	if err != nil {
		return err
	}
	return nil
}

func (controller *Controller) initRoute() error {
	controller.ControllerRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// each service can connect which db it uses. Example it can connect to postgres or mysql if it needs
	// db will be init base on json config file
	// after init service config will be stored all need data
	err := controller.StockService.Init(BasePath + FileName)
	if err != nil {
		return err
	}
	if err = controller.StockService.ConnectToDB("./"); err != nil {
		return err
	}
	err = controller.setRoute(*controller.StockService.ServiceConfig)
	if err != nil {
		return err
	}
	return nil
}

func (controller *Controller) setRoute(apiConfig model.RestApiConfig) error {
	if apiConfig.Prefix != "" {
		router := controller.ControllerRouter.Group("/" + apiConfig.Prefix)
		if err := controller.setHandler(router, apiConfig); err != nil {
			return err
		}
		return nil
	}
	err := controller.setHandler(controller.ControllerRouter, apiConfig)
	if err != nil {
		return err
	}
	return nil
}

// handle logic is complex with reflect
func (controller *Controller) setHandler(router gin.IRouter, apiConfig model.RestApiConfig) error {
	for _, method := range apiConfig.Methods {
		executeMethod := reflect.ValueOf(&controller.StockService).MethodByName(method.HandlerFunc)
		if !executeMethod.IsValid() {
			return errors.New("method not found")
		}
		switch method.MethodType {
		case "GET":
			if method.Params != "" {
				router.GET(apiConfig.EndPoint+"/"+method.Params, api.GetStockByIDApi(executeMethod))
			} else {
				router.GET(apiConfig.EndPoint, api.GetStocksApi(executeMethod))
			}
		case "POST":
			router.POST(apiConfig.EndPoint+"/"+method.Params, api.CreateStockApi(executeMethod))
		case "PATCH":
			router.PATCH(apiConfig.EndPoint+"/"+method.Params, api.UpdatePriceStockApi(executeMethod))
		case "DELETE":
			router.DELETE(apiConfig.EndPoint+"/"+method.Params, api.DeleteStockApi(executeMethod))
		}
	}
	return nil
}
