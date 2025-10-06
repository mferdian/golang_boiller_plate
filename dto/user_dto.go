package dto

import (
	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/model"
)

type (
	UserResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone_number"`
		Address     string    `json:"address"`
	}

	RegisterUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterUserResponse struct {
		ID    uuid.UUID `json:"id"`
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
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
	}

	UpdateUserRequest struct {
		ID          string  `json:"-"`
		Name        *string `json:"name,omitempty"`
		Email       *string `json:"email,omitempty"`
		Password    *string `json:"password,omitempty"`
		PhoneNumber *string `json:"phone_number,omitempty"`
		Address     *string `json:"address,omitempty"`
	}

	DeleteUserRequest struct {
		UserID string `json:"-"`
	}

	UserPaginationRequest struct {
		PaginationRequest
		UserID string `form:"id"`
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
