package http

import (
	"fmt"
	v1 "github.com/cossim/coss-server/internal/user/api/http/v1"
	"github.com/cossim/coss-server/internal/user/app/command"
	"github.com/cossim/coss-server/internal/user/app/query"
	"github.com/cossim/coss-server/internal/user/domain/entity"
	"github.com/cossim/coss-server/pkg/constants"
	"github.com/cossim/coss-server/pkg/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchUser searches for a user by email.
// @Summary 搜索用户
// @Description 搜索用户
// @Tags user
// @Security BearerAuth
// @Param email query string true "用户邮箱"
// @Success 200 {object} v1.Response{} "搜索用户成功"
// @Router /api/v1/user/search [get]
func (h *HttpServer) SearchUser(c *gin.Context, params v1.SearchUserParams) {
	h.logger.Info("Search user", zap.String("email", params.Email))
	searchUser, err := h.app.Queries.GetUser.Handle(c, &query.GetUse{
		CurrentUser: c.Value(constants.UserID).(string),
		//TargetUser:  params.Email,
		TargetEmail: params.Email,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "User found", getUserToResponse(searchUser))
}

// UpdateUserAvatar updates the user's avatar.
// @Summary 修改用户头像
// @Description 修改用户头像
// @Tags user
// @Security BearerAuth
// @Accept multipart/form-data
// @Param file formData string true "头像文件"
// @Success 200 {object} v1.Response{} "修改用户头像成功"
// @Router /api/v1/user/avatar [put]
func (h *HttpServer) UpdateUserAvatar(c *gin.Context) {
	// Parse form data
	if err := c.Request.ParseMultipartForm(25 << 20); // 25 MB limit
	err != nil {
		response.SetFail(c, "Failed to parse form data", nil)
		return
	}

	// Get the file from the form data
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		response.SetFail(c, "Error retrieving the file", nil)
		return
	}
	defer file.Close()

	// Check file type
	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		response.SetFail(c, "Unsupported file type. Only JPEG and PNG are allowed.", nil)
		return
	}

	// Check file size
	if handler.Size > 25<<20 { // 25 MB limit
		response.SetFail(c, "File size exceeds the limit. Maximum allowed size is 25 MB.", nil)
		return
	}

	url, err := h.app.Commands.UpdateUserAvatarHandler.Handle(c, &command.UpdateUserAvatar{
		UserID: c.Value(constants.UserID).(string),
		Avatar: file,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "更新用户头像成功", gin.H{"avatar": url})
}

// GetUser retrieves user information by ID.
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags user
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Success 200 {object} v1.Response{} "获取用户信息成功"
// @Router /api/v1/user/{id} [get]
func (h *HttpServer) GetUser(c *gin.Context, id string) {
	getUser, err := h.app.Queries.GetUser.Handle(c, &query.GetUse{
		CurrentUser: c.Value(constants.UserID).(string),
		TargetUser:  id,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "获取用户信息成功", getUserToResponse(getUser))
}

func getUserToResponse(e *entity.UserInfo) *v1.UserInfo {
	var preferences *v1.Preferences
	if e.Preferences != nil {
		preferences = &v1.Preferences{
			OpenBurnAfterReading:        e.Preferences.OpenBurnAfterReading,
			OpenBurnAfterReadingTimeOut: int(e.Preferences.OpenBurnAfterReadingTimeOut),
			Remark:                      e.Preferences.Remark,
			Silent:                      e.Preferences.SilentNotification,
		}
	}
	return &v1.UserInfo{
		Avatar:         e.Avatar,
		CossId:         e.CossID,
		Email:          e.Email,
		LastLoginTime:  e.LastLoginTime,
		LoginNumber:    int64(e.LoginNumber),
		NewDeviceLogin: e.NewDeviceLogin,
		Nickname:       e.Nickname,
		RelationStatus: e.RelationStatus.Int(),
		Signature:      e.Signature,
		Status:         v1.UserInfoStatus(e.Status),
		Tel:            e.Tel,
		UserId:         e.UserID,
		Preferences:    preferences,
	}
}

// UpdateUser updates user information.
// @Summary 修改用户信息
// @Description 修改用户信息
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Param body v1.UpdateUserJSONRequestBody true "用户信息"
// @Success 200 {object} v1.Response{} "修改用户信息成功"
// @Router /api/v1/user [put]
func (h *HttpServer) UpdateUser(c *gin.Context) {
	var req v1.UpdateUserJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("update user password", zap.Error(err))
		return
	}

	_, err := h.app.Commands.UpdateUser.Handle(c, &command.UpdateUser{
		UserID:    c.Value(constants.UserID).(string),
		Avatar:    req.Avatar,
		CossID:    req.CossId,
		Nickname:  req.Nickname,
		Signature: req.Signature,
		Tel:       req.Tel,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "更新用户信息成功", nil)
}

// UpdateUserPassword updates the user's password.
// @Summary 修改密码
// @Description 修改密码
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Param body v1.UpdateUserPasswordJSONRequestBody true "密码信息"
// @Success 200 {object} v1.Response{} "修改密码成功"
// @Router /api/v1/user/password [put]
func (h *HttpServer) UpdateUserPassword(c *gin.Context) {
	req := &v1.UpdateUserPasswordJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}
	_, err := h.app.Commands.UpdatePassword.Handle(c, &command.UpdatePassword{
		UserID:          c.Value(constants.UserID).(string),
		OldPassword:     req.OldPassword,
		ConfirmPassword: req.ConfirmPassword,
		NewPassword:     req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "更新用户密码成功", nil)
}

// UserActivate activates a user account.
// @Summary 用户激活
// @Description 用户激活
// @Tags user
// @Param user_id query string true "用户ID"
// @Param key query string true "激活密钥"
// @Success 200 {object} v1.Response{} "激活成功"
// @Router /api/v1/user/activate [get]
func (h *HttpServer) UserActivate(c *gin.Context, params v1.UserActivateParams) {
	_, err := h.app.Commands.UserActivate.Handle(c, &command.UserActivate{
		UserID:           params.UserId,
		VerificationCode: params.Key,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "用户激活成功", nil)
}

// GetUserBundle
// @Summary 获取用户信息包
// @Description 获取包含登录信息和用户详细信息的用户信息包
// @Tags user
// @Security BearerAuth
// @Success 200 {object} v1.Response{} "获取用户信息包成功"
// @Router /api/v1/user/bundle [get]
func (h *HttpServer) GetUserBundle(c *gin.Context, id string) {
	getUserBundle, err := h.app.Queries.GetUserBundle.Handle(c, &query.GetUserBundle{
		UserID: id,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "获取用户密钥包成功", &v1.UserSecretBundle{
		UserId:       getUserBundle.UserID,
		SecretBundle: getUserBundle.SecretBundle,
	})
}

func (h *HttpServer) UpdateUserBundle(c *gin.Context) {
	req := &v1.UpdateUserBundleJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	err := h.app.Commands.UpdateUserBundle.Handle(c, &command.UpdateUserBundle{
		UserID:       c.Value(constants.UserID).(string),
		SecretBundle: req.SecretBundle,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "更新用户密钥包成功", nil)
}

func (h *HttpServer) GetUserLoginClients(c *gin.Context) {
	getUserClients, err := h.app.Queries.GetUserLoginClients.Handle(c, &query.GetUserLoginClients{
		UserID: c.Value(constants.UserID).(string),
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "获取用户客户端成功", getUserClientsToResponse(getUserClients))
}

func getUserClientsToResponse(e []*query.GetUserLoginClientsResponse) []*v1.UserLoginClient {
	var userClients []*v1.UserLoginClient
	for _, v := range e {
		userClients = append(userClients, &v1.UserLoginClient{
			ClientIp:   v.ClientIP,
			DriverId:   v.DriverID,
			DriverType: v.DriverType,
			LoginAt:    v.LoginAt,
		})
	}
	return userClients
}

// UserEmailVerification sends an activation email.
// @Summary 发送激活邮件
// @Description 发送激活邮件
// @Tags user
// @Accept application/json
// @Param body v1.UserEmailVerificationJSONRequestBody{} true "发送激活邮件请求参数"
// @Success 200 {object} v1.Response{} "发送激活邮件成功"
// @Router /api/v1/user/email/verification [post]
func (h *HttpServer) UserEmailVerification(c *gin.Context) {
	req := &v1.UserEmailVerificationJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	err := h.app.Commands.SendUserEmailVerification.Handle(c, &command.SendUserEmailVerification{
		//ID: c.Value(constants.ID).(string),
		Email: string(req.Email),
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "发送邮箱验证码成功", nil)
}

// UserLogin authenticates the user and generates an access token.
// @Summary 用户登录
// @Description 用户登录认证并生成访问令牌。
// @Tags user
// @Accept application/json
// @Param body v1.LoginRequest true "登录请求参数"
// @Success 200 {object} v1.Response{} "登录成功"
// @Router /api/v1/user/login [post]
func (h *HttpServer) UserLogin(c *gin.Context) {
	req := &v1.UserLoginJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	userLogin, err := h.app.Commands.UserLogin.Handle(c, &command.UserLogin{
		Email:       req.Email,
		Password:    req.Password,
		DriverID:    req.DriverId,
		ClientIP:    c.ClientIP(),
		DriverToken: req.DriverToken,
		Platform:    req.Platform,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("user_id", userLogin.UserID)
	response.SetSuccess(c, "登录成功", ConversionUserLogin(userLogin))
}

func ConversionUserLogin(userLogin *command.UserLoginResponse) *v1.LoginResponse {
	return &v1.LoginResponse{
		Token: userLogin.Token,
		UserInfo: &struct {
			LastLoginTime  int64  `json:"last_login_time"`
			NewDeviceLogin bool   `json:"new_device_login"`
			UserId         string `json:"user_id"`
		}{LastLoginTime: userLogin.LastLoginTime, NewDeviceLogin: userLogin.NewDeviceLogin, UserId: userLogin.UserID},
	}
}

// UserLogout logs out the user.
// @Summary 退出登录
// @Description 退出登录
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Param body v1.UserLogoutJSONRequestBody true "退出登录请求参数"
// @Success 200 {object} v1.Response{} "退出登录成功"
// @Router /api/v1/user/logout [post]
func (h *HttpServer) UserLogout(c *gin.Context) {
	req := &v1.UserLogoutJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	uid := c.Value(constants.UserID).(string)
	err := h.app.Commands.UserLogout.Handle(c, &command.UserLogout{
		UserID:   uid,
		DriverID: req.DriverId,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "退出登录成功", nil)
}

// SetUserPublicKey sets the user's PGP public key.
// @Summary 设置用户pgp公钥
// @Description 设置用户pgp公钥
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Param body v1.SetUserPublicKeyJSONRequestBody true "设置pgp公钥请求参数"
// @Success 200 {object} v1.Response{} "设置成功"
// @Router /api/v1/user/public_key [post]
func (h *HttpServer) SetUserPublicKey(c *gin.Context) {
	req := &v1.SetUserPublicKeyJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	err := h.app.Commands.SetUserPublicKey.Handle(c, &command.SetUserPublicKey{
		UserID:    c.Value(constants.UserID).(string),
		PublicKey: req.PublicKey,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "设置用户公钥成功", nil)
}

func (h *HttpServer) ResetUserPublicKey(c *gin.Context) {
	req := &v1.ResetUserPublicKeyJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	err := h.app.Commands.ResetUserPublicKey.Handle(c, &command.ResetUserPublicKey{
		UserID:    c.Value(constants.UserID).(string),
		PublicKey: req.PublicKey,
		Code:      req.Code,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "重置用户公钥成功", nil)
}

// UserRegister registers a new user.
// @Summary 用户注册
// @Description 用户注册
// @Tags user
// @Accept application/json
// @Param body  v1.UserRegisterJSONRequestBody true "注册请求参数"
// @Success 200 {object} v1.Response{} "注册成功"
// @Router /api/v1/user/register [post]
func (h *HttpServer) UserRegister(c *gin.Context) {
	req := &v1.UserRegisterJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	userID, err := h.app.Commands.UserRegister.Handle(c, &command.UserRegister{
		Email:       string(req.Email),
		Password:    req.Password,
		ConfirmPass: req.ConfirmPassword,
		Nickname:    req.Nickname,
		PublicKey:   req.PublicKey,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("user_id", userID)
	response.SetSuccess(c, "注册成功", gin.H{"user_id": userID})
}

// GetPGPPublicKey retrieves the user's PGP public key.
// @Summary 获取用户pgp公钥
// @Description 获取用户pgp公钥
// @Tags user
// @Security BearerAuth
// @Success 200 {object} v1.Response{} "获取成功"
// @Router /api/v1/user/public_key [get]
func (h *HttpServer) GetPGPPublicKey(c *gin.Context) {
	response.SetSuccess(c, "获取系统pgp公钥成功", gin.H{"public_key": h.pgpKey})
}

func (h *HttpServer) ConfirmLogin(c *gin.Context, token string) {
	thisID := c.Value(constants.UserID).(string)
	handle, err := h.app.Queries.GetQRCode.Handle(c, &query.GetQrCode{Token: token})
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		return
	}

	if handle.Status == entity.QRCodeStatusExpired {
		response.SetFail(c, "二维码已过期", nil)
		return
	}

	if handle.UserID != thisID || handle.Status != entity.QRCodeStatusScanned {
		response.SetFail(c, "确认登录失败", nil)
		return
	}

	err = h.app.Commands.UpdateQRCode.Handle(c, &entity.QRCode{
		Token:  token,
		Status: entity.QRCodeStatusConfirmed,
		UserID: thisID,
	})

	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "登录确认成功", nil)
}

func (h *HttpServer) GenerateQRCode(c *gin.Context) {
	handle, err := h.app.Commands.GenerateQRCode.Handle(c, &command.GenerateQRCode{})
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		return
	}

	response.SetSuccess(c, "生成二维码成功", v1.GenerateQRCodeResponse{
		QrCode: &handle.QrCode,
		Token:  &handle.Token,
	})
}

func (h *HttpServer) SsoLogin(c *gin.Context, token string) {
	req := &v1.SSOLoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error("参数验证失败", zap.Error(err))
		response.SetFail(c, "参数验证失败", nil)
		return
	}

	handle, err := h.app.Queries.GetQRCode.Handle(c, &query.GetQrCode{Token: token})
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		return
	}

	if handle.Status == entity.QRCodeStatusExpired {
		response.SetFail(c, "二维码已过期", nil)
		return
	}

	if handle.Status != entity.QRCodeStatusConfirmed || handle.UserID == "" {
		response.SetFail(c, "登录失败", nil)
		return
	}

	handle2, err := h.app.Commands.SSOLogin.Handle(c, &command.SSOLogin{
		UserID:      handle.UserID,
		DriverID:    req.DriverId,
		ClientIP:    c.ClientIP(),
		Platform:    req.Platform,
		DriverToken: req.DriverToken,
	})
	if err != nil {
		c.Error(err)
		return
	}

	err = h.app.Commands.UpdateQRCode.Handle(c, &entity.QRCode{
		Token:  token,
		Status: entity.QRCodeStatusExpired,
		UserID: handle2.UserID,
	})

	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "登录成功", ConversionSSOLogin(handle2))
}

func ConversionSSOLogin(userLogin *command.SSOLoginResponse) *v1.LoginResponse {
	return &v1.LoginResponse{
		Token: userLogin.Token,
		UserInfo: &struct {
			LastLoginTime  int64  `json:"last_login_time"`
			NewDeviceLogin bool   `json:"new_device_login"`
			UserId         string `json:"user_id"`
		}{LastLoginTime: userLogin.LastLoginTime, NewDeviceLogin: userLogin.NewDeviceLogin, UserId: userLogin.UserID},
	}
}

func (h *HttpServer) VerifyQRCodeStatus(c *gin.Context, token string) {
	handle, err := h.app.Queries.GetQRCode.Handle(c, &query.GetQrCode{Token: token})
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		return
	}

	response.SetSuccess(c, "获取二维码状态成功", v1.QRCodeStatusResponse{
		Status: v1.QRCodeStatusResponseStatus(handle.Status),
	})
}

func (h *HttpServer) ScanQRCode(c *gin.Context, token string) {
	handle, err := h.app.Queries.GetQRCode.Handle(c, &query.GetQrCode{Token: token})
	if err != nil {
		fmt.Println(err)
		c.Error(err)
		return
	}

	if handle.Status == entity.QRCodeStatusScanned {
		response.SetFail(c, "重复扫码", nil)
		return
	}

	if handle.Status == entity.QRCodeStatusExpired {
		response.SetFail(c, "二维码已过期", nil)
		return
	}

	if handle.Status != entity.QRCodeStatusNotScanned {
		response.SetFail(c, "扫码失败", nil)
		return
	}

	err = h.app.Commands.UpdateQRCode.Handle(c, &entity.QRCode{
		Token:  token,
		Status: entity.QRCodeStatusScanned,
		UserID: c.Value(constants.UserID).(string),
	})

	if err != nil {
		c.Error(err)
		return
	}

	response.SetSuccess(c, "扫码成功", nil)
}
