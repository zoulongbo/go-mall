package repositories

import (
	"database/sql"
	"errors"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"strconv"
)

type Product interface {
	Conn() error
	Insert(product *models.Product) (id int64, err error)
	Delete(id int64) bool
	Update(product *models.Product) error
	SelectByKey(id int64) (product *models.Product, err error)
	SelectAll() (products []*models.Product, err error)
	SubProductNum(productId int64, count int) error
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

//类似构造方法
func NewProductManager() Product {
	return &ProductManager{
		table:     models.ProductTable,
		mysqlConn: common.DB,
	}
}


func (p *ProductManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
		common.DB = mysql
	}

	if p.table == "" {
		p.table = models.ProductTable
	}
	return nil
}

func (p ProductManager) Insert(product *models.Product) (id int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	sql := "INSERT " + p.table + " SET product_name=?, product_num=?, product_image=?, product_url=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (p ProductManager) Delete(id int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM " + p.table + " where id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func (p ProductManager) Update(product *models.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "UPDATE " + p.table + " SET product_name=?, product_num=?, product_image=?, product_url=? WHERE id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, strconv.FormatInt(product.ID, 10))
	if err != nil {
		return err
	}
	return nil
}

func (p ProductManager) SelectByKey(id int64) (productResult *models.Product, err error) {
	product := &models.Product{}
	if err := p.Conn(); err != nil {
		return product, err
	}
	sql := "SELECT * FROM " + p.table + " WHERE id=" + strconv.FormatInt(id, 10)
	row, err := p.mysqlConn.Query(sql)
	if err != nil {
		return product, err
	}
	result := common.GetResultRow(row)
	if len(result) < 1 {
		return product, errors.New("商品详情不存在")
	}
	defer row.Close()
	common.DataToStructByTag(result, product, "sql")
	return product, nil
}

func (p ProductManager) SelectAll() (products []*models.Product, err error) {
	if err := p.Conn(); err != nil {
		return products, err
	}
	sql := "SELECT * FROM " + p.table
	rows, err := p.mysqlConn.Query(sql)
	if err != nil {
		return products, err
	}
	result := common.GetResultRows(rows)
	if len(result) < 1 {
		return products, err
	}
	defer rows.Close()

	for _, row := range result {
		product := &models.Product{}
		common.DataToStructByTag(row, product, "sql")
		products = append(products, product)
	}
	return products, nil
}

func (p ProductManager) SubProductNum(productId int64, count int) error  {
	if err := p.Conn(); err != nil {
		return err
	}
	sql  := "update " + p.table + " set product_num=product_num-" + strconv.Itoa(count) + " where id=" + strconv.FormatInt(productId, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}