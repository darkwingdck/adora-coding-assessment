package dto

type GetMobileCarrierUserStatusCmd struct {
	UserID string
}

type GetMobileCarrierUserStatusResult struct {
	Status MobileCarrierUserStatus
}

type MobileCarrierUserStatus string

const (
	MobileCarrierUserStatusActive   MobileCarrierUserStatus = "active"
	MobileCarrierUserStatusInactive MobileCarrierUserStatus = "inactive"
	MobileCarrierUserStatusApiError MobileCarrierUserStatus = "api_error"
)
