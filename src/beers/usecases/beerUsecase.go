package usecases

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/AkapongAlone/komgrip-test/models"
	_ "github.com/AkapongAlone/komgrip-test/requests"
	"github.com/AkapongAlone/komgrip-test/responses"
	"github.com/AkapongAlone/komgrip-test/src/beers/domains"
	_ "github.com/swaggo/swag/example/celler/httputil"
)

type beerUseCase struct {
	beerRepo domains.BeerRepositories
}

func NewBeerUseCase(repo domains.BeerRepositories) domains.BeerUseCase {
	return &beerUseCase{
		beerRepo: repo,
	}
}

func (t *beerUseCase) GetBeer(name string, limit int, page int) (responses.PaginationBody, error) {
	var items []responses.ItemBody
	offset := (page - 1) * limit
	result, err := t.beerRepo.FindData(name, limit, offset)
	if err != nil {
		return responses.PaginationBody{}, err
	}
	// err = godotenv.Load(".env")
	// if err != nil {
	// 	log.Panicln("consider env var")
	// }
	host := os.Getenv("HOST")
	for _, item := range result {
		itemBody := responses.ItemBody{
			Created_at: item.CreatedAt.String(),
			ID:         item.ID,
			Name:       name,
			Picture:    host + "/file/" + item.Picture,
			Detail:     item.Detail,
			Type:       item.Type,
			Updated_at: item.UpdatedAt.String(),
		}
		items = append(items, itemBody)
	}
	totalItem := t.beerRepo.CountAllData(name)
	totalPage := totalItem / limit
	switch {
	case totalPage == 0:
		totalPage = 1
	case (totalItem%limit != 0):
		totalPage++
	}

	var nextPage, prevPage int
	if page == totalPage {
		nextPage = page
	} else {
		nextPage = page + 1
	}
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = page
	}
	respond := responses.PaginationBody{
		CurrentPage:  page,
		Items:        items,
		NextPage:     nextPage,
		PreviousPage: prevPage,
		SizePerPage:  limit,
		TotalItems:   totalItem,
		TotalPages:   totalPage,
	}
	return respond, nil
}

func (t *beerUseCase) PostBeer(result models.BeerDB) error {

	err := t.beerRepo.InsertData(result)
	if err != nil {
		return err
	}
	return nil

}

func (t *beerUseCase) UpdateBeer(id int, result models.BeerDB) error {
	// var nameBeer models.BeerDB

	err := t.beerRepo.EditData(id, result)
	if err != nil {
		return err
	}
	return nil
}

func (t *beerUseCase) DeleteBeer(id int) error {
	err := t.beerRepo.DeleteData(id)
	if err != nil {
		return err
	}
	return nil
}

func (t *beerUseCase) Upload(header *multipart.FileHeader) (string, error) {

	fileName := header.Filename
	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	out, err := os.Create("./images/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}
	// filePath := "http://localhost:8080/file/" + fileName
	return fileName, nil
}

func (t *beerUseCase) Remove(id int) error {

	fileName, err := t.beerRepo.FindPicture(id)
	fmt.Println(fileName)
	if err != nil {
		return err
	}
	err = os.Remove("./images/" + fileName)
	if err != nil {
		return err
	}
	return nil
}
