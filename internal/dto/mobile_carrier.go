package dto

type GetMobileCarrierUserStatusCmd struct {
	UserID string
}

type GetMobileCarrierUserStatusResult struct {
	Status MobileCarrierUserStatus
}
