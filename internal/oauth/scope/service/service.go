package service

import (
	"net/http"
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/oauth/client"
	"oauth2-console-go/internal/oauth/scope"
	"oauth2-console-go/internal/system/sys_account"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/logr"

	"go.uber.org/zap"
)

type Service struct {
	sysAccRepo  sys_account.Repository
	scopeRepo   scope.Repository
	scopeCache  scope.Cache
	clientCache client.Cache
}

func NewService(sar sys_account.Repository, osr scope.Repository, osc scope.Cache, occ client.Cache) scope.Service {
	return &Service{
		sysAccRepo:  sar,
		scopeRepo:   osr,
		scopeCache:  osc,
		clientCache: occ,
	}
}

func (s *Service) ListScope(sysAccId, page, perPage int) (*apires.ListOauthScope, error) {
	// Check account id exist
	sysAcc := model.SysAccount{Id: sysAccId}
	acc, err := s.sysAccRepo.FindOne(&sysAcc)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return nil, findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "account not found.", err)
		return nil, notFoundErr
	}

	total, err := s.scopeRepo.Count()
	if err != nil {
		unknownErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count scope error.", err)
		return nil, unknownErr
	}

	if page <= 1 {
		page = 1
	}

	if perPage <= 1 {
		perPage = 1
	}

	offset := (page - 1) * perPage

	list, err := s.scopeRepo.Find(perPage, offset)
	if err != nil {
		unknownErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return nil, unknownErr
	}

	res := apires.ListOauthScope{
		List:        list,
		Total:       total,
		CurrentPage: page,
		PerPage:     perPage,
	}

	return &res, nil
}

func (s *Service) GetScope(sysAccId, scopeId int) (*model.OauthScope, error) {
	// Check account id exist
	sysAcc := model.SysAccount{Id: sysAccId}
	acc, err := s.sysAccRepo.FindOne(&sysAcc)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return nil, findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "account not found.", err)
		return nil, notFoundErr
	}

	scp, err := s.scopeRepo.FindOne(&model.OauthScope{Id: scopeId})
	if err != nil {
		unknownErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return nil, unknownErr
	}
	if scp == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "scope not found.", err)
		return nil, notFoundErr
	}

	return scp, nil
}

func (s *Service) AddScope(req *apireq.AddOauthScope) error {
	// Check account id exist
	sysAcc := model.SysAccount{Id: req.AccountId}
	acc, err := s.sysAccRepo.FindOne(&sysAcc)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "account not found.", err)
		return notFoundErr
	}

	// Check scope unique
	scp, err := s.scopeRepo.FindOne(&model.OauthScope{Scope: req.Scope})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return findErr
	}
	if scp != nil {
		duplicateErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "scope duplicate error.", err)
		return duplicateErr
	}

	// Insert scope
	isDisable := false
	m := model.OauthScope{
		Scope:       req.Scope,
		Path:        req.Path,
		Method:      req.Method,
		Name:        req.Name,
		Description: req.Description,
		IsDisable:   &isDisable,
	}

	err = s.scopeRepo.Insert(&m)
	if err != nil {
		insertErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "insert scope error.", err)
		return insertErr
	}

	// 清除所有 client app 的授權列表
	err = s.clientCache.DeleteAllClientScopeList()
	if err != nil {
		logr.L.Error("delete client scope list cache error.", zap.String("error", err.Error()))
	}

	return nil
}

func (s *Service) EditScope(scopeId int, req *apireq.EditOauthScope) error {
	// Check account id exist
	sysAcc := model.SysAccount{Id: req.AccountId}
	acc, err := s.sysAccRepo.FindOne(&sysAcc)
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find account error.", err)
		return findErr
	}
	if acc == nil || acc.IsDisable {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "account not found.", err)
		return notFoundErr
	}

	// Check scope exist
	scp, err := s.scopeRepo.FindOne(&model.OauthScope{Id: scopeId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return findErr
	}
	if scp == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "scope not found.", err)
		return notFoundErr
	}

	// Update scope
	scp.Name = req.Name
	scp.Description = req.Description
	scp.IsDisable = req.IsDisable

	err = s.scopeRepo.Update(scp)
	if err != nil {
		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update scope error.", err)
		return updateErr
	}

	// Delete cache
	err = s.scopeCache.DeleteOne(scp.Path, scp.Method)
	if err != nil {
		logr.L.Error("delete scope cache error.", zap.String("error", err.Error()))
	}

	// 清除所有 client app 的授權列表
	err = s.clientCache.DeleteAllClientScopeList()
	if err != nil {
		logr.L.Error("delete client scope list cache error.", zap.String("error", err.Error()))
	}

	return nil
}
