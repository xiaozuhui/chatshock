package utils

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var JWTSecret = []byte("6c9c48b1-1ef6d091-90491a14-b3914756")

// TokenExpireDuration token失效时间默认为14天
const TokenExpireDuration = time.Hour * 24 * 14

// RefreshExpireDuration refresh失效时间为21天
const RefreshExpireDuration = time.Hour * 24 * 21

type UserClaims struct {
	PhoneNumber string
	jwt.RegisteredClaims
}

type UserRefreshClaims struct {
	PhoneNumber string
	IsFresh     bool
	jwt.RegisteredClaims
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	}
}

// GetToken
/**
 * @description: 在header中获取Bearer Token
 * @param {string} authHeader
 * @return {*}
 * @author: xiaozuhui
 */
func GetToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.WithStack(errors.New("无权限访问，请求未携带token"))
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.WithStack(errors.New("请求头中auth格式有误"))
	}
	return parts[1], nil
}

// GenerateToken
/**
 * @description: 生成Token
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func GenerateToken(phoneNumber string) (string, string, *time.Time, error) {
	token, expireTime, err := generateToken(phoneNumber)
	if err != nil {
		return "", "", nil, err
	}
	refreshToken, err := generateRefreshToken(phoneNumber)
	if err != nil {
		return "", "", nil, err
	}
	return token, refreshToken, &expireTime, nil
}

/** generateToken
 * @description: 生成Token
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func generateToken(phoneNumber string) (string, time.Time, error) {
	expireTime := time.Now().Add(TokenExpireDuration)
	claims := &UserClaims{
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
			Issuer:    "ChatShock",
		},
	}
	// 生成Token，指定签名算法和claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	tokenString, err := token.SignedString(JWTSecret)
	return tokenString, expireTime, err
}

/** generateRefreshToken
 * @description: 生成刷新用Token
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func generateRefreshToken(phoneNumber string) (string, error) {
	claims := &UserRefreshClaims{
		PhoneNumber: phoneNumber,
		IsFresh:     true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
			Issuer:    "ChatShock",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	return tokenString, err
}

// ParseToken
/**
 * @description: 解析Token
 * @param {string} tokenStr
 * @return {*}
 * @author: xiaozuhui
 */
func ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// ParseRefreshToken
/**
 * @description: 解析用于刷新的refreshToken
 * @param {string} tokenStr
 * @return {*}
 * @author: xiaozuhui
 */
func ParseRefreshToken(tokenStr string) (*UserRefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserRefreshClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*UserRefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// RefreshToken
/**
 * @description: 刷新token
 * @return {*}
 * @author: xiaozuhui
 */
func RefreshToken(tokenStr string) (string, error) {
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return "", err
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
		Issuer:    "ChatShock",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	return tokenString, err
}
