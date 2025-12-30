package service

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/constants"
	"github.com/mferdian/golang_boiller_plate/dto"
	"github.com/mferdian/golang_boiller_plate/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	registerFn             func(ctx context.Context, user model.User) error
	getByIDFn              func(ctx context.Context, userID string) (model.User, bool, error)
	getByEmailFn           func(ctx context.Context, email string) (model.User, bool, error)
	getAllFn               func(ctx context.Context, search string) ([]model.User, error)
	getAllWithPaginationFn func(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationRepositoryResponse, error)
	createFn               func(ctx context.Context, user model.User) error
	updateFn               func(ctx context.Context, user model.User) error
	deleteByIDFn           func(ctx context.Context, userID string) error
}

func (m *mockUserRepo) Register(ctx context.Context, _ *gorm.DB, user model.User) error {
	if m.registerFn != nil {
		return m.registerFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, _ *gorm.DB, userID string) (model.User, bool, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, userID)
	}
	return model.User{}, false, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, _ *gorm.DB, email string) (model.User, bool, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return model.User{}, false, nil
}

func (m *mockUserRepo) GetAllUser(ctx context.Context, _ *gorm.DB, search string) ([]model.User, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, search)
	}
	return []model.User{}, nil
}

func (m *mockUserRepo) GetAllUserWithPagination(ctx context.Context, _ *gorm.DB, req dto.UserPaginationRequest) (dto.UserPaginationRepositoryResponse, error) {
	if m.getAllWithPaginationFn != nil {
		return m.getAllWithPaginationFn(ctx, req)
	}
	return dto.UserPaginationRepositoryResponse{}, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, _ *gorm.DB, user model.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, _ *gorm.DB, user model.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) DeleteUserByID(ctx context.Context, _ *gorm.DB, userID string) error {
	if m.deleteByIDFn != nil {
		return m.deleteByIDFn(ctx, userID)
	}
	return nil
}

type mockJWTService struct {
	generateFn      func(userID, role string) (string, string, error)
	validateTokenFn func(token string) (*jwt.Token, *jwtCustomClaims, error)
}

func (m *mockJWTService) GenerateToken(userID, role string) (string, string, error) {
	if m.generateFn != nil {
		return m.generateFn(userID, role)
	}
	return "access", "refresh", nil
}

func (m *mockJWTService) ValidateToken(token string) (*jwt.Token, *jwtCustomClaims, error) {
	if m.validateTokenFn != nil {
		return m.validateTokenFn(token)
	}
	return nil, nil, nil
}

// Unit Test

// Register
func TestUserService_Register_InvalidName(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "abc",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != constants.ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}
func TestUserService_Register_InvalidEmail(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "invalid-email",
		Password: "password123",
	})

	if err != constants.ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}
func TestUserService_Register_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, true, nil
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != constants.ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}
func TestUserService_Register_PasswordToShort(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "12345",
	})

	if err != constants.ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
func TestUserService_Register_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, nil
		},
		registerFn: func(ctx context.Context, user model.User) error {
			return nil
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	resp, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Email != "test@mail.com" {
		t.Errorf("expected email test@mail.com, got %s", resp.Email)
	}
}
func TestUserService_Register_GetUserByEmailError(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, constants.ErrInternal
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != constants.ErrInternal {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}
func TestUserService_Register_Failed(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, nil
		},
		registerFn: func(ctx context.Context, user model.User) error {
			return constants.ErrRegisterUser
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != constants.ErrRegisterUser {
		t.Fatalf("expected ErrRegisterUser got %v", err)
	}
}

// Login
func TestUserService_Login_Success(t *testing.T) {
	userID := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     constants.ENUM_ROLE_USER,
			}, true, nil
		},
	}

	jwt := &mockJWTService{
		generateFn: func(userID, role string) (string, string, error) {
			return "access", "refresh", nil
		},
	}

	us := NewUserService(repo, jwt)

	resp, err := us.Login(context.Background(), dto.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("expected access token")
	}
}
func TestUserService_Login_PasswordMismatch(t *testing.T) {
	userID := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     constants.ENUM_ROLE_USER,
			}, true, nil
		},
	}

	jwt := &mockJWTService{
		generateFn: func(userID, role string) (string, string, error) {
			t.Fatal("GenerateToken should NOT be called")
			return "", "", nil
		},
	}

	us := NewUserService(repo, jwt)

	_, err = us.Login(context.Background(), dto.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "password12344",
	})

	if err == nil {
		t.Fatalf("unexpected error")
	}

	if err != constants.ErrInvalidLoginCredential {
		t.Fatalf("expected ErrInvalidLoginCredential, got %v", err)
	}
}
func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, nil
		},
	}

	jwt := &mockJWTService{
		generateFn: func(userID, role string) (string, string, error) {
			t.Fatal("GenerateToken should NOT be called")
			return "", "", nil
		},
	}

	us := NewUserService(repo, jwt)

	_, err := us.Login(context.Background(), dto.LoginUserRequest{
		Email:    "notfound@mail.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != constants.ErrInvalidLoginCredential {
		t.Fatalf("expected ErrInvalidLoginCredential, got %v", err)
	}
}
func TestUserService_Login_GenerateTokenError(t *testing.T) {
	userID := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{
				ID:       userID,
				Email:    email,
				Password: string(hashedPassword),
				Role:     constants.ENUM_ROLE_USER,
			}, true, nil
		},
	}

	jwt := &mockJWTService{
		generateFn: func(userID, role string) (string, string, error) {
			return "", "", constants.ErrGenerateAccessToken
		},
	}

	us := NewUserService(repo, jwt)

	_, err = us.Login(context.Background(), dto.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != constants.ErrGenerateAccessToken {
		t.Fatalf("expected ErrGenerateAccessToken, got %v", err)
	}
}

