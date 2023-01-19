package payment

import "time"

const (
	minimumNumOfApprovalsToCompletePayment        = 2
	durationToDeactivatePaymentFromLastUpdateTime = 7 * 24 * time.Hour
	defaultPageSize                               = 50
)
