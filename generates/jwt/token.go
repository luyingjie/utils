package jwt

import (
	"time"
	"utils/conf"
	myerror "utils/error"

	"github.com/dgrijalva/jwt-go"
)

func GenToken(userModel *UserModel) string {
	claim := jwt.MapClaims{
		"key":      userModel.UserKey,
		"username": userModel.UserName,
		"timeunix": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err := token.SignedString([]byte(conf.GetByKey("tokenkey").(string)))
	if err != nil {
		myerror.Try(4000, 3, err)
	}
	return tokens
}

// headers: {
// 	'Authorization': 'Bearer ' + token
//   }

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetByKey("tokenkey").(string)), nil
	}
}

func CheckToken(tokenss string) UserModel {
	_user := new(UserModel)
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		myerror.Try(4000, 3, err)
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		myerror.Trys(4000, 3, "无法转换令牌")
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		myerror.Trys(4000, 3, "令牌无效")
	}

	_user.UserKey = claim["key"].(string)
	_user.UserName = claim["username"].(string)
	_user.TimeUnix = int64(claim["timeunix"].(float64))
	return *_user
}

func Token(userModel *UserModel) map[string]interface{} {
	// 验证用户信息
	// 现在用户数据不在DB中，先用配置文件临时存放
	var userStore map[interface{}]interface{}
	if v := conf.Get("user", userModel.UserKey); v == nil {
		myerror.Trys(4000, 3, "找不到用户对应信息")
	} else {
		userStore = v.(map[interface{}]interface{})
	}

	if userStore["PassWord"].(string) != userModel.PassWord {
		myerror.Trys(4000, 3, "用户或者密码错误")
	}

	// 填充用户模型，后面可能要放到缓存中。
	userModel.UserName = userStore["UserName"].(string)

	// 发放token
	token := GenToken(userModel)
	return map[string]interface{}{
		"Token": token,
	}
}
