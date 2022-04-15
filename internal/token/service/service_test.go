package service

import (
	"net/http"
	"oauth2-console-go/config"
	"oauth2-console-go/driver"
	"oauth2-console-go/dto/apireq"
	sysAccRepo "oauth2-console-go/internal/system/sys_account/repository"
	tokenRepo "oauth2-console-go/internal/token/repository"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/valider"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	config.InitEnv()
	valider.Init()
}

func TestService_GenToken(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	rc, _ := driver.NewRedis()
	sar := sysAccRepo.NewRepository(orm)
	tc := tokenRepo.NewRedis(rc)
	ts := NewService(sar, tc)

	// No data
	req := apireq.GetSysAccountToken{
		Account:  "sys_account",
		Password: "123456",
	}

	// Act
	res, err := ts.GenToken(&req)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, res)
	authErr := err.(*er.AppError)
	assert.Equal(t, http.StatusUnauthorized, authErr.StatusCode)
	assert.Equal(t, strconv.Itoa(er.UnauthorizedError), authErr.Code)

	// Has data
	req = apireq.GetSysAccountToken{
		Account:  "sys_account",
		Password: "A12345678",
	}

	// Act
	res, err = ts.GenToken(&req)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
