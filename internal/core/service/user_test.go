package service

import (
	"context"
	"errors"
	"testing"

	"project-api/internal/core/entity"
	"project-api/internal/core/model/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*entity.User) // Ensure that it's a pointer to `entity.User`
	if user == nil {
		// If the user is nil, ensure we return an error that matches the test case expectation
		return nil, errors.New("user not found")
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetUserByName(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	user, _ := args.Get(0).(*entity.User) // Same handling for GetUserByName
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, args.Error(1)
}

func TestNewUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	want := &UserService{repo: mockRepo}

	got := NewUserService(mockRepo)

	assert.Equal(t, want, got)
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := &UserService{repo: mockRepo}
	ctx := context.Background()

	tests := []struct {
		name      string
		userReq   *request.UserRequest
		wantErr   error                                                            // Expecting an error, could be nil
		mockErr   error                                                            // Error from the repository
		setupMock func(repo *MockUserRepository, user *entity.User, mockErr error) // Function to set up the mock
	}{
		{
			name: "Success",
			userReq: &request.UserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: nil,
			mockErr: nil,
			setupMock: func(repo *MockUserRepository, user *entity.User, mockErr error) {
				repo.On("Create", ctx, user).Return(mockErr)
			},
		},
		{
			name: "Validation Error - Empty Email",
			userReq: &request.UserRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: request.ErrEmailRequired, // Expecting specific validation error
			mockErr: nil,                      // No mock error for validation
			setupMock: func(repo *MockUserRepository, user *entity.User, mockErr error) {
				// No mock setup needed for validation errors
			},
		},
		{
			name: "Validation Error - Invalid Email",
			userReq: &request.UserRequest{
				Email:    "invalid-email",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: request.ErrInvalidEmail, // Expecting specific validation error
			mockErr: nil,                     // No mock error for validation
			setupMock: func(repo *MockUserRepository, user *entity.User, mockErr error) {
				// No mock setup needed for validation errors
			},
		},
		{
			name: "Validation Error - Empty Password",
			userReq: &request.UserRequest{
				Email:    "test@example.com",
				Username: "testuser",
			},
			wantErr: request.ErrPasswordRequired,
			mockErr: nil,
			setupMock: func(repo *MockUserRepository, user *entity.User, mockErr error) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userEntity, _ := tt.userReq.ToEntity() // Convert to entity

			tt.setupMock(mockRepo, userEntity, tt.mockErr) // Set up the mock

			err := u.Create(ctx, tt.userReq)

			assert.ErrorIs(t, err, tt.wantErr, "Error mismatch for test: %s", tt.name) // Check if the error *is* the expected error or wraps it

			if tt.mockErr != nil {
				mockRepo.AssertExpectations(t) // Only assert expectations if the mock was used
			}
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := &UserService{repo: mockRepo}
	ctx := context.Background()
	email := "test@example.com"
	expectedUser := &entity.User{Email: email, UserName: "testuser"}

	tests := []struct {
		name     string
		email    string // Use email directly
		want     *entity.User
		wantErr  bool
		mockUser *entity.User
		mockErr  error
	}{
		{
			name:     "User Found",
			email:    email,
			want:     expectedUser,
			wantErr:  false,
			mockUser: expectedUser,
			mockErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUserByEmail", ctx, tt.email).Return(tt.mockUser, tt.mockErr)

			got, err := u.GetUserByEmail(ctx, tt.email)

			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Equal(t, tt.mockErr, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByUserName(t *testing.T) {
	mockRepo := new(MockUserRepository)
	u := &UserService{repo: mockRepo}
	ctx := context.Background()
	username := "testuser"
	expectedUser := &entity.User{UserName: username, Email: "test@example.com"}

	tests := []struct {
		name     string
		username string // Use username directly
		want     *entity.User
		wantErr  bool
		mockUser *entity.User
		mockErr  error
	}{
		{
			name:     "User Found",
			username: username,
			want:     expectedUser,
			wantErr:  false,
			mockUser: expectedUser,
			mockErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUserByName", ctx, tt.username).Return(tt.mockUser, tt.mockErr)

			got, err := u.GetUserByUserName(ctx, tt.username)

			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Equal(t, tt.mockErr, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
