package service

import (
	"net/http"
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/oauth/client"
	"oauth2-console-go/internal/oauth/library"
	"oauth2-console-go/internal/oauth/scope"
	"oauth2-console-go/internal/system/sys_account"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/logr"
	"strings"

	"go.uber.org/zap"
)

type Service struct {
	sysAccRepo  sys_account.Repository
	clientRepo  client.Repository
	clientCache client.Cache
}

func NewService(sar sys_account.Repository, ocr client.Repository, occ client.Cache) client.Service {
	return &Service{
		sysAccRepo:  sar,
		clientRepo:  ocr,
		clientCache: occ,
	}
}

func (s *Service) ListClient(sysAccId int, page, perPage int) (*apires.ListOauthClient, error) {
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

	total, err := s.clientRepo.Count()
	if err != nil {
		unknownErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "count client error.", err)
		return nil, unknownErr
	}

	if page <= 1 {
		page = 1
	}

	if perPage <= 1 {
		perPage = 1
	}

	offset := (page - 1) * perPage

	list, err := s.clientRepo.Find(perPage, offset)
	if err != nil {
		unknownErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find client error.", err)
		return nil, unknownErr
	}

	res := apires.ListOauthClient{
		List:        list,
		Total:       total,
		CurrentPage: page,
		PerPage:     perPage,
	}

	return &res, nil
}

func (s *Service) GetClient(sysAccId int, clientId string, scopeRepo scope.Repository) (*apires.OauthClient, error) {
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

	// 取得 client app 資訊
	clt, err := s.clientRepo.FindOne(&model.OauthClient{Id: clientId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find client error.", err)
		return nil, findErr
	}
	if clt == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "client not found.", err)
		return nil, notFoundErr
	}

	// 取得所有 API 列表
	apis, err := scopeRepo.FindScope()
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return nil, findErr
	}

	// 從 API 列表建立授權清單
	scopeList, err := library.GenerateScopeList(apis)
	if err != nil {
		return nil, err
	}

	// 依據 client app 的 scope 設定授權清單
	scopeList, err = library.GenerateClientScopeList(scopeList, clt.Scope)
	if err != nil {
		return nil, err
	}

	res := apires.OauthClient{
		Id:           clt.Id,
		SysAccountId: clt.SysAccountId,
		Name:         clt.Name,
		Secret:       clt.Secret,
		Domain:       clt.Domain,
		Scope:        clt.Scope,
		IconPath:     clt.IconPath,
		ScopeList:    scopeList,
		CreatedAt:    clt.CreatedAt,
		UpdatedAt:    clt.UpdatedAt,
	}

	return &res, nil
}

func (s *Service) AddClient(req *apireq.AddOauthClientWithFile) error {
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

	// Check client id unique
	clt, err := s.clientRepo.FindOne(&model.OauthClient{Id: req.Id})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find client error.", err)
		return findErr
	}
	if clt != nil {
		duplicateErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "client id duplicate error.", err)
		return duplicateErr
	}

	// 上傳檔案
	// TODO - Upload image file

	// Insert client
	m := model.OauthClient{
		Id:           req.Id,
		SysAccountId: req.AccountId,
		Name:         req.Name,
		Secret:       req.Secret,
		Domain:       req.Domain,
		IconPath:     "",
	}

	err = s.clientRepo.Insert(&m)
	if err != nil {
		// 新增 Client App 失敗，刪除檔案
		// TODO - Delete upload image file
	}

	return nil
}

func (s *Service) EditClient(clientId string, req *apireq.EditOauthClientWithFile, scopeRepo scope.Repository) error {
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

	// Check client exist
	clt, err := s.clientRepo.FindOne(&model.OauthClient{Id: clientId})
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find client error.", err)
		return findErr
	}
	if clt == nil {
		notFoundErr := er.NewAppErr(http.StatusBadRequest, er.ResourceNotFoundError, "client not found.", err)
		return notFoundErr
	}

	// 從上傳的授權清單內找出授權項目
	scopes := library.GetScopesFromScopeList(req.ScopeList)

	// 取得所有 API 列表
	apis, err := scopeRepo.FindScope()
	if err != nil {
		findErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "find scope error.", err)
		return findErr
	}

	// 從 API 列表建立授權清單
	scopeList, err := library.GenerateScopeList(apis)
	if err != nil {
		return err
	}

	// 檢查授權
	validScopes, err := library.ValidateScopes(scopeList, scopes)
	if err != nil {
		return err
	}

	// 若是 has_image 為 true ，則檢查圖片
	// var fileName string
	// var uploadFilePath string

	if *req.HasImage {
		// 上傳檔案
		// TODO - Upload image file
	}

	// Update client
	m := model.OauthClient{
		Id:           clt.Id,
		SysAccountId: req.AccountId,
		Name:         req.Name,
		Secret:       req.Secret,
		Domain:       req.Domain,
		Scope:        strings.Join(validScopes, " "),
		IconPath:     clt.IconPath,
	}

	err = s.clientRepo.Update(&m)
	if err != nil {
		if *req.HasImage {
			// 編輯 Client App 失敗，刪除檔案
			// TODO - Delete upload image file
		}

		updateErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, "update client error.", err)
		return updateErr
	}

	// Delete cache
	err = s.clientCache.DeleteClientScopeList(clientId)
	if err != nil {
		logr.L.Error("delete oauth client scope list cache error.", zap.String("error", err.Error()))
	}

	return nil
}
