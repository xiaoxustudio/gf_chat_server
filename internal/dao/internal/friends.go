// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FriendsDao is the data access object for table friends.
type FriendsDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns FriendsColumns // columns contains all the column names of Table for convenient usage.
}

// FriendsColumns defines and stores column names for table friends.
type FriendsColumns struct {
	Id       string // ID
	UserId   string // 用户用户名
	FriendId string // 好友用户名
	AddTime  string // 添加时间
}

// friendsColumns holds the columns for table friends.
var friendsColumns = FriendsColumns{
	Id:       "id",
	UserId:   "user_id",
	FriendId: "friend_id",
	AddTime:  "add_time",
}

// NewFriendsDao creates and returns a new DAO object for table data access.
func NewFriendsDao() *FriendsDao {
	return &FriendsDao{
		group:   "default",
		table:   "friends",
		columns: friendsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FriendsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FriendsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FriendsDao) Columns() FriendsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *FriendsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FriendsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FriendsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
