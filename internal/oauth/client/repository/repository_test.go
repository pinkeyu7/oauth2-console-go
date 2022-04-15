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
	ocr := NewRepository(orm)

	// Act
	total, err := ocr.Count()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, total)
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	ocr := NewRepository(orm)

	// Act
	clients, err := ocr.Find(10, 0)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, clients)
	assert.Len(t, clients, 2)
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	ocr := NewRepository(orm)

	// Act
	client, err := ocr.FindOne(&model.OauthClient{
		Id: "address-book-go",
	})

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestRepository_Insert(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	ocr := NewRepository(orm)

	id := "test_client"
	secret := "pa@@w0rd"
	domain := "http://localhost:9088"
	info := model.OauthClient{
		Id:           id,
		SysAccountId: 1,
		Name:         "VendorNo9999",
		Secret:       secret,
		Domain:       domain,
		Scope:        "all",
		IconPath:     "image.svg",
	}

	err := ocr.Insert(&info)

	// Assert
	assert.Nil(t, err)

	// TearDown
	_, _ = orm.Where("id = ? ", info.Id).Delete(&model.OauthClient{})
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	ocr := NewRepository(orm)

	id := "address-book-go"
	client, _ := ocr.FindOne(&model.OauthClient{Id: id})

	// Act
	req := model.OauthClient{
		Id:           client.Id,
		SysAccountId: 1,
		Name:         "Test Update Client",
		Secret:       client.Secret,
		Domain:       client.Domain,
		Scope:        client.Scope,
		IconPath:     client.IconPath,
	}

	err := ocr.Update(&req)

	// Assert
	assert.Nil(t, err)

	// Put back
	_, _ = orm.Where("id = ?", client.Id).Update(client)
}
