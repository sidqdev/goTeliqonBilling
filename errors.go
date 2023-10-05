package teliqonBillling

import "errors"

var (
	ErrLowBalance                           = errors.New("low balance")
	ErrCouponAlreadyUsed                    = errors.New("coupon already used")
	ErrIncorrectValue                       = errors.New("incorrect value")
	ErrSameUser                             = errors.New("same user for transfer")
	ErrBillingUserNotFound                  = errors.New("billing user not found")
	ErrEnvironmentNotFound                  = errors.New("environment not found")
	ErrCouponNotFound                       = errors.New("coupon not found")
	ErrOutSystemServicesNotFound            = errors.New("out system service not found")
	ErrTransactionNotFound                  = errors.New("transaction not found")
	ErrSubscriptionPlanNotFound             = errors.New("subscription plan not found")
	ErrBillingUserSubscriptionNotFound      = errors.New("billing user subscription not found")
	ErrBillingUserSubscriptionAlreadyExists = errors.New("billing user subscription already exists")
	ErrTransactionCannotBeReverted          = errors.New("transaction cannot be reverted because it unsuccessful")
	ErrTransactionCannotBeProcessed         = errors.New("transaction cannot be processed because it isnot pending")
	ErrBillingUserAlreadyExists             = errors.New("billing user already exists")

	ErrPanic = errors.New("something wrong")

	ErrorsRegister = map[string]error{
		"0001": ErrLowBalance,
		"0002": ErrCouponAlreadyUsed,
		"0003": ErrSameUser,
		"0004": ErrSameUser,
		"0005": ErrBillingUserNotFound,
		"0006": ErrEnvironmentNotFound,
		"0007": ErrCouponNotFound,
		"0008": ErrOutSystemServicesNotFound,
		"0009": ErrTransactionNotFound,
		"0010": ErrSubscriptionPlanNotFound,
		"0011": ErrBillingUserSubscriptionNotFound,
		"0012": ErrBillingUserSubscriptionAlreadyExists,
		"0013": ErrTransactionCannotBeReverted,
		"0014": ErrTransactionCannotBeProcessed,
		"0015": ErrBillingUserAlreadyExists,
		"xxxx": ErrPanic,
	}

	ErrIncorrectApiUrl = errors.New("url have to ends on /")
)
