// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupConnectDao is the data access object for table group-connect.
type GroupConnectDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns GroupConnectColumns // columns contains all the column names of Table for convenient usage.
}

// GroupConnectColumns defines and stores column names for table group-connect.
type GroupConnectColumns struct {
	Id      string // ID
	UserId  string // 用户名
	GroupId string // 群组名称
	Auth    string // 群组权限
	AddTime string // 加入时间
}

// groupConnectColumns holds the columns for table group-connect.
var groupConnectColumns = GroupConnectColumns{
	Id:      "id",
	UserId:  "user_id",
	GroupId: "group_id",
	Auth:    "auth",
	AddTime: "add_time",
}

// NewGroupConnectDao creates and returns a new DAO object for table data access.
func NewGroupConnectDao() *GroupConnectDao {
	return &GroupConnectDao{
		group:   "default",
		table:   "group-connect",
		columns: groupConnectColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GroupConnectDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GroupConnectDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GroupConnectDao) Columns() GroupConnectColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GroupConnectDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GroupConnectDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GroupConnectDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
