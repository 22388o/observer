package obrpc

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"obwallet/obrpc/user"
	"obwallet/utils"
	"time"
)

// A compile time check to ensure that Server fully implements the RouterServer
// gRPC service.
var _ user.UserServiceServer = (*UserServer)(nil)

type UserServer struct {
	user.UnsafeUserServiceServer
	kycUrl string
	payUrl string
	merNo  string
	merKey string
}

func NewUserServer(kycUrl, payUrl, merNo, merKey string) *UserServer {
	InitTables()
	return &UserServer{
		kycUrl: kycUrl,
		payUrl: payUrl,
		merNo:  merNo,
		merKey: merKey,
	}
}

func InitTables() {
	utils.Gdb.AutoMigrate(user.Token{}, user.UserInfo{}, user.CardInfo{}, user.CardExchangeInfo{}, UploadFile{})
}
func (us *UserServer) SignUp(ctx context.Context, in *user.SignUpRequest) (*user.SignUpResponse, error) {
	//check vcode
	if !utils.VerifyCaptcha(in.VerifyCodeId, in.Vcode) {
		return nil, errors.New("verify code err")
	}

	tuser := user.UserInfo{UserName: in.UserName, Email: in.Email}
	err := utils.Gdb.First(tuser, tuser).Error
	if err == nil {
		return nil, errors.New("userName  or email already used")
	}
	uinfo := &user.UserInfo{
		UserName:     in.UserName,
		Email:        in.Email,
		PasswordHash: getPwdHash(in.Password),
	}
	err = utils.Gdb.Create(uinfo).Error
	if err != nil {
		return nil, err
	}

	token, err := uinfo.GenToken()
	if err != nil {
		return nil, errors.New("signup success, but login fail")
	}
	return &user.SignUpResponse{Token: token}, nil
}
func getPwdHash(pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pwd+"obd1234")))
}

func MustGetTokenUserId(tokenstr string) (int64, error) {
	token := user.Token{Token: tokenstr}
	err := utils.Gdb.First(token, token).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return 0, errors.New("token not found")
	} else if err != nil {
		return 0, err
	} else {
		if token.Disabled {
			return 0, errors.New("token Disabled")
		} else {
			return token.UserId, nil
		}
	}
	return 0, nil
}
func (us *UserServer) VerifyCode(tx context.Context, in *user.VerifyCodeRequest) (*user.VerifyCodeResponse, error) {
	id, err := utils.SendVerifyCodeMail(in.Email)
	if err != nil {
		return nil, err
	}
	return &user.VerifyCodeResponse{VerifyCodeId: id}, nil
}
func (us *UserServer) SignIn(tx context.Context, in *user.SignInRequest) (*user.SignInResponse, error) {
	tuser := user.UserInfo{}
	err := utils.Gdb.First(tuser, "user_name=? or email=?", in.UserName, in.UserName).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	} else {
		if tuser.PasswordHash != getPwdHash(in.Password) {
			return nil, errors.New("password error")
		}
	}
	token, err := tuser.GenToken()
	if err != nil {
		return nil, err
	}

	return &user.SignInResponse{Token: token, KycOk: tuser.KeyinfoOk}, nil
}

func (us *UserServer) GetUserInfo(ctx context.Context, in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	uInfo, err := getContextUser(ctx)
	return &user.GetUserInfoResponse{User: uInfo}, err
}

