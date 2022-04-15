package repository

import (
	"oauth2-console-go/config"
	"oauth2-console-go/driver"
	"oauth2-console-go/dto/model"
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

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	// Act
	total, err := osr.Count()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 4, total)
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	// Act
	scopes, err := osr.Find(10, 0)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, scopes)
	assert.Len(t, scopes, 4)
}

func TestRepository_FindScope(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	// Act
	scopes, err := osr.FindScope()

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, scopes)
	assert.Len(t, scopes, 4)
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	// Act
	scope, err := osr.FindOne(&model.OauthScope{
		Path:   "/v1/users",
		Method: "GET",
	})

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, scope)
}

func TestRepository_Insert(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	isDisable := false
	scope := model.OauthScope{
		Scope:       "user.profile_put",
		Path:        "/v1/users",
		Method:      "PUT",
		Name:        "編輯個人資訊",
		Description: "編輯個人資訊",
		IsDisable:   &isDisable,
	}

	// Act
	err := osr.Insert(&scope)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.Delete(&scope)
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	osr := NewRepository(orm)

	// Act
	scope := model.OauthScope{
		Id: 1,
	}

	_, _ = orm.Get(&scope)

	// Act
	err := osr.Update(&model.OauthScope{
		Name:        "Test update name",
		Description: "Test update description",
	})

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(scope.Id).Update(&scope)
}
