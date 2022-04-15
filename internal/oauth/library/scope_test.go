package library

import (
	"fmt"
	"oauth2-console-go/config"
	"oauth2-console-go/pkg/valider"
	"os"
	"strings"
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

func TestGenerateScopeList(t *testing.T) {
	// Test error scope list
	scopes := []string{
		"user",
		"user.profile.put",
		"lifestyle.list.get",
		"lifestyle.article.get",
		"lifestyle.article.post",
	}

	// Act
	_, err := GenerateScopeList(scopes)

	// Assert
	assert.NotNil(t, err)

	// Test normal scope list
	scopes = []string{
		"user.profile_get",
		"user.profile_put",
		"lifestyle.list_get",
		"lifestyle.article_get",
		"lifestyle.article_post",
		"friendship.list_get",
	}

	// Act
	scopeList, err := GenerateScopeList(scopes)

	// Assert
	assert.Nil(t, err)
	for _, scope := range scopes {
		t.Run(fmt.Sprintf("Find scope:%s", scope), func(t *testing.T) {
			arr := strings.Split(scope, ".")
			assert.NotNil(t, (*scopeList)[arr[0]])
			assert.NotNil(t, (*scopeList)[arr[0]].Items[arr[1]])
		})
	}
}

func TestParseScope(t *testing.T) {
	// Act
	testCases := []struct {
		scope    string
		level    int
		category string
		item     string
		isError  bool
	}{
		{
			"",
			0,
			"",
			"",
			true,
		},
		{
			"user",
			1,
			"user",
			"",
			false,
		},
		{
			"user.profile",
			2,
			"user",
			"profile",
			false,
		},
		{
			"user.profile.get",
			0,
			"",
			"",
			true,
		},
	}

	// Assert
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Parse scope:%s", tc.scope), func(t *testing.T) {
			level, category, item, err := ParseScope(tc.scope)
			if tc.isError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.level, level)
			assert.Equal(t, tc.category, category)
			assert.Equal(t, tc.item, item)
		})
	}
}
