package services

import (
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/repositories"
)

type Product interface {
	GetProductById(id int64) (product *models.Product, err error)
	GetAllProduct() (products []*models.Product, err error)
	InsertProduct(product *models.Product) (id int64, err error)
	DeleteProductById(id int64) bool
	UpdateProduct(product *models.Product) error
	SubNumOne(id int64, count int) error
}

type ProductService struct {
	productRepos repositories.Product
}

func NewProductService() Product {
	return &ProductService{productRepos: repositories.NewProductManager()}
}

func (p *ProductService) GetProductById(id int64) (product *models.Product, err error) {
	return p.productRepos.SelectByKey(id)
}

func (p *ProductService) GetAllProduct() (products []*models.Product, err error) {
	return p.productRepos.SelectAll()
}

func (p *ProductService) InsertProduct(product *models.Product) (id int64, err error) {
	return p.productRepos.Insert(product)
}

func (p *ProductService) DeleteProductById(id int64) bool {
	return p.productRepos.Delete(id)
}

func (p *ProductService) UpdateProduct(product *models.Product) error {
	return p.productRepos.Update(product)
}

func (p *ProductService) SubNumOne(id int64, count int) error {
	return p.productRepos.SubProductNum(id, count)
}
