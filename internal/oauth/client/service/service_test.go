package service

import (
	"encoding/json"
	"oauth2-console-go/config"
	"oauth2-console-go/driver"
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

func TestService_ListClient(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	ocr := clientRepo.NewRepository(orm)
	ocs := NewService(sar, ocr, occ)

	// Act
	res, err := ocs.ListClient(1, 1, 10)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, res.Total)
	assert.Len(t, res.List, 2)
}

func TestService_GetClient(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	osr := scopeRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	ocr := clientRepo.NewRepository(orm)
	ocs := NewService(sar, ocr, occ)

	// Act
	res, err := ocs.GetClient(1, "address-book-go", osr)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestService_EditClient(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	re, _ := driver.NewRedis()

	sar := sysAccRepo.NewRepository(orm)
	osr := scopeRepo.NewRepository(orm)
	occ := clientRepo.NewCache(re)
	ocr := clientRepo.NewRepository(orm)
	ocs := NewService(sar, ocr, occ)

	clientId := "address-book-go"
	hasImage := false

	req := apireq.EditOauthClient{
		AccountId: 1,
		Secret:    "12345678",
		Domain:    "http://abc.test.com/",
		Name:      "Test update client name",
		HasImage:  &hasImage,
		ScopeList: "{\"lifestyle\":{\"name\":\"lifestyle\",\"items\":{\"article_get\":{\"name\":\"article_get\",\"is_auth\":true},\"article_post\":{\"name\":\"article_post\",\"is_auth\":true},\"list_get\":{\"name\":\"list_get\",\"is_auth\":true}},\"is_auth\":false},\"user\":{\"name\":\"user\",\"items\":{\"profile_get\":{\"name\":\"profile_get\",\"is_auth\":true}},\"is_auth\":false}}\"",
	}

	scopeList := model.ScopeList{}
	_ = json.Unmarshal([]byte(req.ScopeList), &scopeList)

	request := apireq.EditOauthClientWithFile{
		EditOauthClient: &req,
		File:            nil,
		FileName:        "",
		FileExtension:   "",
		ScopeList:       &scopeList,
	}

	client := model.OauthClient{
		Id: clientId,
	}

	_, _ = orm.Get(&client)

	// Act
	err := ocs.EditClient(clientId, &request, osr)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(client.Id).Update(&client)
}