func getContextUser(ctx context.Context) (*user.UserInfo, error) {
	uid := utils.GetRpcCtxUserid(ctx)
	uInfo := new(user.UserInfo)
	err := utils.Gdb.First(uInfo, uid).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found:"+err.Error())
	}
	return uInfo, nil
}
func (us *UserServer) Kyc(ctx context.Context, in *user.KycRequest) (*user.KycResponse, error) {
	uInfo, err := getContextUser(ctx)
	if err != nil {
		return nil, err
	}
	if in.IdNum == "" {
		return nil, fmt.Errorf("id_card number is miss")
	}
	if in.Address1 == "" {
		return nil, fmt.Errorf("Address1  is miss")
	}
	if in.Address2 == "" {
		return nil, fmt.Errorf("Address2  is miss")
	}
	uInfo.CountryCode = in.CountryCode
	uInfo.Id1 = in.Id1
	uInfo.Id2 = in.Id2
	uInfo.SocialId = in.SocialId
	uInfo.Address1 = in.Address1
	uInfo.Address2 = in.Address2
	uInfo.Address3 = in.Address3
	err = utils.Gdb.Save(uInfo).Error
	if err != nil {
		return nil, err
	}
	res, err := us.sumbmitKyc(uInfo)
	if err != nil {
		return nil, err
	}
	uInfo.OpenId = res.OpenId
	uInfo.KeyinfoOk = true
	err = utils.Gdb.Save(uInfo).Error
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &user.KycResponse{}, nil
}
func (us *UserServer) ApplyCard(ctx context.Context, in *user.ApplyCardRequest) (*user.ApplyCardResponse, error) {
	uInfo, err := getContextUser(ctx)
	if err != nil {
		return nil, err
	}
	res, err := us.sumbmitApplyCard(uInfo, in.Currency)
	if err != nil {
		return nil, err
	}
	cardInfo := &user.CardInfo{
		CardNo:     res.CardNo,
		Symbol:     in.Currency,
		ExpiryDate: res.ExpiryDate,
		Cvv:        res.Cvv,
		UserId:     uInfo.Id,
	}
	err = utils.Gdb.Create(cardInfo).Error
	if err != nil {
		return nil, err
	}
	return &user.ApplyCardResponse{Card: cardInfo}, nil
}

type UploadFile struct {
	Id        int
	Filepath  string
	Size      int
	CreatedAt time.Time
	UserId    int64
	Tag       string
}

func (us *UserServer) Upload(ctx context.Context, in *user.UploadRequest) (*user.UploadResponse, error) {
	uid := utils.GetRpcCtxUserid(ctx)
	filePath := fmt.Sprintf("%d_%s.jpg", in.Tag, uid)
	//
	finfo := &UploadFile{
		Filepath: filePath,
		Size:     0,
		UserId:   uid,
		Tag:      in.Tag,
	}
	utils.Gdb.Create(finfo)
	return &user.UploadResponse{
		ImageUrl: filePath,
	}, nil
}

type tRes struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}
type tKycRes struct {
	Cid    string `json:"cid"`
	OpenId string `json:"openId"`
	//pend 待审 rejected 失败 passed 成功
	Status string `json:"status"`
}
type tKycObj struct {
	Address struct {
		AddressLine1 string `json:"addressLine1" example:"string"`
		AddressLine2 string `json:"addressLine2" example:"string"`
		City         string `json:"city" example:"string"`
		CountryCode  string `json:"countryCode" example:"string"`
		PostCode     string `json:"postCode" example:"string"`
		State        string `json:"state" example:"string"`
	} `json:"address"`
	Cid       string `json:"cid" example:"string"`
	Dob       string `json:"dob" example:"string"`
	Email     string `json:"email" example:"string"`
	FirstName string `json:"firstName" example:"string"`
	LastName  string `json:"lastName" example:"string"`
	MerNo     string `json:"merNo" example:"string"`
	Mobile    string `json:"mobile" example:"string"`
}

