package repositories

import (
	"database/sql"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"strconv"
)

type Order interface {
	Conn() error
	Insert(order *models.Order) (id int64, err error)
	Delete(id int64) bool
	Update(order *models.Order) error
	SelectByKey(id int64) (order *models.Order, err error)
	SelectAll() (orders []*models.Order, err error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

//类似构造方法
func NewOrderManager() Order {
	return &OrderManager{
		table:     models.OrderTable,
		mysqlConn: common.DB,
	}
}


func (p *OrderManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
		common.DB = mysql
	}

	if p.table == "" {
		p.table = models.OrderTable
	}
	return nil
}

func (o OrderManager) Insert(order *models.Order) (id int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT `" + o.table + "` SET user_id=?, product_id=?, order_status=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, models.OrderStatusWait)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (o OrderManager) Delete(id int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM `" + o.table + "` where id=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func (o OrderManager) Update(order *models.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}

	sql := "UPDATE `" + o.table + "` SET  user_id=?, product_id=?, order_status=? WHERE id=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, strconv.FormatInt(order.ID, 10))
	if err != nil {
		return err
	}
	return nil
}

func (o OrderManager) SelectByKey(id int64) (orderResult *models.Order, err error) {
	order := &models.Order{}
	if err := o.Conn(); err != nil {
		return order, err
	}
	sql := "SELECT * FROM `" + o.table + "` WHERE id=" + strconv.FormatInt(id, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return order, err
	}
	result := common.GetResultRow(row)
	if len(result) < 1 {
		return order, nil
	}
	defer row.Close()
	common.DataToStructByTag(result, order, "sql")
	return order, nil
}

func (o OrderManager) SelectAll() (orders []*models.Order, err error) {
	if err := o.Conn(); err != nil {
		return orders, err
	}
	sql := "SELECT * FROM `" + o.table + "`"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return orders, err
	}
	result := common.GetResultRows(rows)
	if len(result) < 1 {
		return orders, err
	}
	defer rows.Close()

	for _, row := range result {
		order := &models.Order{}
		common.DataToStructByTag(row, order, "sql")
		orders = append(orders, order)
	}
	return orders, nil
}

func (o OrderManager) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT o.id, p.product_name,o.order_status FROM `" + o.table + "` as o left join product as p on o.product_id=p.id"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result) < 1 {
		return nil, err
	}
	return result, nil
}