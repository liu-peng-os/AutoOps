package jwt

import (
	"dodevops-api/api/system/model"
	"dodevops-api/common/constant"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type userStdClaims struct {
	model.JwtAdmin
	jwt.RegisteredClaims
}

const TokenExpireDuration = time.Hour * 24

var Secret = []byte("dodevops-api")

var (
	ErrAbsent  = "token absent"
	ErrInvalid = "token invalid"
)

func GenerateTokenByAdmin(admin model.SysAdmin) (string, error) {
	jwtAdmin := model.JwtAdmin{
		ID:       admin.ID,
		Username: admin.Username,
		Nickname: admin.Nickname,
		Icon:     admin.Icon,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Note:     admin.Note,
	}
	c := userStdClaims{
		jwtAdmin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(Secret)
}

func ValidateToken(tokenString string) (*model.JwtAdmin, error) {
	if tokenString == "" {
		return nil, errors.New(ErrAbsent)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if token == nil {
		return nil, errors.New(ErrInvalid)
	}
	claims := userStdClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	return &claims.JwtAdmin, err
}

func GetAdminId(c *gin.Context) (uint, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return 0, errors.New("can't get admin id")
	}
	admin, ok := u.(*model.JwtAdmin)
	if ok {
		return admin.ID, nil
	}
	return 0, errors.New("context user is not JwtAdmin")
}

func GetAdminName(c *gin.Context) (string, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return "", errors.New("can't get admin name")
	}
	admin, ok := u.(*model.JwtAdmin)
	if ok {
		return admin.Username, nil
	}
	return "", errors.New("context user is not JwtAdmin")
}

func GetAdmin(c *gin.Context) (*model.JwtAdmin, error) {
	u, exist := c.Get(constant.ContextKeyUserObj)
	if !exist {
		return nil, errors.New("can't get api")
	}
	admin, ok := u.(*model.JwtAdmin)
	if ok {
		return admin, nil
	}
	return nil, errors.New("context user is not JwtAdmin")
}
