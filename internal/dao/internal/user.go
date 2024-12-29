// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao is the data access object for table user.
type UserDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns UserColumns // columns contains all the column names of Table for convenient usage.
}

// UserColumns defines and stores column names for table user.
type UserColumns struct {
	Id           string // ID
	Nickname     string // 用户昵称
	Username     string // 用户用户名ID
	Password     string // 用户密码
	Phone        string // 用户手机号
	Email        string // 用户邮箱
	RegisterTime string // 用户注册时间
	LoginTime    string // 用户最后登录时间
	Token        string // 用户token
	Group        string // 用户分组ID
	Avatar       string // 用户头像
	EmailAuth    string // 邮箱验证 0 未验证 1 验证
}

// userColumns holds the columns for table user.
var userColumns = UserColumns{
	Id:           "id",
	Nickname:     "nickname",
	Username:     "username",
	Password:     "password",
	Phone:        "phone",
	Email:        "email",
	RegisterTime: "register_time",
	LoginTime:    "login_time",
	Token:        "token",
	Group:        "group",
	Avatar:       "avatar",
	EmailAuth:    "email_auth",
}

// NewUserDao creates and returns a new DAO object for table data access.
func NewUserDao() *UserDao {
	return &UserDao{
		group:   "default",
		table:   "user",
		columns: userColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserDao) Columns() UserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
