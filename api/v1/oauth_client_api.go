package v1

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"oauth2-console-go/api"
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/model"
	clientRepo "oauth2-console-go/internal/oauth/client/repository"
	clientSrv "oauth2-console-go/internal/oauth/client/service"
	scopeRepo "oauth2-console-go/internal/oauth/scope/repository"
	sysAccRepo "oauth2-console-go/internal/system/sys_account/repository"
	tokenLibrary "oauth2-console-go/internal/token/library"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/helper"
	"oauth2-console-go/pkg/valider"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOauthClient
// @Summary List Oauth Client - Oauth Client App 列表
// @Produce json
// @Accept json
// @Tags Oauth Client
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id query int true "Account ID"
// @Param page query int true "Page"
// @Param per_page query int true "PerPage"
// @Success 200 {object} apires.ListOauthClient
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/clients [get]
func ListOauthClient(c *gin.Context) {
	req := apireq.ListOauthClient{}
	err := c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	sar := sysAccRepo.NewRepository(env.Orm)
	occ := clientRepo.NewCache(env.RedisCluster)
	ocr := clientRepo.NewRepository(env.Orm)
	ocs := clientSrv.NewService(sar, ocr, occ)

	res, err := ocs.ListClient(req.AccountId, req.Page, req.PerPage)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetOauthClient
// @Summary Get Oauth Client 取得 Client APP 資訊
// @Produce json
// @Accept json
// @Tags Oauth Client
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param client_id path string true "Oauth Client ID"
// @Param account_id query int true "Account ID"
// @Success 200 {object} apires.OauthClient
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/clients/{client_id} [get]
func GetOauthClient(c *gin.Context) {
	clientId := c.Param("id")

	accIdStr := c.Query("account_id")
	accId, err := strconv.Atoi(accIdStr)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "account id format error.", err)
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, accId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	sar := sysAccRepo.NewRepository(env.Orm)
	osr := scopeRepo.NewRepository(env.Orm)
	occ := clientRepo.NewCache(env.RedisCluster)
	ocr := clientRepo.NewRepository(env.Orm)
	ocs := clientSrv.NewService(sar, ocr, occ)

	res, err := ocs.GetClient(accId, clientId, osr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddOauthClient
// @Summary Add Oauth Client - 新增 Client APP
// @Produce json
// @Accept json
// @Tags Oauth Client
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id formData int true "Account id"
// @Param id formData string true "Client Id"
// @Param secret formData string true "Client Secret"
// @Param domain formData string true "Client Domain"
// @Param name formData string true "Client Name"
// @Param file formData file true "Client Icon Image"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/clients [post]
func AddOauthClient(c *gin.Context) {
	req := apireq.AddOauthClient{}
	err := c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// validate upload image
	file, fileName, fileExtension, err := helper.CheckFormUploadImage(c, "file", 2) // 2MB
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	sar := sysAccRepo.NewRepository(env.Orm)
	occ := clientRepo.NewCache(env.RedisCluster)
	ocr := clientRepo.NewRepository(env.Orm)
	ocs := clientSrv.NewService(sar, ocr, occ)

	request := apireq.AddOauthClientWithFile{
		AddOauthClient: &req,
		File:           file,
		FileName:       fileName,
		FileExtension:  fileExtension,
	}

	err = ocs.AddClient(&request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// EditOauthClient
// @Summary Edit Oauth Client - 編輯 Client APP
// @Produce json
// @Accept json
// @Tags Oauth Client
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param client_id path string true "Oauth Client ID"
// @Param account_id formData int true "Account id"
// @Param secret formData string true "Client Secret"
// @Param domain formData string true "Client Domain"
// @Param name formData string true "Client Name"
// @Param has_image formData bool true "Upload Image for Update"
// @Param file formData file true "Client Icon Image"
// @Param scope_list formData string true "Client Scope List(After json stringify)"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/clients/{client_id} [put]
func EditOauthClient(c *gin.Context) {
	clientId := c.Param("id")

	req := apireq.EditOauthClient{}
	err := c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 將 json stringify 轉回 ScopeList
	scopeList := model.ScopeList{}
	err = json.Unmarshal([]byte(req.ScopeList), &scopeList)
	if err != nil {
		parseErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "scope list json parse error.", err)
		_ = c.Error(parseErr)
		return
	}

	// 更新 client app 的資料時，icon是否更新由前端提供 has_image 判斷
	// 若是 has_image 為 true ，則檢查圖片
	var file multipart.File
	var fileName string
	var fileExtension string

	if *req.HasImage {
		// validate upload image
		file, fileName, fileExtension, err = helper.CheckFormUploadImage(c, "file", 2) // 2MB
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	// 驗證 jwt user == user_id
	err = tokenLibrary.CheckJWTAccountId(c, req.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	env := api.GetEnv()
	sar := sysAccRepo.NewRepository(env.Orm)
	osr := scopeRepo.NewRepository(env.Orm)
	occ := clientRepo.NewCache(env.RedisCluster)
	ocr := clientRepo.NewRepository(env.Orm)
	ocs := clientSrv.NewService(sar, ocr, occ)

	request := apireq.EditOauthClientWithFile{
		EditOauthClient: &req,
		File:            file,
		FileName:        fileName,
		FileExtension:   fileExtension,
		ScopeList:       &scopeList,
	}

	err = ocs.EditClient(clientId, &request, osr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
