package group

import (
	"errors"
	"fmt"
	"gf_chat_server/internal/consts"
	"gf_chat_server/internal/dao"
	"gf_chat_server/internal/model/entity"
	"gf_chat_server/utility/msgtoken"
	"gf_chat_server/utility/rand"
	"gf_chat_server/utility/token"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type Group struct{}

// 群组数据表模板
type GroupTemplateTable struct {
	Id       int         `json:"id"`
	GroupId  int         `json:"group_id"`
	UserName int         `json:"username"`
	AddTime  *gtime.Time `json:"add_time"`
}

func New() *Group {
	return &Group{}
}
func (c *Group) validToken(req *ghttp.Request) (*token.JWTToken, error) {
	md := g.Model("user")
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

// 删除群组
func (c *Group) DeleteGroup(req *ghttp.Request) {
	tk, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("删除群组失败：%s", err.Error()), nil)))
	}

	group_id := data["group_id"]
	if group_id == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "请完善群号！", nil)))
	}
	group_id = group_id.(string)
	// 是否存在该群
	md := g.Model("groups")
	res, err := md.Where("group_id", group_id).One()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未找到该群！", nil)))
	}
	// 是否是群主
	group_master := res.GMap().Get("group_master").(string)
	if group_master != tk.Username {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "无权操作！", nil)))
	}
	md = dao.Groups.Ctx(req.Context())
	runDel, _ := md.Where("group_id", group_id).Delete()
	affected, err := runDel.RowsAffected()
	if err == nil && affected > 0 {
		_, err = g.Model(fmt.Sprintf("group-%s", group_id)).Delete()
		if err == nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", g.Map{"group_id": group_id})))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除群组错误！", nil)))
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除群组错误！", nil)))
	}
}

// 创建群组
func (c *Group) CreateGroup(req *ghttp.Request) {
	tk, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("创建群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组失败：缺少参数", nil)))
	}
	md := g.Model("groups")
	group_id := rand.GenUniqueID()
	tableName := fmt.Sprintf("`group-%s`", group_id)
	res, err := md.Where("group_id", group_id).One()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组错误！", nil)))
	} else if len(res) > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组错误！（初始化重复群组号）", nil)))
	}

	md = dao.Groups.Ctx(req.Context())
	md2 := dao.GroupConnect.Ctx(req.Context())
	_, err = md.Insert(entity.Groups{
		GroupId:     group_id,
		GroupStatus: 0,
		GroupDesc:   data["group_desc"].(string),
		GroupAvatar: data["group_avatar"].(string),
		GroupMaster: tk.Username,
		GroupName:   data["group_name"].(string),
	})
	_, err2 := md2.Insert(entity.GroupConnect{
		GroupId: group_id,
		UserId:  tk.Username,
		Auth:    2,
		AddTime: gtime.Now(),
	})
	if err == nil && err2 == nil {
		createTableSQL := `
CREATE TABLE ` + tableName + ` (
  ` + "`id`" + ` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  ` + "`user_id`" + ` varchar(255) NOT NULL COMMENT '用户名',
  ` + "`auth`" + ` int(11) NOT NULL COMMENT '用户权限：0 普通 1管理 2群主',
  ` + "`add_time`" + ` datetime(0) NOT NULL COMMENT '加入群聊时间',
  ` + "`last_chat_time`" + ` datetime(0) NULL DEFAULT NULL COMMENT '最后发言时间',
  PRIMARY KEY (` + "`id`" + `) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_unicode_ci ROW_FORMAT = Fixed;
`
		if _, err := g.DB().Exec(req.Context(), createTableSQL); err != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组错误！", nil)))
		} else {
			// 将自己添加进群组
			md = g.Model(tableName)
			_, err = md.Insert(g.Map{
				"user_id":  tk.Username,
				"add_time": gtime.Now(),
				"auth":     2,
			})
			if err == nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", g.Map{"group_id": group_id})))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组错误！", nil)))
		}
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "创建群组错误！", nil)))
	}
}

var allowedMimeTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
}

// 上传图片或头像（获取图片地址）
func (c *Group) UploadGroupAvatar(req *ghttp.Request) {
	files := req.GetUploadFiles("file")
	var targetPath = "./resource/group_avatar"
	if files == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "请选择需要上传的文件", nil)))
	}
	maxFileSize := int64(1 * 1024 * 1024) // 1MB文件大小限制
	// 检查文件类型和大小
	for _, file := range files {
		if !allowedMimeTypes[file.Header.Get("Content-Type")] {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "不支持的文件类型，只允许上传png、jpg、jpeg文件", nil)))
			return
		}
		if file.Size > maxFileSize {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "文件大小超出限制，最大允许1MB", nil)))
			return
		}
	}

	names, err := files.Save(targetPath)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "上传失败："+err.Error(), nil)))
	}
	cacheFilesStrings := make([]string, len(names)-1)
	for _, val := range names {
		cacheFilesStrings = append(cacheFilesStrings, targetPath+"/"+val)
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "上传成功", cacheFilesStrings)))
}

