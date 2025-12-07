package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/constants"
	"github.com/mferdian/golang_boiller_plate/dto"
	"github.com/mferdian/golang_boiller_plate/helpers"
	"github.com/mferdian/golang_boiller_plate/logging"
	"github.com/mferdian/golang_boiller_plate/model"
	"github.com/mferdian/golang_boiller_plate/repository"
)

type (
	IUserService interface {
		Register(ctx context.Context, req dto.RegisterUserRequest) (dto.RegisterUserResponse, error)
		Login(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error)

		CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
		GetuserByID(ctx context.Context, userID string) (dto.UserResponse, error)
		GetAllUser(ctx context.Context, search string) ([]dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error)
		UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error)
		DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.UserResponse, error)
	}

	UserService struct {
		userRepo   repository.IUserRepository
		jwtService InterfaceJWTService
	}
)

func NewUserService(userRepo repository.IUserRepository, jwtService InterfaceJWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (us *UserService) Register(ctx context.Context, req dto.RegisterUserRequest) (dto.RegisterUserResponse, error) {
	if len(req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": name too short")
		return dto.RegisterUserResponse{}, constants.ErrInvalidName
	}

	if !helpers.IsValidEmail(req.Email) {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": invalid email format")
		return dto.RegisterUserResponse{}, constants.ErrInvalidEmail
	}

	_, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err == nil && found {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": email already exists")
		return dto.RegisterUserResponse{}, constants.ErrEmailAlreadyExists
	}

	if len(req.Password) < 8 {
		logging.Log.Warn(constants.MESSAGE_FAILED_REGISTER + ": password too short")
		return dto.RegisterUserResponse{}, constants.ErrInvalidPassword
	}

	user := model.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     constants.ENUM_ROLE_USER,
	}

	err = us.userRepo.Register(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_REGISTER)
		return dto.RegisterUserResponse{}, constants.ErrRegisterUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_REGISTER+": %s", user.Email)

	return dto.RegisterUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (us *UserService) Login(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error) {
	user, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err != nil || !found {
		logging.Log.Warn(constants.MESSAGE_FAILED_LOGIN_USER + ": email not found")
		return dto.LoginResponse{}, constants.ErrInvalidLoginCredential
	}

	if ok, err := helpers.CheckPassword(user.Password, []byte(req.Password)); !ok || err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_LOGIN_USER + ": password mismatch")
		return dto.LoginResponse{}, constants.ErrInvalidLoginCredential
	}

	accessToken, refreshToken, err := us.jwtService.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_LOGIN_USER + ": failed generate token")
		return dto.LoginResponse{}, constants.ErrGenerateAccessToken
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_LOGIN_USER+": %s", user.Email)

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (us *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	if len(req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": name too short")
		return dto.UserResponse{}, constants.ErrInvalidName
	}

	if !helpers.IsValidEmail(req.Email) {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": invalid email")
		return dto.UserResponse{}, constants.ErrInvalidEmail
	}

	_, found, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err == nil && found {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": email already exists")
		return dto.UserResponse{}, constants.ErrEmailAlreadyExists
	}

	if len(req.Password) < 8 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_USER + ": password too short")
		return dto.UserResponse{}, constants.ErrInvalidPassword
	}

	user := model.User{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Role:        constants.ENUM_ROLE_ADMIN,
	}

	err = us.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_USER)
		return dto.UserResponse{}, constants.ErrCreateUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_USER+": %s", user.Email)

	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}

func (us *UserService) GetAllUser(ctx context.Context, search string) ([]dto.UserResponse, error) {
	users, err := us.userRepo.GetAllUser(ctx, nil, search)

	if err != nil {
		return nil, constants.ErrGetAllUser
	}

	var datas []dto.UserResponse
	for _, user := range users {
		data := dto.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
		}

		datas = append(datas, data)
	}
	return datas, nil
}

func (us *UserService) GetAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := us.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_LIST_USER)
		return dto.UserPaginationResponse{}, constants.ErrGetAllUserWithPagination
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_LIST_USER+": page %d", req.Page)

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		datas = append(datas, dto.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
		})
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (us *UserService) GetuserByID(ctx context.Context, userID string) (dto.UserResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_GET_DETAIL_USER + ": invalid UUID")
		return dto.UserResponse{}, constants.ErrInvalidUUID
	}

	user, _, err := us.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", userID).Error(constants.MESSAGE_FAILED_GET_DETAIL_USER)
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_USER+": %s", userID)

	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, _, err := us.userRepo.GetUserByID(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", req.ID).Error(constants.MESSAGE_FAILED_UPDATE_USER)
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	if req.Name != nil && len(*req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": invalid name")
		return dto.UserResponse{}, constants.ErrInvalidName
	} else if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		if !helpers.IsValidEmail(*req.Email) {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": invalid email format")
			return dto.UserResponse{}, constants.ErrInvalidEmail
		}

		existingUser, found, err := us.userRepo.GetUserByEmail(ctx, nil, *req.Email)
		if err == nil && found && existingUser.ID != user.ID {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": email already used by other user")
			return dto.UserResponse{}, constants.ErrEmailAlreadyExists
		}

		user.Email = *req.Email
	}

	if req.Password != nil {
		if ok, _ := helpers.CheckPassword(user.Password, []byte(*req.Password)); ok {
			logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_USER + ": new password same as old")
			return dto.UserResponse{}, constants.ErrPasswordSame
		}

		hashed, err := helpers.HashPassword(*req.Password)
		if err != nil {
			logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_USER + ": hash password error")
			return dto.UserResponse{}, constants.ErrHashPassword
		}
		user.Password = hashed
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	if req.Address != nil {
		user.Address = *req.Address
	}

	err = us.userRepo.UpdateUser(ctx, nil, user)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_USER)
		return dto.UserResponse{}, constants.ErrUpdateUser
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_USER+": %s", user.ID)

	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}

func (us *UserService) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.UserResponse, error) {
	user, _, err := us.userRepo.GetUserByID(ctx, nil, req.UserID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	err = us.userRepo.DeleteUserByID(ctx, nil, req.UserID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_USER)
		return dto.UserResponse{}, constants.ErrDeleteUserByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", req.UserID)

	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}
