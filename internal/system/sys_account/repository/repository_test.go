package repository

import (
	"fmt"
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

func TestRepository_Insert(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	m := model.SysAccount{
		Account:  "test_account",
		Phone:    "test_phone",
		Email:    "test_email",
		Password: "test_password",
		Name:     "test_name",
	}

	// Act
	err := sar.Insert(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(m.Id).Delete(&model.SysAccount{})
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	// Act
	testCases := []struct {
		Limit     int
		Offset    int
		WantCount int
	}{
		{
			2,
			0,
			1,
		},
		{
			10,
			10,
			0,
		},
	}
	// Act
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Find Sys Account,Offset:%d,Limit:%d", tc.Offset, tc.Limit), func(t *testing.T) {
			data, err := sar.Find(tc.Offset, tc.Limit)
			assert.Nil(t, err)
			assert.Len(t, data, tc.WantCount)
		})
	}
}

func TestRepository_FindOne(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	// No data
	// Act
	res, err := sar.FindOne(&model.SysAccount{Id: 100})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, res)

	// Has data
	// Act
	res, err = sar.FindOne(&model.SysAccount{Id: 1})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, res.Id)
}

func TestRepository_Count(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	// Act
	count, err := sar.Count()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	orm, _ := driver.NewXorm()
	sar := NewRepository(orm)

	acc := model.SysAccount{Id: 1}
	_, _ = orm.Get(&acc)

	m := model.SysAccount{
		Account:  "test_account",
		Phone:    "test_phone",
		Email:    "test_email",
		Password: "test_password",
		Name:     "test_name",
	}

	// Act
	err := sar.Update(&m)

	// Assert
	assert.Nil(t, err)

	// Teardown
	_, _ = orm.ID(acc.Id).Update(&acc)
}
