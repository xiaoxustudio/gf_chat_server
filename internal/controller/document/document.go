package document

import (
	"encoding/json"
	"errors"
	"fmt"
	"gf_chat_server/internal/consts"
	"gf_chat_server/internal/dao"
	"gf_chat_server/internal/model/entity"
	"gf_chat_server/utility/msgtoken"
	"gf_chat_server/utility/rand"
	"gf_chat_server/utility/token"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type Document struct{}

// 页面和文件夹共用一个表，通过类型字段区分

func New() *Document {
	return &Document{}
}

func (d *Document) validToken(req *ghttp.Request) (*token.JWTToken, error) {
	md := dao.User.Ctx(req.Context())
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		return &token.JWTToken{}, errors.New("token校验失败！")
	}
	vaild, err := token.ValidToken(tok)
	if vaild && err == nil {
		// 进行数据库比对
		r, err := md.Where("token", tok).All()
		// 解析
		if err == nil && len(r) > 0 {
			tokVal, _ := token.ParseJwt(tok)
			return tokVal, nil
		} else {
			return &token.JWTToken{}, errors.New("token校验失败！")
		}
	} else {
		return &token.JWTToken{}, errors.New("token校验失败！" + err.Error())
	}
}

// 创建页面
func (d *Document) CreatePage(req *ghttp.Request) {
	tk, err := d.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("创建页面失败：%s", err.Error()), nil)))
	}
	if data == nil || data["type"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败：缺少参数", nil)))
	}

	typeVal, ok := data["type"].(int)
	if !ok {
		toInt, ok := strconv.Atoi(string(data["type"].(json.Number)))
		if ok != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败：类型参数类型错误", nil)))
		}
		typeVal = toInt
	}

	if !(typeVal == 0 || typeVal == 1 || typeVal == 2) {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败：类型错误:"+err.Error(), nil)))
	}
	var udata entity.User
	md := dao.User.Ctx(req.Context())
	err = md.Where("username", tk.Username).Scan(&udata)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "用户异常", nil)))
	}
	if udata.EmailAuth != 1 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "用户邮箱未验证，无法创建页面", nil)))
	}
	md = dao.Documents.Ctx(req.Context())
	document_id := rand.GenUniqueID()
	tableName := fmt.Sprintf("`document-%s`", document_id)
	res, err := md.Clone().Where("block", document_id).One()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面错误！", nil)))
	} else if len(res) > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面错误（初始化重复页面ID）", nil)))
	}
	_, err = md.Clone().Insert(entity.Documents{
		Block:     document_id,
		UserId:    tk.Username,
		Type:      typeVal,
		Status:    1,
		Content:   "",
		BlockDesc: "",
		BlockName: "无标题",
	})
	if err == nil {
		createTableSQL := `
CREATE TABLE ` + tableName + ` (
  ` + "`id`" + ` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  ` + "`user_id`" + ` varchar(255) NOT NULL COMMENT '用户名',
  ` + "`auth`" + ` int(11) NOT NULL COMMENT '用户权限：0 可查看 1可编辑 2可管理',
  ` + "`add_time`" + ` datetime(0) NOT NULL COMMENT '添加时间',
  PRIMARY KEY (` + "`id`" + `) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_unicode_ci ROW_FORMAT = Fixed;
`
		_, err := g.DB().Exec(createTableSQL)
		if err != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面错误！", nil)))
		} else {
			// 将自己添加进页面
			md := g.Model(tableName)
			_, err = md.Insert(g.Map{
				"user_id":  tk.Username,
				"add_time": gtime.Now(),
				"auth":     2,
			})
			if err == nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", g.Map{"block": document_id})))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面错误！", nil)))
		}
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面错误！", nil)))
	}
}

