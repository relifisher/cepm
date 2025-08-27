package repositories

import (
	"cepm-backend/database"
	"cepm-backend/models"

	"gorm.io/gorm"
)

type PerformanceReviewRepository interface {
	Create(review *models.PerformanceReview) error
	GetByID(id uint) (*models.PerformanceReview, error)
	ListByUserID(userID uint) ([]models.PerformanceReview, error)
	ListByManagerID(managerID uint) ([]models.PerformanceReview, error)
	ListAllSubmittedReviews() ([]models.PerformanceReview, error)
	UpdateWithItems(review *models.PerformanceReview, items []models.PerformanceItem) error
	UpdateStatus(reviewID uint, newStatus string) error
	UpdateStatusAndAddApproval(reviewID uint, newStatus string, approverID uint, comment string) error
	GetByUserIDAndPeriod(userID uint, period string) (*models.PerformanceReview, error)
	Update(review *models.PerformanceReview) error
	FindAllReviewsByPeriod(period string) ([]models.PerformanceReview, error)
}

type dbPerformanceReviewRepository struct {
	db *gorm.DB
}

func NewPerformanceReviewRepository() PerformanceReviewRepository {
	return &dbPerformanceReviewRepository{db: database.DB}
}

func (r *dbPerformanceReviewRepository) Create(review *models.PerformanceReview) error {
	return r.db.Create(review).Error
}

// GetByID retrieves a single performance review with its items and user preloaded.
func (r *dbPerformanceReviewRepository) GetByID(id uint) (*models.PerformanceReview, error) {
	var review models.PerformanceReview
	err := r.db.Preload("Items").Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// ListByUserID retrieves all performance reviews for a given user.
func (r *dbPerformanceReviewRepository) ListByUserID(userID uint) ([]models.PerformanceReview, error) {
	var reviews []models.PerformanceReview
	err := r.db.Preload("User.Department").Preload("User.Role").Preload("Items").Where("user_id = ?", userID).Order("period desc").Find(&reviews).Error
	return reviews, err
}

// ListByManagerID retrieves all performance reviews for users reporting to a given manager.
func (r *dbPerformanceReviewRepository) ListByManagerID(managerID uint) ([]models.PerformanceReview, error) {
	var reviews []models.PerformanceReview
	// Find user IDs that report to the manager
	var userIDs []uint
	r.db.Model(&models.User{}).Where("manager_id = ?", managerID).Pluck("id", &userIDs)

	if len(userIDs) == 0 {
		return reviews, nil // Return empty slice if manager has no reports
	}

	// Find all reviews for those user IDs, and preload the User info for display, excluding '草稿' status
	err := r.db.Preload("User").Where("user_id IN ? AND status != ?", userIDs, "草稿").Order("period desc").Find(&reviews).Error
	return reviews, err
}

// ListAllSubmittedReviews retrieves all performance reviews that are not in '草稿' status.
func (r *dbPerformanceReviewRepository) ListAllSubmittedReviews() ([]models.PerformanceReview, error) {
	var reviews []models.PerformanceReview
	// Define statuses that are considered 'submitted' (i.e., not '草稿')
	submittedStatuses := []string{"待审批", "已批准", "已完成", "待人事确认", "已驳回"}

	err := r.db.Preload("User.Department").Preload("User.Role").Preload("Items").Where("status IN ?", submittedStatuses).Order("period desc, user_id asc").Find(&reviews).Error
	return reviews, err
}

// UpdateWithItems updates a review and its associated items in a single transaction.
func (r *dbPerformanceReviewRepository) UpdateWithItems(review *models.PerformanceReview, items []models.PerformanceItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update each performance item
		for _, item := range items {
			if err := tx.Model(&models.PerformanceItem{}).Where("id = ?", item.ID).Updates(models.PerformanceItem{
				CompletionDetails: item.CompletionDetails,
				Score:             item.Score,
			}).Error; err != nil {
				return err
			}
		}

		// 2. Update the parent review with the total score and new status
		if err := tx.Save(review).Error; err != nil {
				return err
		}

		return nil
	})
}

// UpdateStatus updates the status of a performance review.
func (r *dbPerformanceReviewRepository) UpdateStatus(reviewID uint, newStatus string) error {
	return r.db.Model(&models.PerformanceReview{}).Where("id = ?", reviewID).Update("status", newStatus).Error
}

// UpdateStatusAndAddApproval updates the status of a review and adds an approval history entry.
func (r *dbPerformanceReviewRepository) UpdateStatusAndAddApproval(reviewID uint, newStatus string, approverID uint, comment string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update review status
		if err := tx.Model(&models.PerformanceReview{}).Where("id = ?", reviewID).Update("status", newStatus).Error; err != nil {
			return err
		}

		// Add approval history entry
		approval := models.ApprovalHistory{
			ReviewID:   reviewID,
			ApproverID: approverID,
			Status:     newStatus, // Use the new status as the approval status
			Comment:    comment,
		}
		if err := tx.Create(&approval).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByUserIDAndPeriod retrieves a single performance review for a given user and period.
func (r *dbPerformanceReviewRepository) GetByUserIDAndPeriod(userID uint, period string) (*models.PerformanceReview, error) {
	var review models.PerformanceReview
	// Preload nested associations for the user details
	err := r.db.Preload("Items").Preload("User.Department").Preload("User.Role").Where("user_id = ? AND period = ?", userID, period).First(&review).Error
	if err != nil {
		return nil, err // Can be gorm.ErrRecordNotFound
	}
	return &review, nil
}

// Update updates a review and its associated items in a single transaction.
// It replaces all old items with the new ones provided.
func (r *dbPerformanceReviewRepository) Update(review *models.PerformanceReview) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update the parent review object's top-level fields (e.g., period, status)
		// We use Select("*") to ensure all fields are updated, even if they are zero-valued.
		if err := tx.Select("*").Omit("Items").Save(review).Error; err != nil {
			return err
		}

		// 2. Delete all existing items for this review
		if err := tx.Where("review_id = ?", review.ID).Delete(&models.PerformanceItem{}).Error; err != nil {
			return err
		}

		// 3. Create the new items from the review's Items slice
		// GORM will automatically set the ReviewID for each item.
		if len(review.Items) > 0 {
			// We need to make sure the ReviewID is set on each item before creating
			for i := range review.Items {
				review.Items[i].ReviewID = review.ID
			}
			if err := tx.Create(&review.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindAllReviewsByPeriod retrieves all performance reviews for a given period, regardless of status.
func (r *dbPerformanceReviewRepository) FindAllReviewsByPeriod(period string) ([]models.PerformanceReview, error) {
	var reviews []models.PerformanceReview
	err := r.db.Preload("User.Department").Preload("User.Role").Preload("Items").Where("period = ? AND status != ?", period, "草稿").Order("user_id asc").Find(&reviews).Error
	return reviews, err
}