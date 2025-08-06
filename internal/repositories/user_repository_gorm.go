package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/saengepch-sys/iple-backend/internal/models"
	"strings" // สำหรับ Search และ Email Lowercase

	"gorm.io/gorm"
)

// gormUserRepository คือ struct ที่ implement UserRepository โดยใช้ GORM
type gormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository คือ Factory function สำหรับสร้าง gormUserRepository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	// ควรมีการตรวจสอบว่า db ไม่ใช่ nil ที่นี่ หรือในส่วนที่เรียกใช้
	return &gormUserRepository{db: db}
}

// --- Implementation of UserRepository methods ---

func (r *gormUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// ควรมีการ Hash Password และ Validate ข้อมูลเบื้องต้นใน Service Layer ก่อนมาถึงจุดนี้
	user.ID = 0 // Ensure DB generates ID
	// ทำให้ Email เป็นตัวพิมพ์เล็กก่อนบันทึก เพื่อให้ unique index ทำงานได้ถูกต้อง
	user.Email = strings.ToLower(user.Email)

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		// ตรวจสอบ Unique Constraint Violation (การตรวจสอบนี้อาจต้องปรับตาม Error จริงของ PostgreSQL/GORM)
		// GORM อาจคืน ErrDuplicatedKey หรืออาจต้องเช็คจาก Error message/code
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			// อาจเจาะจง constraint name เช่น "users_email_key" หรือ "idx_users_email"
			if strings.Contains(result.Error.Error(), "users_email_key") || strings.Contains(result.Error.Error(), "idx_users_email") {
				return fmt.Errorf("creating user failed: email already exists: %w", ErrDuplicateRecord)
			}
			if strings.Contains(result.Error.Error(), "users_clerk_user_id_key") { // สมมติชื่อ constraint
				return fmt.Errorf("creating user failed: clerk user id already exists: %w", ErrDuplicateRecord)
			}
			// Generic duplicate error
			return fmt.Errorf("creating user failed: %w: %w", ErrDuplicateRecord, result.Error)
		}
		return fmt.Errorf("creating user failed: %w", result.Error)
	}
	// result.RowsAffected ควรเป็น 1
	return nil
}

func (r *gormUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// ค้นหาแบบ Case-insensitive
	result := r.db.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound // คืนค่า Error ที่กำหนดเอง
		}
		return nil, fmt.Errorf("finding user by email '%s' failed: %w", email, result.Error)
	}
	return &user, nil
}

func (r *gormUserRepository) GetUserByID(ctx context.Context, id uint64, opts *UserGetOptions) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx)

	// --- Preloading Logic ---
	if opts != nil {
		if opts.PreloadEnrollments {
			// โหลด Enrollments และข้อมูล Course ที่เกี่ยวข้องกับ Enrollment นั้นๆ
			query = query.Preload("Enrollments.Course")
		}
		if opts.PreloadCourses { // TaughtCourses สำหรับ Teacher
			query = query.Preload("TaughtCourses")
		}
		if opts.PreloadLinks {
			// โหลด Links และข้อมูล User ที่เกี่ยวข้อง (Parent/Child)
			// เลือกระดับความลึกของการ Preload ตามต้องการ
			query = query.Preload("ChildrenLinks.Child", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "email") // เลือกเฉพาะฟิลด์ที่จำเป็น
			}).Preload("ParentLinks.Parent", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "email") // เลือกเฉพาะฟิลด์ที่จำเป็น
			})
		}
		// เพิ่มเงื่อนไข Preload อื่นๆ ตามต้องการ
	}

	// GORM จะใช้ Primary Key โดยอัตโนมัติเมื่อใส่ ID เป็น argument สุดท้ายของ First()
	result := query.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound // คืนค่า Error ที่กำหนดเอง
		}
		return nil, fmt.Errorf("finding user by ID %d failed: %w", id, result.Error)
	}
	return &user, nil
}

func (r *gormUserRepository) ListUsers(ctx context.Context, opts UserListOptions) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	// --- Filtering ---
	if opts.Role != nil && *opts.Role != "" {
		query = query.Where("role = ?", *opts.Role)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}
	if opts.Search != nil && *opts.Search != "" {
		searchTerm := "%" + strings.ToLower(*opts.Search) + "%"
		// ค้นหาจากหลายฟิลด์ (ปรับปรุงตามต้องการ)
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ?", searchTerm, searchTerm, searchTerm)
	}
	// Add other filters...

	// --- Counting ---
	// ทำ Count บน query ที่มี Filter แล้ว แต่ *ก่อน* Limit/Offset
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("counting users failed: %w", err)
	}

	// ถ้า totalCount เป็น 0 ก็ไม่จำเป็นต้อง Query ข้อมูลต่อ
	if totalCount == 0 {
		return users, 0, nil // Return empty slice and 0 count
	}

	// --- Pagination ---
	// ตั้งค่า Default Limit ถ้าไม่ได้ระบุมา หรือระบุมาเป็น 0 หรือค่าน้อยไป
	limit := opts.Limit
	if limit <= 0 {
		limit = 10 // Default page size
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}
	query = query.Limit(limit).Offset(offset)

	// --- Sorting ---
	// สามารถเพิ่มเงื่อนไข Sorting ตาม opts ได้ เช่น opts.SortBy = "email_asc"
	query = query.Order("created_at DESC") // Default sort

	// --- Execution ---
	err = query.Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listing users failed: %w", err)
	}

	return users, totalCount, nil
}

