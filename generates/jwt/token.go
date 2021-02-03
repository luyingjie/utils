package jwt

import (
	"net/http"
	"strings"
	"time"
	"utils/conv"
	verror "utils/os/error"

	"github.com/dgrijalva/jwt-go"
)

const (
	DEFAULT_Time_Unix int64 = 3600
)

// GenToken 生成一个Token
func GenToken(userModel *UserModel, key string) (string, error) {
	claim := jwt.MapClaims{
		"key":      userModel.UserKey,
		"username": userModel.UserName,
		"timeunix": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokens, nil
}

// headers: {
// 	'Authorization': 'Bearer ' + token
//   }

func secret(key string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}
}

// CheckToken 检查Token
func CheckToken(tokenss string, key string, timeunix ...int64) (UserModel, error) {
	_user := new(UserModel)
	token, err := jwt.Parse(tokenss, secret(key))
	if err != nil {
		return *_user, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return *_user, verror.New("无法转换令牌")
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		return *_user, verror.New("令牌无效")
	}

	// 判断时间戳
	tu := conv.Int64(claim["timeunix"])

	// 判断时效
	_timeUnix := DEFAULT_Time_Unix
	if len(timeunix) > 0 && timeunix[0] != 0 {
		_timeUnix = timeunix[0]
	}
	if _timeUnix != 0 {
		if time.Now().Unix()-tu > _timeUnix {
			return *_user, verror.New("用户认证信息失效！")
		}
	}
	if claim["key"].(string) == "" {
		return *_user, verror.New("用户认证信息错误！")
	}

	_user.UserKey = claim["key"].(string)
	_user.UserName = claim["username"].(string)
	_user.TimeUnix = tu
	return *_user, nil
}

// BearerAuth 从http中解析出Token
func BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}