func (us *UserServer) sumbmitKyc(info *user.UserInfo) (res *tKycRes, err error) {
	obj := tKycObj{
		Address: struct {
			AddressLine1 string `json:"addressLine1" example:"string"`
			AddressLine2 string `json:"addressLine2" example:"string"`
			City         string `json:"city" example:"string"`
			CountryCode  string `json:"countryCode" example:"string"`
			PostCode     string `json:"postCode" example:"string"`
			State        string `json:"state" example:"string"`
		}{
			AddressLine1: info.Address1,
			AddressLine2: info.Address2,
			City:         "1234",
			CountryCode:  info.CountryCode.String(),
			PostCode:     "1234",
			State:        "1234",
		},
		Cid:       info.UserName,
		Dob:       "xx",
		Email:     info.Email,
		FirstName: "xx",
		LastName:  "xx",
		MerNo:     us.merNo,
		Mobile:    "13344445555",
	}
	resObj := new(tKycRes)
	err = utils.PSPostObj1(us.kycUrl+"/api/mg/kyc/v1", "", nil, obj, func() any {
		return resObj
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type tApplyCardObj struct {
	MerNo    string `json:"merNo"`
	OpenId   string `json:"openId"`
	IcNo     string `json:"icNo"`
	Currency string `json:"currency"`
}
type tApplyCardRes struct {
	CardNo     string `json:"cardNo"`
	ExpiryDate string `json:"expiryDate"`
	Cvv        string `json:"cvv"`
}

func (us *UserServer) sumbmitApplyCard(info *user.UserInfo, currencyCode user.CurrencyCode) (res *tApplyCardRes, err error) {
	obj := &tApplyCardObj{
		MerNo:    us.merNo,
		OpenId:   info.OpenId,
		IcNo:     info.IdNum,
		Currency: currencyCode.String(),
	}
	resObj := new(tApplyCardRes)
	err = utils.PSPostObj1(us.kycUrl+"/api/mg/bcard/add", "", nil, obj, func() any {
		return resObj
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type tCardDetailRes struct {
	CardNo         string `json:"cardNo"`
	ExpiryData     string `json:"expiryData"`
	Cvv            string `json:"cvv"`
	State          int    `json:"state"`
	Currency       string `json:"currency"`
	SettleCurrency string `json:"settleCurrency"`
	Amount         int    `json:"amount"`
}

func (us *UserServer) sumbmitCardDetail(cardNo string) (res *tCardDetailRes, err error) {
	obj := map[string]string{
		"merNo":  us.merNo,
		"cardNo": cardNo,
	}
	resObj := new(tCardDetailRes)
	err = utils.PSPostObj1(us.kycUrl+"/api/mg/bcard/detail", "", nil, obj, func() any {
		return resObj
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type tChargeRes struct {
	TradeNo      string      `json:"tradeNo"`
	MerOrderNo   string      `json:"merOrderNo"`
	PayCode      int         `json:"payCode"`
	PayUrl       interface{} `json:"payUrl"`
	PayDesc      string      `json:"payDesc"`
	SourceAmount string      `json:"sourceAmount"`
	CurrencyCode string      `json:"currencyCode"`
	TradeMerNo   int         `json:"tradeMerNo"`
	Sign         string      `json:"sign"`
}

func (us *UserServer) sumbmitCharge(cardNo string, amt float64) (res *tChargeRes, err error) {
	cardinfo := &user.CardInfo{CardNo: cardNo}
	err = utils.Gdb.First(cardinfo, cardinfo).Error
	if err != nil {
		return nil, err
	}
	currencyCode := cardinfo.Symbol.String()
	merOrderNo := shortuuid.New()
	//tradeMerNo、merOrderNo、currencyCode、sourceAmount、Merchant key
	signData := fmt.Sprintf("%v%v%v%.2f%v", us.merNo, merOrderNo, currencyCode, amt, us.merKey)
	sign := fmt.Sprintf("%X", md5.Sum([]byte(signData)))
	obj := map[string]any{
		"tradeMerNo":   us.merNo,
		"merOrderNo":   merOrderNo,
		"currencyCode": currencyCode,
		"sourceAmount": amt,
		//"returnUrl":         "http://baidu.com",
		"sign":              sign,
		"cardNo":            cardNo,
		"cardExpireMonth":   "06",
		"cardExpireYear":    "2020",
		"cardSecurityCode":  "361",
		"billingFirstName":  "三",
		"billingLastName":   "张",
		"billingAddress1":   "东方明珠",
		"billingCity":       "上海",
		"billingState":      "上海",
		"billingCountry":    "CN",
		"billingZipCode":    "200000",
		"billingPhone":      "13688888888",
		"shippingFirstName": "三",
		"shippingLastName":  "张",
		"shippingAddress1":  "东方明珠",
		"shippingCity":      "上海",
		"shippingState":     "上海",
		"shippingCountry":   "CN",
		"shippingZipCode":   "200000",
		"shippingPhone":     "13688888888",
		"userAgent":         "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/6.0)",
		"ipAddress":         "116.235.134.86",
		"browser": map[string]any{
			"acceptHeader":   "/**",
			"userAgent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
			"language":       "en",
			"timeZoneOffset": 10,
			"colorDepth":     10,
			"screenHeight":   100,
			"screenWidth":    1300,
			"javaEnabled":    true,
		},
		"holderName": "test",
		"productInfoList": []map[string]string{
			map[string]string{
				"sku":         "1",
				"productName": "mac pro",
				"price":       "12000",
				"quantity":    "1",
			},
		},
	}
	resObj := new(tCardDetailRes)
	err = utils.PSPostObj(us.payUrl+"/payment", "", nil, obj, func() any {
		return resObj
	})
	if err != nil {
		return nil, err
	}
	return res, nil
	return nil, err
}
