package cmd

import (
	"context"
	"log"
	"obwallet/obrpc/user"
	"testing"
)

func TestSignUp(t *testing.T) {
	in := &user.SignUpRequest{
		UserName:        "wxf",
		Password:        "1234",
		ConfirmPassword: "1234",
		Email:           "wxf4150@163.com",
		Vcode:           "895036",
		VerifyCodeId:    "kTz8COKOGVS1dDjs4jDY",
	}
	res, err := gClient.SignUp(context.TODO(), in)
	if err != nil {
		log.Fatalln(err)
	}
	printRespJSON(res)
}
