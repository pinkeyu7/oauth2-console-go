package v1

import (
	"net/http"
	"oauth2-console-go/api"
	"oauth2-console-go/dto/apireq"
	clientRepo "oauth2-console-go/internal/oauth/client/repository"
	scopeRepo "oauth2-console-go/internal/oauth/scope/repository"
	scopeSrv "oauth2-console-go/internal/oauth/scope/service"
	sysAccRepo "oauth2-console-go/internal/system/sys_account/repository"
	tokenLibrary "oauth2-console-go/internal/token/library"
	"oauth2-console-go/pkg/er"
	"oauth2-console-go/pkg/valider"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOauthScope
// @Summary List Oauth Scope - Oauth Scope API 列表
// @Produce json
// @Accept json
// @Tags Oauth Scope
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param account_id query int true "Account ID"
// @Param page query int true "Page"
// @Param per_page query int true "PerPage"
// @Success 200 {object} apires.ListOauthScope
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/scopes [get]
func ListOauthScope(c *gin.Context) {
	req := apireq.ListOauthScope{}
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
	osc := scopeRepo.NewCache(env.RedisCluster)
	osr := scopeRepo.NewRepository(env.Orm)
	oss := scopeSrv.NewService(sar, osr, osc, occ)

	res, err := oss.ListScope(req.AccountId, req.Page, req.PerPage)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetOauthScope
// @Summary Get Oauth Scope 取得 API 資訊
// @Produce json
// @Accept json
// @Tags Oauth Scope
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param scope_id path int true "Oauth Scope ID"
// @Param account_id query int true "Account ID"
// @Success 200 {object} model.OauthScope
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/scopes/{scope_id} [get]
func GetOauthScope(c *gin.Context) {
	id := c.Param("id")
	scopeId, err := strconv.Atoi(id)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "scope id format error.", err)
		_ = c.Error(err)
		return
	}

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
	occ := clientRepo.NewCache(env.RedisCluster)
	osc := scopeRepo.NewCache(env.RedisCluster)
	osr := scopeRepo.NewRepository(env.Orm)
	oss := scopeSrv.NewService(sar, osr, osc, occ)

	res, err := oss.GetScope(accId, scopeId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddOauthScope
// @Summary Add Oauth Scope - 新增 API
// @Produce json
// @Accept json
// @Tags Oauth Scope
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param Body body apireq.AddOauthScope true "Request Add Oauth Scope"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/scopes [post]
func AddOauthScope(c *gin.Context) {
	req := apireq.AddOauthScope{}
	err := c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
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
	osc := scopeRepo.NewCache(env.RedisCluster)
	osr := scopeRepo.NewRepository(env.Orm)
	oss := scopeSrv.NewService(sar, osr, osc, occ)

	err = oss.AddScope(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

// EditOauthScope
// @Summary Edit Oauth Scope - 編輯 API
// @Produce json
// @Accept json
// @Tags Oauth Scope
// @Security Bearer
// @Param Bearer header string true "JWT Token"
// @Param scope_id path string true "Oauth Scope ID"
// @Param Body body apireq.EditOauthScope true "Request Edit Oauth Scope"
// @Success 200 {string} string "{}"
// @Failure 400 {object} er.AppErrorMsg "{"code":"400400","message":"Wrong parameter format or invalid"}"
// @Failure 401 {object} er.AppErrorMsg "{"code":"400401","message":"Unauthorized"}"
// @Failure 403 {object} er.AppErrorMsg "{"code":"400403","message":"Permission denied"}"
// @Failure 404 {object} er.AppErrorMsg "{"code":"400404","message":"Resource not found"}"
// @Failure 500 {object} er.AppErrorMsg "{"code":"500000","message":"Database unknown error"}"
// @Router /v1/oauth/scopes/{scope_id} [put]
func EditOauthScope(c *gin.Context) {
	id := c.Param("id")
	scopeId, err := strconv.Atoi(id)
	if err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, "scope id format error.", err)
		_ = c.Error(paramErr)
		return
	}

	req := apireq.EditOauthScope{}
	err = c.Bind(&req)
	if err != nil {
		err = er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(err)
		return
	}

	// 參數驗證
	err = valider.Validate.Struct(req)
	if err != nil {
		paramErr := er.NewAppErr(400, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
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
	osc := scopeRepo.NewCache(env.RedisCluster)
	osr := scopeRepo.NewRepository(env.Orm)
	oss := scopeSrv.NewService(sar, osr, osc, occ)

	err = oss.EditScope(scopeId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
