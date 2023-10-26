package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-account/internal/inconst"
	"github.com/stellar-payment/sp-account/internal/indto"
	"github.com/stellar-payment/sp-account/internal/model"
	"github.com/stellar-payment/sp-account/internal/util/structutil"
	"github.com/stellar-payment/sp-account/pkg/dto"
	"github.com/stellar-payment/sp-account/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) RegisterCustomer(ctx context.Context, payload *dto.RegisterCustomerPayload) (err error) {
	logger := log.Ctx(ctx)

	if val := structutil.CheckMandatoryField(payload); val != "" {
		logger.Error().Msgf("field %s is missing a value", val)
		return errs.New(errs.ErrMissingRequiredAttribute, val)
	}

	if exists, err := s.repository.FindUser(ctx, &indto.UserParams{Username: payload.Username}); err != nil {
		logger.Error().Err(err).Send()
		return err
	} else if exists != nil {
		return errs.ErrDuplicatedResources
	}

	userModel := &model.User{
		UserID:   uuid.NewString(),
		Username: payload.Username,
		Password: "",
		RoleID:   inconst.ROLE_CUSTOMER,
	}

	if hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost); err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return err
	} else {
		userModel.Password = string(hashed)
	}

	err = s.repository.InsertUser(ctx, userModel)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	err = s.publishEvent(ctx, inconst.TOPIC_CREATE_CUSTOMER, &indto.Customer{
		UserID:       userModel.UserID,
		LegalName:    payload.LegalName,
		Phone:        payload.Phone,
		Email:        payload.Email,
		Birthdate:    payload.Birthdate,
		Address:      payload.Address,
		PhotoProfile: payload.PhotoProfile,
	})
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
