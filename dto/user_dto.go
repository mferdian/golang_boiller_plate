package dto

import (
	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/model"
)

type (
	UserResponse struct {
		ID      uuid.UUID `json:"user_id"`
		Name    string    `json:"user_name"`
		Email   string    `json:"user_email"`
		NoTelp  string    `json:"user_no_telp"`
		Address string    `json:"user_address"`
	}

	RegisterUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterUserResponse struct {
		ID    uuid.UUID `json:"user_id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	CreateUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		NoTelp   string `json:"user_no_telp"`
		Address  string `json:"user_address"`
	}

	UpdateUserRequest struct {
		ID       string  `json:"-"`
		Name     *string `json:"name,omitempty"`
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
		NoTelp   *string `json:"user_no_telp,omitempty"`
		Address  *string `json:"user_address,omitempty"`
	}

	DeleteUserRequest struct {
		UserID string `json:"-"`
	}

	UserPaginationRequest struct {
		PaginationRequest
		UserID string `form:"user_id"`
	}

	UserPaginationResponse struct {
		PaginationResponse
		Data []UserResponse `json:"data"`
	}

	UserPaginationRepositoryResponse struct {
		PaginationResponse
		Users []model.User
	}
)
