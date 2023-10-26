package indto

// This DTO is copies of actual dto used for JSON marshaling
// and unmarshaling using native Go's struct

type Customer struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	LegalName    string `json:"legal_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Birthdate    string `json:"birth_date"`
	Address      string `json:"address"`
	PhotoProfile string `json:"photo_profile"`
}

type Merchant struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	LegalName    string `json:"name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	PICName      string `json:"pic_name"`
	PICEmail     string `json:"pic_email"`
	PICPhone     string `json:"pic_phone"`
	PhotoProfile string `json:"photo_profile"`
}