func (r *gormUserRepository) UpdateUser(ctx context.Context, id uint64, updates map[string]interface{}) error {
	// --- Input Sanitization (สำคัญ!) ---
	// Service layer ควร validate และ sanitize 'updates' map ก่อนส่งมาที่นี่
	// Repository ไม่ควรต้องรู้ business logic ว่าฟิลด์ไหนแก้ได้/ไม่ได้
	// แต่ Repository ควรป้องกันการแก้ไขฟิลด์ที่ไม่ควรแก้ไข เช่น ID, CreatedAt
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "password") // Password ต้องอัปเดตผ่าน Service ที่ Hash ใหม่
	delete(updates, "email")    // การเปลี่ยน Email อาจต้องมีกระบวนการ Verify แยก

	if len(updates) == 0 {
		return errors.New("no valid fields provided for update")
	}
	// ทำให้ Email ใน updates map เป็นตัวพิมพ์เล็ก ถ้ามีการส่งมา (แต่ปกติไม่ควรอนุญาตให้อัปเดต email ตรงๆ)
	// if emailVal, ok := updates["email"]; ok {
	//      if emailStr, okStr := emailVal.(string); okStr {
	//           updates["email"] = strings.ToLower(emailStr)
	//      }
	// }

	// ใช้ .Model().Where() เพื่อระบุ record และ .Updates() เพื่ออัปเดตเฉพาะฟิลด์
	// .Updates() จะไม่อัปเดต Zero Value โดย default ถ้าใช้ map,
	// ถ้าต้องการอัปเดตเป็น Zero Value ต้องใช้ .Select("*").Updates() หรือระบุ field ใน .Select()
	// GORM จะอัปเดต `updated_at` โดยอัตโนมัติ
	tx := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id)

	// ใช้ Select("*") เพื่อให้สามารถอัปเดตฟิลด์เป็น Zero Value ได้ (เช่น IsActive เป็น false)
	// หรือระบุฟิลด์ที่จะอัปเดตใน Select() ถ้าต้องการควบคุมละเอียดขึ้น
	// keys := make([]string, 0, len(updates))
	// for k := range updates {
	//     keys = append(keys, k)
	// }
	// result := tx.Select(keys).Updates(updates)
	result := tx.Updates(updates) // อัปเดตตาม Map โดยตรง (อาจไม่อัปเดต zero value) - ต้องทดสอบพฤติกรรม

	if result.Error != nil {
		// Check for unique constraint violation (e.g., if clerk_user_id is updated)
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("updating user %d failed: %w: %w", id, ErrDuplicateRecord, result.Error)
		}
		return fmt.Errorf("updating user %d failed: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		// ตรวจสอบว่า User ID นั้นมีอยู่จริงหรือไม่ เพื่อแยกแยะระหว่าง Not Found กับ No Changes Needed
		var exists int64
		checkResult := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Count(&exists)
		// ไม่สนใจ error ในการ check เพราะถ้า error ตอน update ก็น่าจะ fail ไปแล้ว
		if checkResult.Error == nil && exists == 0 {
			return fmt.Errorf("updating user %d failed: %w", id, ErrRecordNotFound)
		}
		// ถ้า user มีอยู่แต่ไม่มีแถวไหนถูก update อาจจะเพราะข้อมูลเหมือนเดิม หรือ มี GORM Hook ป้องกัน
		// ในกรณีส่วนใหญ่ การไม่เกิด error ถือว่าเพียงพอ
		return nil // No error, but no rows affected (potentially okay)
	}

	return nil
}

func (r *gormUserRepository) DeleteUser(ctx context.Context, id uint64) error {
	// GORM ทำ Soft Delete อัตโนมัติเพราะ model มี gorm.DeletedAt
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("deleting user %d failed: %w", result.Error)
	}

	// ถ้า RowsAffected เป็น 0 หมายความว่าไม่พบ User ID นั้น
	if result.RowsAffected == 0 {
		return fmt.Errorf("deleting user %d failed: %w", id, ErrRecordNotFound)
	}

	return nil
}

// (Optional) Implement GetUserByClerkID
// func (r *gormUserRepository) GetUserByClerkID(ctx context.Context, clerkID string) (*models.User, error) { ... }
