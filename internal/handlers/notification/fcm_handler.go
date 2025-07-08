package notification

import (
	"context"
	"errors"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-coursework/internal/dto/notification"
	"go-coursework/internal/models"
	pkgerr "go-coursework/pkg/errors"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"log"
)

type FCMHandler struct {
	client *messaging.Client
	ctx    context.Context
	rctx   *models.RouterContext
}

func NewFCMHandler(serviceAccountPath string, rctx *models.RouterContext) (
	*FCMHandler,
	error,
) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing FCM client: %v", err)
	}

	return &FCMHandler{
		client: client,
		ctx:    ctx,
		rctx:   rctx,
	}, nil
}

func (h *FCMHandler) SendToSingleDevice(token, title, body string) error {
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := h.client.Send(h.ctx, msg)
	if err != nil {
		return fmt.Errorf("error sending FCM message: %v", err)
	}

	log.Printf("Successfully sent message: %s\n", response)
	return nil
}

func (h *FCMHandler) SaveFCMToken(ctx *fiber.Ctx) error {
	const op = "handler.NewFCMToken.SaveFCMToken"

	var req notification.SaveFCMTokenReq
	if err := ctx.BodyParser(&req); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusBadRequest, op, err, pkgerr.ErrBodyParse.Message, pkgerr.ErrBodyParse.Details)
	}

	var fcm = models.FCMTokens{
		UserID: req.UserID,
		Token:  req.Token,
	}

	if err := h.rctx.DB.Create(fcm).Error; err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	if err := h.SendToSingleDevice(req.Token, "Successfully SignUp Account", "Thank you for your SignUp account go-coursework! Good Job Yeah...."); err != nil {
		return h.rctx.Logger.LogRequestError(ctx, fiber.StatusInternalServerError, op, err, pkgerr.ErrInternalServer.Message, pkgerr.ErrInternalServer.Details)
	}

	return nil
}

func (h *FCMHandler) GetFCMToken(userID int) (
	string,
	error,
) {
	var fcm models.FCMTokens
	if err := h.rctx.DB.Where("user_id = ?", userID).First(&fcm).Error; err != nil {
		return "", fmt.Errorf("error getting FCM tokens: %v", err)
	}

	return fcm.Token, nil
}

func (h *FCMHandler) UpdateFCMToken(userID int, token string) error {
	var fcm models.FCMTokens
	if err := h.rctx.DB.Where("user_id = ?", userID).First(&fcm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error getting FCM tokens, userID is not found: %v", err)
		}
		return fmt.Errorf("error getting FCM tokens: %v", err)
	}

	fcm.Token = token
	if err := h.rctx.DB.Save(&fcm).Error; err != nil {
		return fmt.Errorf("error updating FCM tokens: %v", err)
	}

	return nil
}

func (h *FCMHandler) DeleteFCMToken(userID int) error {
	// this method are use for delete fcm token, if user delete application
	return nil
}