// 获取加入的群聊
func (c *Group) GetJoinGroup(req *ghttp.Request) {
	_, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败：缺少参数", nil)))
	}

	md := dao.GroupConnect.Ctx(req.Context())
	var gdata []*entity.GroupConnect
	err = md.Where("user_id", data["user"].(string)).With(entity.Groups{}).Scan(&gdata)
	if err == nil {
		if gdata == nil {
			gdata = make([]*entity.GroupConnect, 0)
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", gdata)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败", nil)))
}

// 根据群聊ID获取群聊信息
func (c *Group) GetGroup(req *ghttp.Request) {
	_, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败：缺少参数", nil)))
	}
	md := g.Model("groups")
	res, err := md.Where("group_id", data["group"].(string)).One()
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败!", nil)))
}

// 根据群聊ID获取群所有成员
func (c *Group) GetGroupMembers(req *ghttp.Request) {
	_, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取成员失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取成员失败：缺少参数", nil)))
	}
	md := g.Model("groups")
	group_id := data["group"].(string)
	res, err := md.Where("group_id", group_id).One()
	if err != nil || len(res) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未找到该群！!", nil)))
	}
	md = g.Model(fmt.Sprintf("group-%s", group_id))
	var adata []entity.GroupTemplate
	err = md.With(entity.User{}).Scan(&adata)
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", adata)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取成员失败!"+err.Error(), nil)))
}

// 根据群聊名称/ID搜索群聊
func (c *Group) SearchGroup(req *ghttp.Request) {
	_, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("获取群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败：缺少参数", nil)))
	}
	group_id := data["group"].(string)
	res, err := g.Model("groups").Where("group_id", group_id).All()
	if err == nil && len(res) > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	res, err = g.Model("groups").WhereLike("group_name", fmt.Sprintf("%%%s%%", group_id)).All()
	if err == nil && len(res) > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", res)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "获取群组失败!", nil)))
}

// 加入群组
func (c *Group) JoinGroup(req *ghttp.Request) {
	tok, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("加入群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "加入群组失败：缺少参数", nil)))
	}
	md := g.Model("groups")
	group_id := data["group"].(string)
	res, err := md.Where("group_id", group_id).One()
	if err == nil {
		if len(res) > 0 {
			md = g.Model(fmt.Sprintf("group-%s", group_id))
			res, err := md.Where("user_id", tok.Username).All()
			if len(res) > 0 && err == nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "已经加入进群了！", nil)))
			}
			_, err = md.Insert(&entity.GroupTemplate{
				UserId:  tok.Username,
				AddTime: gtime.Now()})
			_, err1 := g.Model("group-connect").Insert(&entity.GroupConnect{
				UserId:  tok.Username,
				GroupId: group_id,
				AddTime: gtime.Now(), Auth: 0})
			if err == nil && err1 == nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "加入群组成功！", nil)))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "加入群组失败！", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未找到该群！", nil)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "加入群组失败!", nil)))
}

// 根据群聊ID退出群聊
func (c *Group) ExitGroup(req *ghttp.Request) {
	tok, err := c.validToken(req)
	data := req.GetFormMap()
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, fmt.Sprintf("退出群组失败：%s", err.Error()), nil)))
	}
	if data == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "退出群组失败：缺少参数", nil)))
	}
	group_id := data["group"].(string)
	// 判断是否有该群聊
	resG, err := g.Model("groups").Where("group_id", group_id).All()
	if err == nil && len(resG) > 0 {
		// 判断是否加入了该群
		mdG := g.Model(fmt.Sprintf("group-%s", group_id))
		res, err := mdG.Where("user_id", tok.Username).All()
		if err == nil && len(res) > 0 {
			_, err := mdG.Delete("user_id", tok.Username)
			_, err1 := g.Model("group-connect").Where("user_id", tok.Username).Where("group_id", group_id).Delete()
			if err == nil && err1 == nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "ok", nil)))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "退出群组失败!", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "退出群组失败：未加入群聊!", nil)))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "退出群组失败：未找到群聊!", nil)))
}
