package user

import (
	"errors"
	"fmt"
	"gf_chat_server/internal/consts"
	emailsend "gf_chat_server/internal/controller/emailSend"
	"gf_chat_server/internal/dao"
	"gf_chat_server/internal/model/entity"
	"gf_chat_server/utility/iptool"
	"gf_chat_server/utility/msgtoken"
	"gf_chat_server/utility/rand"
	"gf_chat_server/utility/token"
	"gf_chat_server/utility/verifiy"
	"gf_chat_server/utility/xtime"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type User struct {
	EmailIns emailsend.EmailSendIns
}

func New() *User {
	return &User{EmailIns: emailsend.New()}
}

func (c *User) ValidToken(req *ghttp.Request) {
	md := g.Model("user")
	data := req.GetFormMap()
	tok := data["token"].(string)
	if data["token"] == nil || len(tok) == 0 {
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
	if len(data) == 0 || data["user"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "空数据", nil)))
	}
	user_id := data["user"].(string)
	res, err := md.Where("username", user_id).All()
	if err == nil && res.Len() > 0 {
		singleData := res[0].Map()
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", g.Map{
			"nickname":      singleData["nickname"],
			"username":      singleData["username"],
			"phone":         singleData["phone"],
			"email":         singleData["email"],
			"avatar":        singleData["avatar"],
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
	if len(data) == 0 || data["user"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	user_id := data["user"].(string)
	var fData []entity.Friends
	md := dao.Friends.Ctx(req.Context())
	err = md.Where("user_id", user_id).FieldsEx("friend_data.password").With(entity.User{}).Scan(&fData)
	if err == nil {
		var fDataEx []*entity.Friends
		if len(fData) > 0 {
			fDataEx = make([]*entity.Friends, len(fData)-1)
			for _, i := range fData {
				i.FriendData.Password = ""
				fDataEx = append(fDataEx, &i)
			}
		} else {
			fDataEx = make([]*entity.Friends, 0)
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
	if len(data) == 0 || data["user"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	user_id := data["user"].(string)
	res, err := md.Where("username", user_id).All()
	if err == nil && res.Len() > 0 {
		jtoken, _ := token.ParseJwt(tok)
		if user_id != jtoken.Username {
			res, _ = g.Model("friends").Where("user_id", jtoken.Username).Where("friend_id", user_id).All()
			if res.Len() > 0 {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "已经添加过好友了!", nil)))
			}
			md = dao.Friends.Ctx(req.Context())
			_, err := md.Insert(g.Map{
				"user_id":   jtoken.Username,
				"friend_id": user_id,
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

func (c *User) DeleteFriend(req *ghttp.Request) {
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
	if len(data) == 0 || data["user"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	user_id := data["user"].(string)
	res, err := md.Where("username", user_id).All()
	if err == nil && res.Len() > 0 {
		jtoken, _ := token.ParseJwt(tok)
		if user_id != jtoken.Username {
			res, _ = g.Model("friends").Where("user_id", jtoken.Username).Where("friend_id", user_id).All()
			if res.Len() == 0 {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "未添加该好友!", nil)))
			}
			md = dao.Friends.Ctx(req.Context())
			_, err := md.Delete(g.Map{
				"user_id":   jtoken.Username,
				"friend_id": user_id,
			})
			if err != nil {
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "删除好友失败！", nil)))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "删除好友成功！", nil)))
		} else {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "不能删除自己", nil)))
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
	if len(data) == 0 || data["nickname"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "空数据", nil)))
	}
	nickname := data["nickname"].(string)
	res, err := md.WhereLike("nickname", fmt.Sprintf("%%%s%%", nickname)).All()
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

var allowedMimeTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
}

// 上传图片或头像（获取图片地址）
func (c *User) UploadImg(req *ghttp.Request) {
	files := req.GetUploadFiles("file")
	var targetPath = getTemp()
	if files == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "请选择需要上传的文件", nil)))
	}
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

// 修改头像
func (c *User) ChangeAvatar(req *ghttp.Request) {
	md := g.Model("user")
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "token校验失败:"+err.Error(), nil)))
	}
	files := req.GetUploadFiles("file")
	if files == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "请选择需要上传的文件", nil)))
	}

	targetPath := "./resource/avatar"
	names, err := files.Save(targetPath)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "上传失败："+err.Error(), nil)))
	}
	cacheFilesStrings := make([]string, len(names)-1)
	for _, val := range names {
		cacheFilesStrings = append(cacheFilesStrings, targetPath+"/"+val)
	}

	_, err = md.Where("token", tok).Update(g.Map{"avatar": cacheFilesStrings[0]})
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(1, "ok", cacheFilesStrings[0])))
	} else {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "修改失败！"+err.Error(), nil)))
	}
}

