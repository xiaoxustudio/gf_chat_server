// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocumentTemplateDao is the data access object for table document-template.
type DocumentTemplateDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns DocumentTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// DocumentTemplateColumns defines and stores column names for table document-template.
type DocumentTemplateColumns struct {
	Id      string // ID
	UserId  string // 用户名
	Auth    string // 用户权限：0 可查看 1可编辑 2可管理
	AddTime string // 添加时间
}

// documentTemplateColumns holds the columns for table document-template.
var documentTemplateColumns = DocumentTemplateColumns{
	Id:      "id",
	UserId:  "user_id",
	Auth:    "auth",
	AddTime: "add_time",
}

// NewDocumentTemplateDao creates and returns a new DAO object for table data access.
func NewDocumentTemplateDao() *DocumentTemplateDao {
	return &DocumentTemplateDao{
		group:   "default",
		table:   "document-template",
		columns: documentTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DocumentTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DocumentTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DocumentTemplateDao) Columns() DocumentTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DocumentTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DocumentTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DocumentTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
