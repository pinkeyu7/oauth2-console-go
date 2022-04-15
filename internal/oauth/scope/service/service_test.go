package service

import (
	"oauth2-console-go/config"
	"oauth2-console-go/driver"
	_ "oauth2-console-go/driver"
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/model"
	clientRepo "oauth2-console-go/internal/oauth/client/repository"
	scopeRepo "oauth2-console-go/internal/oauth/scope/repository"
	sysAccRepo "oauth2-console-go/internal/system/sys_account/repository"
	"oauth2-console-go/pkg/valider"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
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

func TestService_ListScope(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	osc := scopeRepo.NewCache(re)
	osr := scopeRepo.NewRepository(orm)
	oss := NewService(sar, osr, osc, occ)

	// Act
	res, err := oss.ListScope(1, 1, 10)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 4, res.Total)
	assert.Len(t, res.List, 4)
}

func TestService_GetScope(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	osc := scopeRepo.NewCache(re)
	osr := scopeRepo.NewRepository(orm)
	oss := NewService(sar, osr, osc, occ)

	// Act
	res, err := oss.GetScope(1, 1)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestOauth2ScopeServiceAddScope(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	osc := scopeRepo.NewCache(re)
	osr := scopeRepo.NewRepository(orm)
	oss := NewService(sar, osr, osc, occ)

	req := apireq.AddOauthScope{
		AccountId:   1,
		Scope:       "user.profile_put",
		Path:        "/v1/users",
		Method:      "PUT",
		Name:        "編輯個人資訊",
		Description: "編輯個人資訊",
	}

	// Act
	err := oss.AddScope(&req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Where(" scope = ? ", req.Scope).Delete(&model.OauthScope{})
}

func TestOauth2ScopeServiceEditScope(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	osc := scopeRepo.NewCache(re)
	osr := scopeRepo.NewRepository(orm)
	oss := NewService(sar, osr, osc, occ)

	scopeId := 1
	isDisable := true
	req := apireq.EditOauthScope{
		AccountId:   1,
		Name:        "Test update name",
		Description: "Test update description",
		IsDisable:   &isDisable,
	}

	scope := model.OauthScope{
		Id: scopeId,
	}

	_, _ = orm.Get(&scope)

	// Act
	err := oss.EditScope(scopeId, &req)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(scope.Id).Update(&scope)
}