// Create User
func TestUserService_CreateUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, nil
		},
		registerFn: func(ctx context.Context, user model.User) error {
			return nil
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	resp, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "Som User",
		Email:       "test@mail.com",
		Password:    "password123",
		PhoneNumber: "0862323213331",
		Address:     "San Diego, Mexico Selatan",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Email != "test@mail.com" {
		t.Errorf("expected email test@mail.com, got %s", resp.Email)
	}
}
func TestUserService_CreateUser_InvalidName(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "abc",
		Email:       "test@mail.com",
		Password:    "password123",
		PhoneNumber: "0862323213331",
		Address:     "San Diego, Mexico Selatan",
	})

	if err != constants.ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}
func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "Som User",
		Email:       "invalid-email",
		Password:    "password123",
		PhoneNumber: "081543625372",
		Address:     "Mexico Barat Selatan Timur dan Utara",
	})

	if err != constants.ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}
func TestUserService_CreateUser_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, true, nil
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "Som User",
		Email:       "test@mail.com",
		Password:    "password123",
		PhoneNumber: "087637347123",
		Address:     "Mexico Utara Selatan Barat",
	})

	if err != constants.ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}
func TestUserService_CreateUser_PasswordToShort(t *testing.T) {
	us := NewUserService(
		&mockUserRepo{},
		&mockJWTService{},
	)

	_, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "Som User",
		Email:       "test@mail.com",
		Password:    "aaaa",
		PhoneNumber: "0854983283923",
		Address:     "Mexico Utara Selatan Barat TImur",
	})

	if err != constants.ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
func TestUserService_CreateUser_GetUserByEmailError(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, constants.ErrInternal
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.CreateUser(context.Background(), dto.CreateUserRequest{
		Name:        "Som User",
		Email:       "test@mail.com",
		Password:    "password123",
		PhoneNumber: "085626372638",
		Address:     "Mexico Selatan Barat Timur Utara",
	})

	if err != constants.ErrInternal {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}
func TestUserService_CreateUser_Failed(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (model.User, bool, error) {
			return model.User{}, false, nil
		},
		registerFn: func(ctx context.Context, user model.User) error {
			return constants.ErrRegisterUser
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	_, err := us.Register(context.Background(), dto.RegisterUserRequest{
		Name:     "Som User",
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != constants.ErrRegisterUser {
		t.Fatalf("expected ErrRegisterUser got %v", err)
	}
}


// GetAllUser
func TestUserService_GetAllUser_Error(t *testing.T) {
	repo := &mockUserRepo{
		getAllFn: func(ctx context.Context, search string) ([]model.User, error) {
			return nil, constants.ErrInternal
		},
	}

	us := NewUserService(repo, &mockJWTService{})

	resp, err := us.GetAllUser(context.Background(), "som")

	if err != constants.ErrGetAllUser {
		t.Fatalf("expected ErrGetAllUser, got %v", err)
	}

	if resp != nil {
		t.Fatalf("expected nil response, got %v", resp)
	}
}
