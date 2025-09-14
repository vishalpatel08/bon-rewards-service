package models

import "time"

type BillStatus string

// Possible statuses for a Bill
const (
	StatusUnpaid     BillStatus = "UNPAID"
	StatusPaidOnTime BillStatus = "PAID_ON_TIME"
	StatusPaidLate   BillStatus = "PAID_LATE"
)

type Bill struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Amount      int64      `json:"amount"`
	DueDate     time.Time  `json:"due_date"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
	Status      BillStatus `json:"status"`
}
