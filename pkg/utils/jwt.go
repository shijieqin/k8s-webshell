package utils

import (
	"k8s-webshell/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username      string `json:"username"`
	Password      string `json:"password`
	PodNs         string `json:"podNs"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`

	jwt.StandardClaims
}

func GenerateToken(username, password, podNs, podName, containerName string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		podNs,
		podName,
		containerName,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "k8s-webshell",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if Claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return Claims, nil
		}
	}
	return nil, err
}
