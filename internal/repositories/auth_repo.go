package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-coursework/constants"
	"go-coursework/internal/dto/auth"
	"go-coursework/internal/models"
	"go-coursework/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationRepo struct {
	db *gorm.DB
}

func NewAuthenticationRepo(db *gorm.DB) *AuthenticationRepo {
	return &AuthenticationRepo{db: db}
}

func (r *AuthenticationRepo) SignIn(req *auth.SignInRequest) (resp auth.UserSignInResponse, code int, opRepo string, err error, msg string, details []string) {
	opRepo = "repositories.Auth.SignIn"

	var (
		msgNotFound                   = "User Not Found"
		msgDetailsNotFound            = []string{"Make sure you have registered", "This mostly happens because you haven't registered yet"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
		msgWrongPassword              = "Wrong Password"
		msgWrongPasswordDetails       = []string{"Make sure your password is correct", "Please press forgot password, if you really forgot your password"}
	)

	var user models.Users

	if err := r.db.
		Preload("Role").
		Preload("StudyProgram").
		Preload("ContactVerification").
		First(&user, "email = ?", req.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, fiber.StatusNotFound, opRepo, err, msgNotFound, msgDetailsNotFound
		}

		return resp, fiber.StatusInternalServerError, opRepo, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return resp, fiber.StatusUnauthorized, opRepo, err, msgWrongPassword, msgWrongPasswordDetails
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role.Name)
	if err != nil {
		return resp, fiber.StatusInternalServerError, opRepo, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	resp = auth.UserSignInResponse{
		Email:             user.Email,
		EmailVerified:     user.ContactVerification.EmailVerified,
		Telephone:         user.Telephone,
		TelephoneVerified: user.ContactVerification.TelephoneVerified,
		Role:              user.Role.Name,
		Token:             token,
	}

	return resp, fiber.StatusOK, opRepo, nil, "Successfully signed in", nil
}

func (r *AuthenticationRepo) SignUp(req *auth.SignUpRequest) (resp auth.UserSignUpResponse, code int, opRepo string, err error, msg string, details []string) {
	const op = "repositories.Auth.SignUp"

	var (
		msgEmailConflict              = "Email Already Exists"
		msgEmailConflictDetails       = []string{"Try changing the email you use"}
		msgUsernameConflict           = "Username Already Exists"
		msgUsernameConflictDetails    = []string{"Try changing the Username you use"}
		msgNotFoundStudy              = "Study Program Not Found"
		msgNotFoundStudyDetail        = []string{"Please enter the study program ID correctly"}
		msgNotFoundRole               = "Role Not Found"
		msgNotFoundRoleDetail         = []string{"Please enter the Role ID correctly"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	var existingUser models.Users
	if err := r.db.Where("username = ?", req.Username).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, fiber.StatusConflict, op, err, msgUsernameConflict, msgUsernameConflictDetails
	}

	if err := r.db.Where("email = ?", req.Email).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, fiber.StatusConflict, op, err, msgEmailConflict, msgEmailConflictDetails
	}

	var role constants.Roles
	if err := r.db.Where("id = ?", req.RoleID).First(&role).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, fiber.StatusNotFound, op, err, msgNotFoundRole, msgNotFoundRoleDetail
	}

	var studyProgram constants.StudyPrograms
	if err := r.db.Where("id = ?", req.StudyProgramID).First(&studyProgram).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, fiber.StatusNotFound, op, err, msgNotFoundStudy, msgNotFoundStudyDetail
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	req.Password = string(hashPass)

	var user = models.Users{
		Username:       req.Username,
		Email:          req.Email,
		StudyProgramID: req.StudyProgramID,
		Password:       req.Password,
		RoleID:         req.RoleID,
		Batch:          req.Batch,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	if err := r.db.
		Preload("Role").
		Preload("StudyProgram").
		Preload("ContactVerification").Where("id = ?", user.ID).
		First(&existingUser).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	resp = auth.UserSignUpResponse{
		Username:          existingUser.Username,
		Email:             existingUser.Email,
		EmailVerified:     existingUser.ContactVerification.EmailVerified,
		Telephone:         req.Telephone,
		TelephoneVerified: existingUser.ContactVerification.TelephoneVerified,
		StudyProgram:      existingUser.StudyProgram.Name,
		Role:              existingUser.Role.Name,
		Batch:             req.Batch,
		Profile:           existingUser.Profile,
	}

	var contactVerification = models.UserContactVerification{
		UserID:            user.ID,
		EmailVerified:     resp.EmailVerified,
		TelephoneVerified: resp.TelephoneVerified,
	}

	if err := r.db.Create(&contactVerification).Error; err != nil {
		return resp, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails
	}

	return resp, fiber.StatusOK, op, nil, "Successfully signed up", nil
}
