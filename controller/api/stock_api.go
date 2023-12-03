package api

import (
	"github.com/WkBqVN/BackendTest/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type API struct{}

// @Summary      get all Stock api
// @Description  get list of stocks
// @Accept       json
// @Produce      json
// @Param        q    query   string  false  "name search by q"  Format(email)
// @Success      200  {array} model.Stock
// @Router       /api [get]
func GetStocksApi(executeMethod reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.Query("limit")
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
		}
		page := c.Query("page")
		pageInt, err := strconv.Atoi(page)
		if err != nil {

		}
		methodArgs := []reflect.Value{reflect.ValueOf(limitInt), reflect.ValueOf(pageInt)}
		result := executeMethod.Call(methodArgs)
		if len(result) > 0 { // > 0 mean found data
			responseData := result[0].Interface()
			c.JSON(http.StatusOK, responseData)
		} else {
			c.JSON(http.StatusNoContent, gin.H{"message": "No content"})
		}
	}
}

func GetStockByIDApi(executeMethod reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ID is not correct"})
			return
		}
		var abc *model.Stock
		result := executeMethod.Call([]reflect.Value{reflect.ValueOf(uint(idInt))})
		if len(result) > 0 {
			responseData := result[0].Interface()
			if responseData == abc {
				c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			} else {
				c.JSON(http.StatusOK, responseData)
			}
		}
	}
}

// @Summary      create Stock api
// @Description  create new stock
// @Accept       json
// @Produce      json
// @Success      200  {array} model.Stock
// @Router       /api [get]
func CreateStockApi(executeMethod reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonObj := &model.Stock{}
		if err := c.ShouldBind(&jsonObj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request params not good",
			})
		}
		jsonObj.LastUpdate = time.Now()
		result := executeMethod.Call([]reflect.Value{reflect.ValueOf(jsonObj)})
		if len(result) > 0 { // > 0 mean execute success
			c.JSON(http.StatusCreated, gin.H{"message": "Success create stock"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No content"})
		}
	}
}

// @Summary      update Stock api
// @Description  update price of a stock
// @Accept       json
// @Produce      json
// @Param        q    query   string  false  "name search by q"  Format(email)
// @Success      200  {array} model.Stock
// @Router       /api [post]
func UpdatePriceStockApi(executeMethod reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonObj := &model.Stock{}
		if err := c.ShouldBind(jsonObj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Input Data"})
		}
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
		}
		jsonObj.StockID = uint(idInt)
		result := executeMethod.Call(
			[]reflect.Value{reflect.ValueOf(jsonObj.StockPrice), reflect.ValueOf(jsonObj.StockID)})
		if len(result) > 0 { // > 0 mean found data
			responseData := result[0].Interface()
			c.JSON(http.StatusOK, responseData)
		} else {
			c.JSON(http.StatusNoContent, gin.H{"message": "No content"})
		}
	}
}

func DeleteStockApi(executeMethod reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
		}
		result := executeMethod.Call([]reflect.Value{reflect.ValueOf(uint(idInt))})
		if len(result) > 0 { // > 0 mean found data
			responseData := result[0].Interface()
			c.JSON(http.StatusOK, responseData)
		} else {
			c.JSON(http.StatusNoContent, gin.H{"message": "No content"})
		}
	}
}
