package jwt

import (
	"time"
	"utils/conf"
	myerror "utils/error"

	"github.com/dgrijalva/jwt-go"
)

// GenToken 生成一个Token
func GenToken(userModel *UserModel) (string, error) {
	claim := jwt.MapClaims{
		"key":      userModel.UserKey,
		"username": userModel.UserName,
		"timeunix": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err := token.SignedString([]byte(conf.GetByKey("tokenkey").(string)))
	if err != nil {
		return "", err
	}
	return tokens, nil
}

// headers: {
// 	'Authorization': 'Bearer ' + token
//   }

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetByKey("tokenkey").(string)), nil
	}
}

// CheckToken 检查Token
func CheckToken(tokenss string) (UserModel, error) {
	_user := new(UserModel)
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		return *_user, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return *_user, myerror.New("无法转换令牌")
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		return *_user, myerror.New("令牌无效")
	}

	_user.UserKey = claim["key"].(string)
	_user.UserName = claim["username"].(string)
	_user.TimeUnix = int64(claim["timeunix"].(float64))
	return *_user, nil
}

// Token 解析Token信息
func Token(userModel *UserModel) (map[string]interface{}, error) {
	// 验证用户信息
	// 现在用户数据不在DB中，先用配置文件临时存放
	var userStore map[string]interface{}
	if v := conf.Get("user", userModel.UserKey); v == nil {
		return nil, myerror.New("找不到用户对应信息")
	} else {
		userStore = v.(map[string]interface{})
	}

	if userStore["PassWord"].(string) != userModel.PassWord {
		return nil, myerror.New("用户或者密码错误")
	}

	// 填充用户模型，后面可能要放到缓存中。
	userModel.UserName = userStore["UserName"].(string)

	// 发放token
	token, err := GenToken(userModel)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"Token": token,
	}, nil
}
