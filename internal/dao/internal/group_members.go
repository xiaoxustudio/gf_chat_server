// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupMembersDao is the data access object for table group_members.
type GroupMembersDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns GroupMembersColumns // columns contains all the column names of Table for convenient usage.
}

// GroupMembersColumns defines and stores column names for table group_members.
type GroupMembersColumns struct {
	Id      string // ID
	UserId  string // 用户名称
	GroupId string // 加入组ID
	AddTime string // 加入时间
	Auth    string // 权限：0普通，1管理员，2群主
}

// groupMembersColumns holds the columns for table group_members.
var groupMembersColumns = GroupMembersColumns{
	Id:      "id",
	UserId:  "user_id",
	GroupId: "group_id",
	AddTime: "add_time",
	Auth:    "auth",
}

// NewGroupMembersDao creates and returns a new DAO object for table data access.
func NewGroupMembersDao() *GroupMembersDao {
	return &GroupMembersDao{
		group:   "default",
		table:   "group_members",
		columns: groupMembersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GroupMembersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GroupMembersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GroupMembersDao) Columns() GroupMembersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GroupMembersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GroupMembersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GroupMembersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
