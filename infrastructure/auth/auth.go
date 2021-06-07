package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type AuthInterface interface {
	CreateAuth(uint64, *TokenDetail) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteToken(*AccessDetail) error
}

type ClientData struct {
	client *redis.Client
}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{
		client: client,
	}
}

var _ AuthInterface = &ClientData{}

type AccessDetail struct {
	TokenUUID string
	UserID    uint64
}
type TokenDetail struct {
	AccessToken  string
	RefreshToken string
	TokenUUID    string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func (cd *ClientData) CreateAuth(id uint64, token *TokenDetail) error {
	at := time.Unix(token.AtExpires, 0)
	rt := time.Unix(token.RtExpires, 0)
	now := time.Now()

	// save token on redis
	atCreated, err := cd.client.Set(token.TokenUUID, strconv.Itoa(int(id)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := cd.client.Set(token.RefreshUUID, strconv.Itoa(int(id)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

func (cd *ClientData) FetchAuth(tokenUUID string) (uint64, error) {
	userid, err := cd.client.Get(tokenUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}
func (cd *ClientData) DeleteRefresh(refreshToken string) error {
	deleted, err := cd.client.Del(refreshToken).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
func (cd *ClientData) DeleteToken(authD *AccessDetail) error {
	refershToken := fmt.Sprintf("%s++%d", authD.TokenUUID, authD.UserID)
	deletedAt, err := cd.client.Del(refershToken).Result()
	if err != nil {
		return err
	}
	deletedRt, err := cd.client.Del(authD.TokenUUID).Result()
	if err != nil {
		return err
	}
	if deletedAt == 0 || deletedRt == 0 {
		return errors.New("something went wrong")
	}
	return nil
}