// 指定邮件人发送验证码（内部）
func (c *User) sendEmailToken(req *ghttp.Request, addr string) (string, error) {
	ip, err := iptool.GetIP(req.Request)
	if err != nil {
		return "", errors.New("发送邮件失败！")
	}
	tk := rand.GetID(6)
	if len(tk) == 0 {
		return "", errors.New("发送邮件失败！")
	}
	md := g.Model("tokens")
	var tdata []entity.Tokens
	err = md.Clone().Where("target_email", addr).WithAll().Scan(&tdata)
	if err != nil {
		return "", errors.New("发送邮件失败！")
	}
	if len(tdata) == 3 {
		_, err = md.Clone().Insert(entity.Tokens{
			Token:       tk,
			CreateTime:  gtime.Now(),
			FailureTime: gtime.Now().Add(time.Duration(60) * time.Second), // 1分钟
			TargetEmail: addr,
			Ip:          ip,
		})
		if err != nil {
			return "", errors.New("发送邮件失败Db：超过次数！")
		}
		return "", errors.New("发送邮件失败：超过次数！")
	} else if len(tdata) == 4 {
		// 判断是否到1分钟，1分钟后可重新发送验证码
		dur_time := time.Duration(60) * time.Second
		lastRow := tdata[len(tdata)-1]
		timeInstan := time.Now().Add(dur_time)
		p := lastRow.CreateTime.Time.Before(timeInstan)
		if !p {
			return "", errors.New("发送邮件失败：请过1分钟后再试！")
		}
		// 先清空全部，再发送
		_, err = md.Clone().Where("target_email", addr).Delete()
		if err != nil {
			return "", errors.New("发送邮件失败：请过1分钟后再试！")
		}
		return c.sendEmailToken(req, addr)
	}
	// 发送邮件
	c.EmailIns.Send(addr, fmt.Sprintf("您的验证码为：%s \n IP : %s \n 1分钟内有效", tk, ip))
	// 记录token
	_, err = md.Clone().Insert(entity.Tokens{
		Token:       tk,
		CreateTime:  gtime.Now(),
		FailureTime: gtime.Now().Add(time.Duration(60) * time.Second), // 1分钟
		TargetEmail: addr,
		Ip:          ip,
	})
	if err != nil {
		return "", errors.New("发送邮件失败Db")
	}
	return ip, nil
}

// 向指定邮箱发送验证码链接
func (c *User) SendEmail(req *ghttp.Request) {
	data := req.GetFormMap()
	if len(data) == 0 || data["email"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	email := data["email"].(string)
	res, err := c.sendEmailToken(req, email)
	if err == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "发送成功!", g.Map{"ip": res})))
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "发送失败："+err.Error(), nil)))
}

// 验证邮箱
func (c *User) ValidEmail(req *ghttp.Request) {
	data := req.GetFormMap()
	tok := req.Header.Get("Authorization")
	if len(tok) == 0 {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenInValid, "token校验失败！", nil)))
	}
	_, err := token.ValidToken(tok)
	if err != nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenExpired, "token校验失败:"+err.Error(), nil)))
	}
	if len(data) == 0 || data["token"] == nil || data["user"] == nil {
		req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "空数据", nil)))
	}
	user_id := data["user"].(string)
	token := data["token"].(string)
	userMd := dao.User.Ctx(req.Context())
	tokensMd := dao.Tokens.Ctx(req.Context())
	var userSingle entity.User
	err = userMd.Clone().Where("username", user_id).Scan(&userSingle)
	if err == nil {
		if email_auth := userSingle.EmailAuth; email_auth == 1 {
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "该用户已验证！", nil)))
		} else {
			// 未验证
			email := userSingle.Email
			var tokenRes entity.Tokens
			err := tokensMd.Clone().Where("token", token).Where("target_email", email).Scan(&tokenRes)
			if err == nil {
				// 判断有效期
				createTime := tokenRes.CreateTime
				failureTime := tokenRes.FailureTime // 失效期
				if !createTime.Before(failureTime) {
					req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证码过期！", nil)))
				}
				// 说明验证码，邮箱，有效期都是对的
				_, err = tokensMd.Clone().Where("target_email", email).Delete()
				_, err1 := userMd.Clone().Where("username", user_id).Update(g.Map{"email_auth": 1})
				if err1 == nil && err == nil {
					req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.Success, "验证成功！", nil)))
				}
				req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证失败！", nil)))
			}
			req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(0, "验证码错误！"+err.Error(), nil)))
		}
	}
	req.Response.WriteJsonExit(msgtoken.ToGMap(msgtoken.MsgToken(consts.TokenEmpty, "未找到该用户！"+err.Error(), nil)))
}
