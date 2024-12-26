// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupTemplateDao is the data access object for table group-template.
type GroupTemplateDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns GroupTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// GroupTemplateColumns defines and stores column names for table group-template.
type GroupTemplateColumns struct {
	Id           string // ID
	UserId       string // 用户名
	Auth         string // 用户权限：0 普通 1管理 2群主
	AddTime      string // 加入群聊时间
	LastChatTime string // 最后发言时间
}

// groupTemplateColumns holds the columns for table group-template.
var groupTemplateColumns = GroupTemplateColumns{
	Id:           "id",
	UserId:       "user_id",
	Auth:         "auth",
	AddTime:      "add_time",
	LastChatTime: "last_chat_time",
}

// NewGroupTemplateDao creates and returns a new DAO object for table data access.
func NewGroupTemplateDao() *GroupTemplateDao {
	return &GroupTemplateDao{
		group:   "default",
		table:   "group-template",
		columns: groupTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GroupTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GroupTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GroupTemplateDao) Columns() GroupTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GroupTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GroupTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GroupTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
