// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TokensDao is the data access object for table tokens.
type TokensDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns TokensColumns // columns contains all the column names of Table for convenient usage.
}

// TokensColumns defines and stores column names for table tokens.
type TokensColumns struct {
	Id          string // ID
	Token       string // Token
	CreateTime  string // 创建时间
	FailureTime string // 失效时间
	TargetEmail string // 目标邮箱
	Ip          string // IP
}

// tokensColumns holds the columns for table tokens.
var tokensColumns = TokensColumns{
	Id:          "id",
	Token:       "token",
	CreateTime:  "create_time",
	FailureTime: "failure_time",
	TargetEmail: "target_email",
	Ip:          "ip",
}

// NewTokensDao creates and returns a new DAO object for table data access.
func NewTokensDao() *TokensDao {
	return &TokensDao{
		group:   "default",
		table:   "tokens",
		columns: tokensColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TokensDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TokensDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TokensDao) Columns() TokensColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TokensDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TokensDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TokensDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
