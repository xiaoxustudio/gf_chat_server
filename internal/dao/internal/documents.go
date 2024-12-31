// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DocumentsDao is the data access object for table documents.
type DocumentsDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns DocumentsColumns // columns contains all the column names of Table for convenient usage.
}

// DocumentsColumns defines and stores column names for table documents.
type DocumentsColumns struct {
	Id        string // ID
	Block     string // blockID
	BlockDesc string // 页面简介
	BlockName string // 页面名称
	UserId    string // 属于用户
	Type      string // 类型 0 文件夹 1 页面
	Content   string // 页面内容
	Status    string // block状态 0 锁定 1 可编辑
}

// documentsColumns holds the columns for table documents.
var documentsColumns = DocumentsColumns{
	Id:        "id",
	Block:     "block",
	BlockDesc: "block_desc",
	BlockName: "block_name",
	UserId:    "user_id",
	Type:      "type",
	Content:   "content",
	Status:    "status",
}

// NewDocumentsDao creates and returns a new DAO object for table data access.
func NewDocumentsDao() *DocumentsDao {
	return &DocumentsDao{
		group:   "default",
		table:   "documents",
		columns: documentsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DocumentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DocumentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DocumentsDao) Columns() DocumentsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DocumentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DocumentsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DocumentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
