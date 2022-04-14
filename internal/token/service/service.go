package service

import (
	"net/http"
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/system/sys_account"
	"oauth2-console-go/internal/token"
	tokenLibrary "oauth2-console-go/internal/token/library"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/helper"
	"time"
)

type Service struct {
	sysAccRepo sys_account.Repository
	tokenCache token.Cache
}

func NewService(sar sys_account.Repository, tc token.Cache) token.Service {
	return &Service{
		sysAccRepo: sar,
		tokenCache: tc,
	}
}

func (s *Service) GenToken(req *apireq.GetSysAccountToken) (*apires.SysAccountToken, error) {
	// Check Account Exist
	acc, err := s.sysAccRepo.FindOne(&model.SysAccount{Account: req.Account})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return nil, findErr
	}
	if acc == nil || acc.IsDisable {
		authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", nil)
		return nil, authErr
	}

	// Password not matched
	pw := helper.ScryptStr(req.Password)
	if acc.Password != pw {
		authErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", nil)
		return nil, authErr
	}

	oToken, expiredAt, err := tokenLibrary.GenToken(acc.Id)
	if err != nil {
		tokenErr := er.NewAppErr(http.StatusUnauthorized, er.UnauthorizedError, "", err)
		return nil, tokenErr
	}

	// Set iat
	iat := time.Now().UTC().Unix()
	_ = s.tokenCache.SetTokenIat(acc.Id, iat)

	mapData := map[string]interface{}{}
	mapData["name"] = acc.Name
	mapData["email"] = acc.Email

	res := apires.SysAccountToken{
		Token:     oToken,
		ExpiredAt: expiredAt,
		Data:      mapData,
	}

	return &res, nil
}