// 获取指定用户所有的page
func (d *Document) GetPages(req *ghttp.Request) {
	tk, err := d.validToken(req)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取页面失败：%s", err.Error()), nil)))
	}
	md := dao.Documents.Ctx(req.Context())
	var Ddata []entity.Documents
	err = md.Clone().WithAll().Scan(&Ddata)
	if err == nil {
		if len(Ddata) == 0 {
			Ddata = make([]entity.Documents, 0)
		}
		var DData = make([]entity.Documents, 0)
		for _, v := range Ddata {
			res, err := g.Model(fmt.Sprintf("document-%s", v.Block)).Where("user_id", tk.Username).All()
			v.DocData = res
			if v.UserId == tk.Username {
				// 如果是自己创建的
				DData = append(DData, v)
			} else {
				// 判断是否是分享给我们的
				if err == nil && len(res) > 0 {
					DData = append(DData, v)
				}
			}
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", DData)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取页面失败", nil)))
}

// 获取指定用户指定的page信息
func (d *Document) GetPage(req *ghttp.Request) {
	tk, err := d.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("创建页面失败：%s", err.Error()), nil)))
	}
	if data == nil || data["block"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败：缺少参数", nil)))
	}
	block_id := data["block"].(string)
	md := dao.Documents.Ctx(req.Context())
	res, err := md.Clone().Where(g.Map{"user_id": tk.Username, "block": block_id}).One()
	if err == nil {
		if res.IsEmpty() {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败：未找到页面", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建页面失败", nil)))
}

func (d *Document) DeletePage(req *ghttp.Request) {
	tk, err := d.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("删除页面失败：%s", err.Error()), nil)))
	}
	if data == nil || data["block"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除页面失败：缺少参数", nil)))
	}
	block_id := data["block"].(string)
	md := dao.Documents.Ctx(req.Context())
	res, err := md.Clone().Where(g.Map{"user_id": tk.Username, "block": block_id}).One()
	if err == nil {
		if res.IsEmpty() {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除页面失败：未找到页面", nil)))
		}
		_, err = md.Clone().Where(g.Map{"user_id": tk.Username, "block": block_id}).Delete()
		if err != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除页面失败：异常"+err.Error(), nil)))
		}
		sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `document-%s`;", block_id)
		_, err1 := g.DB().Exec(sqlStr)

		if err1 != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除页面失败：异常"+err1.Error(), nil)))
		}

		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除页面失败", nil)))
}

// 获取页面全部协作者
func (d *Document) GetPageCollaborators(req *ghttp.Request) {
	_, err := d.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取协作者失败：%s", err.Error()), nil)))
	}
	if data == nil || data["block"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取协作者失败：缺少参数", nil)))
	}
	block_id := data["block"].(string)
	tableName := fmt.Sprintf("document-%s", block_id)
	md := g.Model(tableName)
	res, err := md.Clone().All()
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取协作者失败："+err.Error(), nil)))
}

// 邀请协作
func (d *Document) InvitePeople(req *ghttp.Request) {
	tk, err := d.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("邀请失败：%s", err.Error()), nil)))
	}
	if data == nil || data["block"] == nil || data["users"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "邀请页面失败：缺少参数", nil)))
	}
	block_id := data["block"].(string)
	var members = make([]string, 0)
	for _, v := range data["users"].([]interface{}) {
		members = append(members, v.(string))
	}
	md := dao.Documents.Ctx(req.Context())
	var bdata entity.Documents
	err = md.Clone().Where(g.Map{"user_id": tk.Username, "block": block_id}).Scan(&bdata)
	if err == nil {
		tableName := fmt.Sprintf("document-%s", block_id)
		md := g.Model(tableName)
		_, err = md.Clone().Where("auth", 2).Where("user_id", tk.Username).One()
		if err == nil {
			for _, v := range members {
				_, err := g.Model("user").Where("username", v).One()
				if err != nil {
					// 用户不存在则跳过
					continue
				}
				_, err = md.Clone().Where("user_id", v).One()
				if err != nil {
					// 用户已存在在页面中则跳过
					continue
				}
				_, err = md.Clone().Insert(g.Map{
					"user_id":  v,
					"add_time": gtime.Now(),
					"auth":     1,
				})
				if err != nil {
					req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "邀请失败："+err.Error(), nil)))
				}
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "邀请失败：无管理权限", nil)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "邀请失败"+err.Error(), nil)))
}
