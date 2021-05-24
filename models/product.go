package models

type Product struct {
	ID           int64  `json:"id" sql:"id" info:"id"`
	ProductName  string `json:"productName" sql:"product_name"  info:"productName"`
	ProductNum   int64  `json:"productNum" sql:"product_num" info:"productNum"`
	ProductImage string `json:"productImage" sql:"product_image" info:"productImage"`
	ProductUrl   string `json:"productUrl" sql:"product_url" info:"productUrl"`
}

const (
	ProductTable = "product"
)