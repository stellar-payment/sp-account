package dto

type RegisterCustomerPayload struct {
	LegalName    string `json:"legal_name" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Birthdate    string `json:"birth_date" validate:"required"`
	Address      string `json:"address" validate:"required"`
	PhotoProfile string `json:"photo_profile" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type RegisterMerchantPayload struct {
	LegalName    string `json:"name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	PICName      string `json:"pic_name"`
	PICEmail     string `json:"pic_email"`
	PICPhone     string `json:"pic_phone"`
	PhotoProfile string `json:"photo_profile"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
}
