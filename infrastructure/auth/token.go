package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type Token struct{}

func NewToken() *Token {
	return &Token{}
}

type TokenInterface interface {
	CreateToken(userid uint64) (*TokenDetail, error)
	ExtractTokenMeta(*http.Request) (*AccessDetail, error)
}

var _ TokenInterface = &Token{}

func (t *Token) CreateToken(userid uint64) (*TokenDetail, error) {
	td := &TokenDetail{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.TokenUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 7 * 24).Unix()
	td.RefreshToken = fmt.Sprintf("%s++%d", td.TokenUUID, userid)

	var err error
	// create access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.TokenUUID
	atClaims["user_id"] = userid
	atClaims["axp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	// create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(req *http.Request) string {
	bearToken := req.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
func VerifyToken(req *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted siging method: %s", t.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func TokenValid(req *http.Request) error {
	token, err := VerifyToken(req)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
func (t *Token) ExtractTokenMeta(req *http.Request) (*AccessDetail, error) {
	token, err := VerifyToken(req)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_token"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetail{TokenUUID: accessUUID, UserID: userID}, nil
	}
	return nil, err
}
