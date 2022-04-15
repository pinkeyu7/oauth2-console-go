package library

import (
	"fmt"
	"oauth2-console-go/dto/model"
	"oauth2-console-go/pkg/er"
	"strings"
)

func ValidateScopes(scopeList *model.ScopeList, scopes []string) ([]string, error) {
	scpList := *scopeList
	validScopes := make([]string, 0)

	for _, scope := range scopes {
		// 檢查權限的欄位格式
		level, category, item, err := ParseScope(scope)
		if err != nil {
			return nil, err
		}

		switch level {
		case 1:
			// 檢查 Level 1 授權
			if scpList[category] == nil {
				notFoundErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope not found error.", nil)
				return nil, notFoundErr
			}
			if scpList[category] != nil {
				validScopes = append(validScopes, scope)
			}
		case 2:
			// 檢查 Level 2 授權
			if scpList[category] == nil {
				notFoundErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope not found error.", nil)
				return nil, notFoundErr
			}
			if scpList[category].Items[item] == nil {
				notFoundErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope not found error.", nil)
				return nil, notFoundErr
			}
			if scpList[category].Items[item] != nil {
				validScopes = append(validScopes, scope)
			}
		}
	}

	return validScopes, nil
}

func GetScopesFromScopeList(scopeList *model.ScopeList) []string {
	scpList := *scopeList
	scopes := make([]string, 0)
	for _, category := range scpList {
		if category.IsAuth {
			scope := fmt.Sprintf("%s", category.Name)
			scopes = append(scopes, scope)
			continue
		}
		for _, item := range category.Items {
			if item.IsAuth {
				scope := fmt.Sprintf("%s.%s", category.Name, item.Name)
				scopes = append(scopes, scope)
			}
		}
	}

	return scopes
}

func CheckScope(scopeList *model.ScopeList, scope string) (bool, error) {
	scpList := *scopeList
	// 檢查權限的欄位格式
	_, category, item, err := ParseScope(scope)
	if err != nil {
		return false, err
	}

	// 檢查 Level 1 授權
	if scpList[category] == nil {
		return false, nil
	}
	if scpList[category].IsAuth {
		return true, nil
	}

	// 檢查 Level 2 授權
	if scpList[category].Items[item] == nil {
		return false, nil
	}
	if scpList[category].Items[item].IsAuth {
		return true, nil
	}

	return false, nil
}

func GenerateClientScopeList(scopeList *model.ScopeList, clientScopeStr string) (*model.ScopeList, error) {
	scpList := *scopeList
	clientScopes := strings.Split(clientScopeStr, " ")
	for _, scope := range clientScopes {
		// 檢查權限的欄位格式
		level, category, item, err := ParseScope(scope)
		if err != nil {
			continue
		}

		switch level {
		case 1:
			// 設定 Level 1 授權
			if scpList[category] == nil {
				continue
			}
			if scpList[category] != nil {
				scpList[category].IsAuth = true
			}
		case 2:
			// 設定 Level 2 授權
			if scpList[category] == nil {
				continue
			}
			if scpList[category].Items[item] == nil {
				continue
			}
			if scpList[category].Items[item] != nil {
				scpList[category].Items[item].IsAuth = true
			}
		}
	}

	return scopeList, nil
}

func GenerateScopeList(scopes []string) (*model.ScopeList, error) {
	scopeList := make(model.ScopeList, 0)

	for _, scope := range scopes {
		// 檢查權限的欄位格式
		level, category, item, err := ParseScope(scope)
		if err != nil {
			return nil, err
		}
		if level != 2 {
			notMatchErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope level error.", nil)
			return nil, notMatchErr
		}

		if scopeList[category] == nil {
			scopeList[category] = &model.ScopeCategory{
				Name:   category,
				Items:  make(map[string]*model.ScopeItem, 0),
				IsAuth: false,
			}
		}

		if scopeList[category].Items[item] == nil {
			scopeList[category].Items[item] = &model.ScopeItem{
				Name:   item,
				IsAuth: false,
			}
		}
	}

	return &scopeList, nil
}

func ParseScope(scope string) (int, string, string, error) {
	var level int
	var category string
	var item string

	if len(scope) == 0 {
		emptyErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope empty error.", nil)
		return 0, "", "", emptyErr
	}

	arr := strings.Split(scope, ".")
	level = len(arr)

	switch len(arr) {
	case 1:
		category = arr[0]
	case 2:
		category = arr[0]
		item = arr[1]
	default:
		formatErr := er.NewAppErr(400, er.ErrorParamInvalid, "scope format error.", nil)
		return 0, "", "", formatErr
	}
	return level, category, item, nil
}
