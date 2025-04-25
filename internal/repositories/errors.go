package repositories

import "errors"

// Common repository errors
var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrUpdateFailed    = errors.New("update failed")
	ErrDeleteFailed    = errors.New("delete failed")
	ErrDuplicateRecord = errors.New("duplicate record")
	// เพิ่ม error อื่นๆ ตามต้องการ
)
