package user

import (
	"context"
	"fmt"
	"gf_chat_server/internal/consts"
	"gf_chat_server/internal/dao"
	"gf_chat_server/internal/model/entity"
	"gf_chat_server/utility/msgtoken"
	"gf_chat_server/utility/token"
	"gf_chat_server/utility/tw"
	"gf_chat_server/utility/verifiy"
	"gf_chat_server/utility/xtime"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

type User struct{}

func New() *User {
	return &User{}
}
func (c *User) ValidToken(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()
	tok := data["token"].(string)
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
	}
	vaild, err := token.ValidToken(tok)
	if vaild && err == nil {
		// 进行数据库比对
		r, err := md.Where("token", tok).All()
		if err == nil && len(r) > 0 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "OK", g.Map{
				"token": tok,
			})))
		} else {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
		}
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenExpired, "token校验失败:"+err.Error(), nil)))
	}
}

func (c *User) Register(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()

	tw.Tw(context.TODO(), "%v", data) // 打印
	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.DataEmpty, "空数据", nil)))
	}
	// 第一步：校验信息是否有效
	_, verifiyErr := verifiy.Exec(data, []string{})
	if verifiyErr != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证信息错误："+verifiyErr.Error(), nil)))
	}
	data["nickname"] = "默认用户" + (gconv.String(time.Now().Unix()))[0:6]
	data["register_time"] = xtime.NowDate()
	data["login_time"] = xtime.NowDate()
	data["group"] = 0
	// 第二步：校验用户是否存在
	usernameRes, _ := md.Where("username", data["username"]).All()
	if usernameRes.Len() > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "用户已存在！", nil)))
	}
	// 分发token
	token, _ := token.Token(data["username"].(string), 24)
	data["token"] = token
	// 通过
	_, err := md.Insert(data)
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "注册成功！", g.Map{"token": data["token"]})))
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "服务发生错误："+err.Error(), nil)))
	}
}
func (c *User) Login(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()

	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.DataEmpty, "空数据", nil)))
	}
	// 第一步：校验信息是否有效
	_, verifiyErr := verifiy.Exec(data, []string{"email"})
	if verifiyErr != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证信息错误："+verifiyErr.Error(), nil)))
	}
	username := data["username"].(string)
	password := data["password"].(string)
	usernameRes, err := md.Where("username", username).Where("password", password).All()
	if err == nil && usernameRes.Len() > 0 {
		// 分发token
		token, _ := token.Token(data["username"].(string), 24)
		md.Update(g.Map{"token": token})
		singleData := usernameRes[0].Map()
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "登录成功！", g.Map{
			"username": singleData["username"],
			"nickname": singleData["nickname"],
			"token":    token,
		})))
	} else {
		if err != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "登录失败！数据异常", nil)))
		} else {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "登录失败！账户或密码错误", nil)))
		}
	}
}

func (c *User) GetUser(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败:"+err.Error(), nil)))
	}
	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "空数据", nil)))
	}
	res, err := md.Where("username", data["user"]).All()
	if err == nil && res.Len() > 0 {
		singleData := res[0].Map()
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", g.Map{
			"nickname":      singleData["nickname"],
			"username":      singleData["username"],
			"phone":         singleData["phone"],
			"email":         singleData["email"],
			"register_time": singleData["register_time"],
			"login_time":    singleData["login_time"],
			"group":         singleData["group"],
		})))
	} else {
		if res.Len() == 0 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未查询到该用户！", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证失败", nil)))
	}
}
func (c *User) GetFriend(req *ghttp.Request) {
	data := req.GetFormMap()
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenExpired, "token校验失败:"+err.Error(), nil)))
	}
	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	var fData []entity.Friends
	md := dao.Friends.Ctx(req.Context())
	err = md.Where("user_id", data["user"].(string)).FieldsEx("friend_data.password").With(entity.User{}).Scan(&fData)
	if err == nil {
		var fDataEx = make([]*entity.Friends, len(fData)-1)
		for _, i := range fData {
			i.FriendData.Password = ""
			fDataEx = append(fDataEx, &i)
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "OK", fDataEx)))
	} else {
		if len(fData) == 0 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未找到联系人", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证失败", nil)))
	}
}
func (c *User) AddFriend(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenExpired, "token校验失败:"+err.Error(), nil)))
	}
	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	res, err := md.Where("username", data["user"]).All()
	if err == nil && res.Len() > 0 {
		jtoken, _ := token.ParseJwt(tok)
		if data["user"] != jtoken.Username {
			token_username, _ := token.ParseJwt(tok)
			res, _ = g.Model("friends").Where("user_id", token_username.Username).Where("friend_id", data["user"]).All()
			if res.Len() > 0 {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "已经添加过好友了!", nil)))
			}
			md = dao.Friends.Ctx(req.Context())
			_, err := md.Insert(g.Map{
				"user_id":   token_username.Username,
				"friend_id": data["user"],
				"add_time":  time.Now(),
			})
			if err != nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "添加好友失败！", nil)))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "添加好友成功！", nil)))
		} else {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "不能添加自己", nil)))
		}
	} else {
		if res.Len() == 0 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "对方不存在", nil)))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证失败", nil)))
	}
}
func (c *User) SearchUsers(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败:"+err.Error(), nil)))
	}
	if len(data) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "空数据", nil)))
	}
	res, err := md.WhereLike("nickname", fmt.Sprintf("%%%s%%", data["nickname"])).All()
	if err == nil && res.Len() > 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", res)))
	} else {
		if res.Len() == 0 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "没有搜索到！", [...]string{})))
		}
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证失败", nil)))
	}
}

func getTemp() string {
	return "./temp"
}

// 清空指定的目录
func clearTempDir(dir string) error {
	// 使用os.ReadDir读取目录下的所有文件和目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// 遍历文件和目录，并删除它们
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// 上传图片或头像
func (c *User) UploadImg(req *ghttp.Request) {
	files := req.GetUploadFiles("file")
	typeVal := req.GetFormMap()["type"]

	var targetPath string
	if typeVal == nil || typeVal != "avatar" {
		targetPath = getTemp()
	} else {
		targetPath = "./resource/avatar"
	}
	if files == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "请选择需要上传的文件", nil)))
	}

	if typeVal != "avatar" {
		maxFiles := 5 // 文件阈值
		// 检测目录文件数量是否超出，超出则清空
		TempFiles, err := os.ReadDir(targetPath)
		if err != nil {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "上传失败："+err.Error(), nil)))
			return
		}
		// 清空操作
		if len(TempFiles) > maxFiles {
			// 清空temp目录
			err = clearTempDir(targetPath)
			if err != nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "上传失败："+err.Error(), nil)))
				return
			}
		}

	}

	allowedMimeTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/jpg":  true,
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
