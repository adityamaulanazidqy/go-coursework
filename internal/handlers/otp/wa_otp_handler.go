package otp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	notification2 "go-coursework/internal/dto/notification"
	"go-coursework/internal/handlers/notification"
	"go-coursework/internal/helpers"
	"go-coursework/internal/logger"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type WaVerificationHandler struct {
	db *gorm.DB
	rc *redis.Client
	lg *logger.ErrorLogger
	fh *notification.FonteHandler
}

func NewWaVerificationHandler(rctx *models.RouterContext) *WaVerificationHandler {
	return &WaVerificationHandler{
		db: rctx.DB,
		rc: rctx.RedisClient,
		lg: rctx.Logger,
		fh: notification.NewFonteHandler(rctx.DB),
	}
}

func (h *WaVerificationHandler) generateToken() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(fmt.Sprintf("%06d", rand.Intn(1000000)))
}

func (h *WaVerificationHandler) VerificationWaTokenService(ctx *fiber.Ctx) error {
	const op = "handler.NewWaVerificationHandler.TokenService"

	u, e := helpers.GetUserClaims(ctx)
	if e != nil {
		h.lg.Logger.Error(fmt.Sprintf("error get user claims from token service: %s", e))
		return h.lg.LogRequestError(ctx, fiber.StatusInternalServerError, op, e, pkgerr.ErrMissingClaims.Message, pkgerr.ErrMissingClaims.Details)
	}

	var req notification2.WaVerificationReq
	if e := ctx.BodyParser(&req); e != nil {
		h.lg.LogUserError(u.Email, e, pkgerr.ErrBodyParse.Message)
		return h.lg.LogRequestError(ctx, fiber.StatusBadRequest, op, e, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	m, e := h.VerificationWaToken(u.UserID, strconv.Itoa(req.Otp))
	if e != nil {
		h.lg.LogUserError(u.Email, e, pkgerr.ErrVerificationOTPWa.Message)
		return h.lg.LogRequestError(ctx, fiber.StatusBadRequest, op, e, pkgerr.ErrVerificationOTPWa.Message, pkgerr.ErrVerificationOTPWa.Details)
	}

	var ucv models.UserContactVerification
	if err := h.db.Model(&ucv).Where("user_id = ?", u.UserID).Update("telephone_verified", true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return h.lg.LogRequestError(ctx, fiber.StatusNotFound, op, err, pkgerr.ErrUserNotFound.Message, pkgerr.ErrUserNotFound.Details)
		}

		h.lg.LogUserError(u.Email, fiber.ErrInternalServerError, pkgerr.ErrInternalServer.Message)
		return h.lg.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": m,
	})
}

func (h *WaVerificationHandler) SendMessageWaOTP(i int, n int) (
	string,
	error,
) {
	t := h.generateToken()

	if err := godotenv.Load("../../../.env"); err != nil {
		h.lg.Logger.Error(fmt.Sprintf("error loading .env file for user id: %d", i))
		return "error load .env file", err
	}

	tf, err := strconv.Atoi(os.Getenv("TOKEN_FONTE"))
	if err != nil {
		h.lg.Logger.Error(fmt.Sprintf("error converting TOKEN_FONTE to integer for user id: %d", i))
		return "error formater no telephone str to int", err
	}

	if err := h.fh.SendMessage(n, tf, fmt.Sprintf("Your code otp is: %s", t)); err != nil {
		h.lg.Logger.Error(fmt.Sprintf("error send message to user id: %d", i))
		return "error send message", err
	}

	if err := h.rc.Set(context.Background(), fmt.Sprintf("%d:%s", n, t), false, 2*time.Minute).Err(); err != nil {
		h.lg.Logger.Error(fmt.Sprintf("error set telephone in redis db from user id: %d", i))
		return "error set telephone in redis", err
	}

	return "success send message", nil
}

func (h *WaVerificationHandler) VerificationWaToken(i int, t string) (
	string,
	error,
) {
	n, err := h.fh.GetNumber(i)
	if err != nil {
		h.lg.Logger.Error(fmt.Sprintf("error get no. telephone for user id: %d", i))
		return "error get no. telephone", err
	}

	var c int64
	if c, err = h.rc.Exists(context.Background(), fmt.Sprintf("%d:%s", n, t)).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			h.lg.Logger.Error(fmt.Sprintf("error check telephone in redis, otp is expired for  user id: %d", i))
			return "error check code otp telephone, your otp is expired or not valid", err
		}

		h.lg.Logger.Error(fmt.Sprintf("error check telephone in redis db from user id: %d", i))
		return "error check code otp telephone", err
	}

	ts := strings.Split(fmt.Sprintf("%d", c), ":")

	if t != ts[1] {
		h.lg.Logger.Error(fmt.Sprintf("error otp telephone in redis, otp telephone is wrong: %s", t))
		return "error your code otp not valid", err
	}

	return "success verification otp telephone", nil
}
