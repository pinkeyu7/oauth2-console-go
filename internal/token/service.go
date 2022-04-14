package token

import (
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
)

type Service interface {
	GenToken(req *apireq.GetSysAccountToken) (*apires.SysAccountToken, error)
}
