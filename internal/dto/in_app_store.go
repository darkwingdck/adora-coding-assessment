package dto

type UpdateUserEntitlementCmd struct {
	EventID     string
	UserID      string
	Type        EventType
	EventTimeMs int64
	ProductID   ProductID
}
