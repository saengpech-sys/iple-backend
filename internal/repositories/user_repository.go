package repositories

import (
	"context"
	"iple-backend/internal/models" // **สำคัญ:** แก้ your_module_path เป็น Path จริงของโปรเจกต์คุณ
)

// UserListOptions กำหนดตัวเลือกสำหรับการดึงรายการผู้ใช้ (Pagination, Filters)
type UserListOptions struct {
	Limit  int
	Offset int
	// ตัวอย่าง Filters (เพิ่มได้ตามต้องการ)
	Role     *string // ใช้ Pointer เพื่อแยกแยะ ไม่ระบุ Filter กับ Filter ค่าว่าง
	IsActive *bool
	Search   *string // สำหรับค้นหาตามชื่อ หรือ อีเมล
}

// UserGetOptions กำหนดตัวเลือกสำหรับการดึงข้อมูลผู้ใช้คนเดียว (เช่น การ Preload)
type UserGetOptions struct {
	PreloadEnrollments bool // โหลดข้อมูลการลงทะเบียนด้วยหรือไม่
	PreloadCourses     bool // โหลดข้อมูลคอร์สที่สอน (ถ้าเป็น Teacher)
	PreloadLinks       bool // โหลดข้อมูล Parent/Child Links
}

// UserRepository คือ Interface สำหรับจัดการข้อมูล User ในฐานข้อมูล
type UserRepository interface {
	// CreateUser สร้างผู้ใช้ใหม่ในฐานข้อมูล
	CreateUser(ctx context.Context, user *models.User) error

	// GetUserByEmail ค้นหาผู้ใช้ด้วย Email (Case-insensitive ควรทำใน Query)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// GetUserByID ค้นหาผู้ใช้ด้วย ID พร้อมตัวเลือก Preload
	GetUserByID(ctx context.Context, id uint64, opts *UserGetOptions) (*models.User, error)

	// ListUsers ดึงรายการผู้ใช้ตามเงื่อนไข พร้อมจำนวนทั้งหมดสำหรับ Pagination
	ListUsers(ctx context.Context, opts UserListOptions) (users []models.User, totalCount int64, err error)

	// UpdateUser อัปเดตข้อมูลผู้ใช้เฉพาะฟิลด์ที่ระบุ
	// updates คือ map[string]interface{} ที่ key คือชื่อคอลัมน์ DB หรือชื่อฟิลด์ Struct (ตาม Naming Strategy)
	UpdateUser(ctx context.Context, id uint64, updates map[string]interface{}) error

	// DeleteUser ทำการ Soft Delete ผู้ใช้ด้วย ID
	DeleteUser(ctx context.Context, id uint64) error

	// (Optional) GetUserByClerkID ค้นหาผู้ใช้ด้วย Clerk User ID
	// GetUserByClerkID(ctx context.Context, clerkID string) (*models.User, error)
}
