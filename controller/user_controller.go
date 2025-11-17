package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/constants"
	"github.com/mferdian/golang_boiller_plate/dto"
	"github.com/mferdian/golang_boiller_plate/logging"
	"github.com/mferdian/golang_boiller_plate/service"
	"github.com/mferdian/golang_boiller_plate/utils"
)

type (
	IUserController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)

		CreateUser(ctx *gin.Context)
		GetAllUser(ctx *gin.Context)
		GetUserByID(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
		DeleteUser(ctx *gin.Context)
	}

	UserController struct {
		userService service.IUserService
	}
)

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var payload dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.Register(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_REGISTER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_REGISTER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_REGISTER+": %s", result.Email)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_REGISTER, result)
	ctx.JSON(http.StatusCreated, res)
}

func (uc *UserController) Login(ctx *gin.Context) {
	var payload dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.Login(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_LOGIN_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_LOGIN_USER+": %s", payload.Email)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_LOGIN_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var payload dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.CreateUser(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_CREATE_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_USER+": %s", result.Email)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_USER, result)
	ctx.JSON(http.StatusCreated, res)
}

func (uc *UserController) GetAllUser(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	search := ctx.DefaultQuery("search", "")

	if !usePagination {
		// Tanpa pagination
		result, err := uc.userService.GetAllUser(ctx, search)
		if err != nil {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_LIST_USER, result)
		ctx.AbortWithStatusJSON(http.StatusOK, res)
		return
	}

	var query dto.UserPaginationRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.GetAllUserWithPagination(ctx.Request.Context(), query)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_LIST_USER+": page %d", query.Page)
	res := utils.Response{
		Status:   true,
		Messsage: constants.MESSAGE_SUCCESS_GET_LIST_USER,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID := ctx.GetString("id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idStr {
		logging.Log.Warn("unauthorized access: user trying to access another user's data")
		res := utils.BuildResponseFailed("unauthorized", "you can only get your own account", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	if _, err := uuid.Parse(idStr); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.GetuserByID(ctx.Request.Context(), idStr)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_DETAIL_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_USER+": %s", idStr)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userID := ctx.GetString("id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idParam {
		logging.Log.Warn("unauthorized update attempt by user")
		res := utils.BuildResponseFailed("unauthorized", "you can only update your own account", nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	var payload dto.UpdateUserRequest
	payload.ID = idParam

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.UpdateUser(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_USER+": %s", result.ID)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if _, err := uuid.Parse(idParam); err != nil {
		logging.Log.WithError(err).Warn(constants.MESSAGE_FAILED_UUID_FORMAT)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userID := ctx.GetString("id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idParam {
		logging.Log.Warn("unauthorized delete attempt by user")
		res := utils.BuildResponseFailed("unauthorized", "you can only delete your own account", nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	payload := dto.DeleteUserRequest{UserID: idParam}

	result, err := uc.userService.DeleteUser(ctx.Request.Context(), payload)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", idParam)
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_USER, result)
	ctx.JSON(http.StatusOK, res)
}
