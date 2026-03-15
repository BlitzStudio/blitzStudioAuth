package types

type User struct {
	// DeviceId string `json:"deviceId" `
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
