package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AkapongAlone/komgrip-test/models"
	"github.com/AkapongAlone/komgrip-test/requests"
	"github.com/AkapongAlone/komgrip-test/src/beers/domains"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type BeerHandler struct {
	beerUseCase domains.BeerUseCase
}

func NewBeerHandler(usecase domains.BeerUseCase) *BeerHandler {
	return &BeerHandler{beerUseCase: usecase}
}

// @summary
// @description
// @tags beer
// @Param name query string false "Insert Name"
// @Param limit query string false "Insert limit of items"
// @Param page query string false "Insert current page"
// @response 200 {object} responses.PaginationBody 
// @router /beer [get]
func (t *BeerHandler) GetBeer(c *gin.Context) {
	nameParam := c.Query("name")
	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	if limitParam == ""{
		limitParam = "10"
	}
	if pageParam == "" {
		pageParam ="1"
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	result, err := t.beerUseCase.GetBeer(nameParam, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// @summary
// @description
// @tags beer
// @param Beer body requests.Beer true "Beer data"
// @response 200 {object} responses.NoDataResponse "OK"
// @router /beer [post]
func (t *BeerHandler) PostBeer(c *gin.Context) {
	var beerInputData requests.Beer
	if err := c.ShouldBind(&beerInputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "hello",
		})
		return
	}

	var beer models.BeerDB
	copier.Copy(&beer, beerInputData)
	// file, header, err := c.Request.FormFile("file")
	file, err := c.FormFile("picture")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	beer.Picture, err = t.beerUseCase.Upload(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}
	fmt.Println(beer)

	err = t.beerUseCase.PostBeer(beer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// @summary
// @description
// @tags beer
// @param Beer body requests.Beer true "Beer data"
// @Param ID query string true "Insert ID"
// @response 200 {object} responses.NoDataResponse "OK"
// @router /beer [put]
func (t *BeerHandler) UpdateBeer(c *gin.Context) {
	var result requests.Beer
	if err := c.ShouldBind(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var beer models.BeerDB
	copier.Copy(&beer, result)
	// file, header, err := c.Request.FormFile("file")
	file, err := c.FormFile("picture")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	t.beerUseCase.Remove(id)
	beer.Picture, err = t.beerUseCase.Upload(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}
	
	
	
	fmt.Println(id)
	err = t.beerUseCase.UpdateBeer(id, beer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": beer,
	})
}

// @summary
// @description
// @tags beer
// @Param ID query string true "Insert ID"
// @response 200
// @router /beer [delete]
func (t *BeerHandler) DeleteBeer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	err = t.beerUseCase.DeleteBeer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Deleted",
	})
}
