package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go-coursework/internal/dto"
	"go-coursework/internal/dto/otp"
	"go-coursework/internal/logger"
	"go-coursework/internal/repositories"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"time"
)

type EmailVerificationHandler struct {
	db          *gorm.DB
	logLogrus   *logger.ErrorLogger
	redisClient *redis.Client
	otpRepo     *repositories.OtpRepo
}

func NewEmailVerification(db *gorm.DB, logLogrus *logger.ErrorLogger, redisClient *redis.Client) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		db:          db,
		logLogrus:   logLogrus,
		redisClient: redisClient,
		otpRepo:     repositories.NewOtpRepo(db),
	}
}

var (
	smtpUser string
	smtpPass string
)

func InitOtpEmail() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASSWORD")
}

func (h *EmailVerificationHandler) generateOtpEmail() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (h *EmailVerificationHandler) sendEmail(to, otp string) (dto.Response, int, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "adityamaullana234@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Kode OTP - Go-Coursework")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<title>Kode OTP Anda</title>
			<style>
				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					background-color: #f9fafb;
					margin: 0;
					padding: 0;
					color: #333;
				}
				.container {
					max-width: 600px;
					margin: 40px auto;
					background-color: #ffffff;
					padding: 30px;
					border-radius: 10px;
					box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
				}
				.header {
					text-align: center;
					border-bottom: 1px solid #eeeeee;
					padding-bottom: 20px;
					margin-bottom: 20px;
				}
				.header h1 {
					color: #2563eb;
					font-size: 24px;
					margin: 0;
				}
				.content p {
					font-size: 16px;
					line-height: 1.6;
				}
				.otp {
					font-size: 28px;
					font-weight: bold;
					color: #10b981;
					text-align: center;
					margin: 20px 0;
				}
				.footer {
					text-align: center;
					font-size: 12px;
					color: #999999;
					margin-top: 30px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Go-Coursework</h1>
				</div>
				<div class="content">
					<p>Halo,</p>
					<p>Berikut adalah <strong>Kode OTP</strong> Anda untuk melanjutkan proses verifikasi akun di <strong>Go-Coursework</strong>:</p>
					<div class="otp">%s</div>
					<p>Jangan berikan kode ini kepada siapa pun. Kode ini hanya berlaku untuk beberapa menit ke depan.</p>
				</div>
				<div class="footer">
					<p>Jika Anda tidak meminta kode ini, Anda bisa mengabaikan email ini.</p>
					<p>&copy; 2025 Go-Coursework</p>
				</div>
			</div>
		</body>
		</html>
	`, otp)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("OTP:", otp, "To:", to)
		return dto.Response{Message: "Failed to send email", Data: nil}, fiber.StatusInternalServerError, err
	}

	return dto.Response{Message: "Successfully sent email", Data: otp}, fiber.StatusOK, nil
}

func (h *EmailVerificationHandler) SendOtpEmail(ctx *fiber.Ctx) error {
	const op = "handler.otp.SendOtpEmail"

	var (
		msgBodyParse                  = "Failed to parse body"
		msgBodyParseDetails           = []string{"One of the requests is not eligible"}
		msgErrorJsonMarshal           = "Failed to marshal json"
		msgNotFoundDetails            = []string{"Sometimes it's because you haven't signed up first"}
		msgErrorSaveRedis             = "Error saving email in redis"
		msgErrorSaveRedisDetails      = []string{"There was an error while saving it in the redis database. There may be an error on the server from our side"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
	)

	var req otp.SendOtpEmail

	if err := ctx.BodyParser(&req); err != nil {
		h.logLogrus.LogUserError("-", err, msgBodyParse)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgBodyParse, msgBodyParseDetails)
	}

	msg, err := h.otpRepo.CheckEmail(req.Email)
	if err != nil {
		h.logLogrus.LogUserError(req.Email, err, msg)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusNotFound, op, err, msg, msgNotFoundDetails)
	}

	var otpEmail = h.generateOtpEmail()
	otpJson, err := json.Marshal(&otpEmail)
	if err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgErrorJsonMarshal)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}

	ctxBackground, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := h.redisClient.Set(ctxBackground, fmt.Sprintf("Otp %s: ", req.Email), otpJson, 2*time.Minute).Err(); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgErrorSaveRedis)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgErrorSaveRedis, msgErrorSaveRedisDetails)
	}

	resp, code, err := h.sendEmail(req.Email, otpEmail)
	if err != nil {
		return h.logLogrus.LogRequestError(ctx, code, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": resp.Message,
		"data":    fmt.Sprintf("Email penerima otp: %s", req.Email),
	})
}

func (h *EmailVerificationHandler) VerifyOtpEmail(ctx *fiber.Ctx) error {
	const op = "handler.EmailVerificationHandler.VerifyOtpEmail"

	var (
		msgBodyParse                  = "Failed to parse body"
		msgBodyParseDetails           = []string{"One of the requests is not eligible"}
		msgOtpNotFound                = "Otp not found in the redis database"
		msgOtpNotFoundDetails         = []string{"It looks like your OTP has expired"}
		msgInternalServerError        = "Internal Server Error"
		msgInternalServerErrorDetails = []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"}
		msgErrorJsonUnMarshal         = "Failed to unmarshal json"
		msgWrongOtpEmail              = "Wrong otp email"
		msgWrongOtpEmailDetails       = []string{"The OTP code you entered is incorrect. Please send a new OTP code."}
	)

	var req otp.VerifyOtpEmail

	if err := ctx.BodyParser(&req); err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgBodyParse)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusBadRequest, op, err, msgBodyParse, msgBodyParseDetails)
	}

	ctxBackground, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	otpJson, err := h.redisClient.Get(ctxBackground, fmt.Sprintf("Otp %s: ", req.Email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgOtpNotFound, msgOtpNotFoundDetails)
		}

		h.logLogrus.LogUserError(req.Email, err, msgInternalServerError)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, msgInternalServerError, msgInternalServerErrorDetails)
	}

	var otpEmail string
	err = json.Unmarshal([]byte(otpJson), &otpEmail)
	if err != nil {
		h.logLogrus.LogUserError(req.Email, err, msgErrorJsonUnMarshal)
	}

	if otpEmail != req.OTP {
		h.logLogrus.LogUserError(req.Email, err, msgWrongOtpEmail)
		return h.logLogrus.LogRequestError(ctx, fiber.StatusUnauthorized, op, err, msgWrongOtpEmail, msgWrongOtpEmailDetails)
	}

	var resp = otp.VerifyOtpEmail{
		Email: req.Email,
		OTP:   otpEmail,
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Successfully verified your email",
		"data":    resp,
	})
}
