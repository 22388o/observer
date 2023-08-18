package user

import (
	"github.com/lithammer/shortuuid/v4"
	"obwallet/utils"
	"time"
)

func (u *UserInfo) GenToken() (string, error) {
	tokenId := shortuuid.New()
	tokenObj := &Token{
		Token:  tokenId,
		UserId: u.Id,
	}
	utils.Gdb.Exec("update tokens set disabled =1 where user_id=? and updated_at=now()", u.Id)

	err := utils.Gdb.Create(tokenObj)
	if err != nil {
		return "", nil
	}
	return tokenId, nil
}

type Token struct {
	Id         int
	Token      string
	UserId     int64
	Updated_At time.Time
	Disabled   bool
}
