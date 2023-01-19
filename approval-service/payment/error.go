package payment

import "errors"

var (
	errNonPendingStatus   = errors.New("non pending status")
	errUnknownReviewEvent = errors.New("unknown review event")
)
