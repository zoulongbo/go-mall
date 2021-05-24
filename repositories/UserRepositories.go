package repositories

import (
	"database/sql"
	"errors"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"strconv"
)

type User interface {
	Conn() error
	Select(username string) (user *models.User, err error)
	Insert(user *models.User) (id int64, err error)
	SelectByKey(id int64) (user *models.User, err error)
}

type UserManager struct {
	table string
	mysqlConn *sql.DB
}

func NewUserManager() User  {
	return &UserManager{
		table:     models.UserTable,
		mysqlConn: common.DB,
	}
}


func (u *UserManager) Conn() error {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
		common.DB = mysql
	}

	if u.table == "" {
		u.table = models.UserTable
	}
	return nil
}

func (u *UserManager) Select(username string) (user *models.User, err error) {
	userInfo := &models.User{}
	if err := u.Conn(); err != nil {
		return userInfo, err
	}
	sql := "SELECT * FROM `" + u.table + "` WHERE username=?"
	rows, err := u.mysqlConn.Query(sql, username)
	if err != nil {
		return userInfo, err
	}
	defer rows.Close()
	result := common.GetResultRow(rows)
	if len(result) < 1 {
		return userInfo, errors.New("用户不存在")
	}
	common.DataToStructByTag(result, userInfo, "sql")
	return userInfo, nil
}

func (u *UserManager) Insert(user *models.User) (id int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT `" + u.table + "` SET nickname=?, username=?, password=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(user.Nickname, user.Username, user.HashPassword)
	if err != nil {
		return
	}
	return result.LastInsertId()
}


func (u *UserManager) SelectByKey(id int64) (user *models.User, err error) {
	userInfo := &models.User{}
	if err := u.Conn(); err != nil {
		return userInfo, err
	}
	sql := "SELECT * FROM `" + u.table + "` WHERE id=" + strconv.FormatInt(id, 10)
	row, err := u.mysqlConn.Query(sql)
	if err != nil {
		return userInfo, err
	}
	result := common.GetResultRow(row)
	if len(result) < 1 {
		return userInfo, nil
	}
	defer row.Close()
	common.DataToStructByTag(result, userInfo, "sql")
	return userInfo, nil
}
