package domain

type StoneStatusEnum string

const (
	StoneStatusEnumPending   StoneStatusEnum = "pending"
	StoneStatusEnumCompleted StoneStatusEnum = "completed"
	StoneStatusEnumFailed    StoneStatusEnum = "failed"
)
