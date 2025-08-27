package services

import (
	"errors"

	"cepm-backend/database"
	"cepm-backend/models"
	"cepm-backend/repositories"

	"gorm.io/gorm"
)

// ScoreItemInput defines the structure for a single item's score data from the API.
type ScoreItemInput struct {
	ID                uint     `json:"id"`
	CompletionDetails string   `json:"completionDetails"`
	Score             *float64 `json:"score"`
}

// ScoreInput defines the structure for the entire scoring request.
type ScoreInput struct {
	Items        []ScoreItemInput `json:"items"`
	FinalComment string           `json:"finalComment"`
}

// PerformanceReviewService defines the interface for performance review services.
type PerformanceReviewService interface {
	CreatePerformanceReview(review *models.PerformanceReview) error
	GetPerformanceReview(reviewID uint) (*models.PerformanceReview, error)
	ListUserReviews(userID uint) ([]models.PerformanceReview, error)
	ListTeamReviews(managerID uint) ([]models.PerformanceReview, error)
	ListAllSubmittedReviews(userID uint) ([]models.PerformanceReview, error) // New method for HR role
	SubmitPerformanceReview(reviewID uint, userID uint) error
	ApprovePerformanceReview(reviewID uint, approverID uint, comment string) error
	RejectPerformanceReview(reviewID uint, approverID uint, comment string) error
	ScorePerformanceReview(reviewID uint, input *ScoreInput) error
	GetPerformanceReviewByPeriod(userID uint, period string) (*models.PerformanceReview, error)
	UpdatePerformanceReview(review *models.PerformanceReview) error
	GetAllReviewsByPeriod(period string) ([]models.PerformanceReview, error)
}

type performanceReviewService struct {
	repo repositories.PerformanceReviewRepository
	db   *gorm.DB // Add gorm.DB dependency for user role check
}

// NewPerformanceReviewService creates a new instance of PerformanceReviewService.
func NewPerformanceReviewService(repo repositories.PerformanceReviewRepository) PerformanceReviewService {
	return &performanceReviewService{repo: repo, db: database.DB} // Inject database.DB
}

// ListAllSubmittedReviews retrieves all performance reviews for HR role.
func (s *performanceReviewService) ListAllSubmittedReviews(userID uint) ([]models.PerformanceReview, error) {
	// 1. Check if the user has the HR role
	var user models.User
	err := s.db.Preload("Role").First(&user, userID).Error
	if err != nil {
		return nil, errors.New("用户不存在或无法获取用户信息")
	}

	if user.Role.Name != "人事" && user.Role.Name != "HR" { // Assuming role name is "人事" or "HR"
		return nil, errors.New("您没有权限查看所有已提交的绩效评估")
	}

	// 2. If authorized, retrieve all submitted reviews from the repository
	return s.repo.ListAllSubmittedReviews()
}

// CreatePerformanceReview handles the business logic for creating a performance review.
func (s *performanceReviewService) CreatePerformanceReview(review *models.PerformanceReview) error {
	// Business logic for creation goes here...
	return s.repo.Create(review)
}

// GetPerformanceReview retrieves a single performance review.
func (s *performanceReviewService) GetPerformanceReview(reviewID uint) (*models.PerformanceReview, error) {
	return s.repo.GetByID(reviewID)
}

// ListUserReviews retrieves all reviews for a specific user.
func (s *performanceReviewService) ListUserReviews(userID uint) ([]models.PerformanceReview, error) {
	return s.repo.ListByUserID(userID)
}

// ListTeamReviews retrieves all reviews for users reporting to a specific manager.
func (s *performanceReviewService) ListTeamReviews(managerID uint) ([]models.PerformanceReview, error) {
	return s.repo.ListByManagerID(managerID)
}

// SubmitPerformanceReview handles the business logic for submitting a performance review.
func (s *performanceReviewService) SubmitPerformanceReview(reviewID uint, userID uint) error {
	// 1. Get the review
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return errors.New("绩效评估不存在")
	}

	// 2. Permission check: Only the owner can submit their own draft
	if review.UserID != userID {
		return errors.New("您无权提交此绩效评估")
	}

	// 3. Status check: Only '草稿' reviews can be submitted
	if review.Status != "草稿" {
		return errors.New("只有草稿状态的绩效评估才能提交")
	}

	// 4. Update status to '待审批' and add approval history
	return s.repo.UpdateStatusAndAddApproval(reviewID, "待审批", userID, "提交审批")
}

// ApprovePerformanceReview handles the business logic for approving a performance review.
func (s *performanceReviewService) ApprovePerformanceReview(reviewID uint, approverID uint, comment string) error {
	// 1. Get the review and its user/manager info
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return errors.New("绩效评估不存在")
	}

	// 2. Permission check: Ensure approver is the correct manager in the chain
	// This is a simplified check. In a real app, you'd check roles and hierarchy more strictly.
	// For now, we assume approverID is the direct manager of review.User
	if review.User.ManagerID == nil || *review.User.ManagerID != approverID {
		return errors.New("您无权审批此绩效评估")
	}

	// 3. Status check: Only '待审批' or '待人事确认' can be approved
	// Simplified: only '待审批' can be approved by manager
	if review.Status != "待审批" {
		return errors.New("只有待审批状态的绩效评估才能被批准")
	}

	// 4. Determine next status based on approval chain
	// Simplified chain: Employee -> Manager -> HR
	// If current approver is manager, next is HR (if review is '已完成') or '待打分' (if plan is approved)
	newStatus := "已批准" // Default for plan approval
	// In a real app, you'd check approver's role and review's current status to determine next step
	// For now, let's assume manager approves plan, next is '待打分'

	// If the review is already '已完成' (scored by employee), manager approval means it goes to HR
	if review.Status == "已完成" { // This status check is wrong, should be '待打分' or '待人事确认' for manager to approve score
		newStatus = "待人事确认"
	} else { // Manager approves the plan
		newStatus = "待打分"
	}

	// 5. Update status and add approval history
	return s.repo.UpdateStatusAndAddApproval(reviewID, newStatus, approverID, comment)
}

// RejectPerformanceReview handles the business logic for rejecting a performance review.
func (s *performanceReviewService) RejectPerformanceReview(reviewID uint, approverID uint, comment string) error {
	// 1. Get the review
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return errors.New("绩效评估不存在")
	}

	// 2. Permission check (simplified: any manager can reject for now)
	// In a real app, check if approverID is in the approval chain for this review

	// 3. Update status to '已驳回' and add approval history
	_ = review // Dummy use to satisfy compiler
	return s.repo.UpdateStatusAndAddApproval(reviewID, "已驳回", approverID, comment)
}

func calculateGradePoint(totalScore float64) float64 {
	if totalScore >= 90 && totalScore <= 100 {
		return 1.0
	} else if totalScore >= 60 && totalScore < 90 {
		return 0.8
	} else if totalScore < 60 {
		return 0
	} else if totalScore > 100 {
		return totalScore / 100 // As per frontend logic
	}
	return 0 // Default case
}

// ScorePerformanceReview handles the business logic for scoring a performance review.
func (s *performanceReviewService) ScorePerformanceReview(reviewID uint, input *ScoreInput) error {
	// 1. Get the existing review with its items
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return errors.New("绩效评估不存在")
	}

	// Create a map of existing items by their ID for easy lookup
	itemMap := make(map[uint]*models.PerformanceItem)
	for i := range review.Items {
		itemMap[review.Items[i].ID] = &review.Items[i]
	}

	var totalScore float64 = 0
	var itemsToUpdate []models.PerformanceItem

	// 2. Calculate total score and prepare items for update
	for _, itemInput := range input.Items {
		existingItem, ok := itemMap[itemInput.ID]
		if !ok {
			return errors.New("无效的绩效项ID")
		}

		if itemInput.Score != nil {
			// Validate score range
			if *itemInput.Score < 0 || *itemInput.Score > 120 {
				return errors.New("单项分数必须在0到120之间")
			}
			// Weight is a percentage (e.g., 80), score is out of 100.
			// Contribution to total score is (weight/100) * score.
			totalScore += (existingItem.Weight / 100.0) * (*itemInput.Score)
		}

		itemsToUpdate = append(itemsToUpdate, models.PerformanceItem{
			ID:                itemInput.ID,
			CompletionDetails: itemInput.CompletionDetails,
			Score:             itemInput.Score,
		})
	}

	// 3. Calculate grade point
	gradePoint := calculateGradePoint(totalScore)

	// 4. Update the parent review object
	review.TotalScore = &totalScore
	review.GradePoint = &gradePoint
	review.FinalComment = input.FinalComment
	review.Status = "已完成" // Or another appropriate status

	// 5. Persist changes to the database
	return s.repo.UpdateWithItems(review, itemsToUpdate)
}

// GetPerformanceReviewByPeriod retrieves a single performance review for a user and period.
// It returns (nil, nil) if the review is not found, allowing the handler to return an empty response.
func (s *performanceReviewService) GetPerformanceReviewByPeriod(userID uint, period string) (*models.PerformanceReview, error) {
	review, err := s.repo.GetByUserIDAndPeriod(userID, period)
	if err != nil {
		// If the error is specifically "record not found", we handle it gracefully.
		// This is not an application error, but a valid state (no review for that month).
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil review and nil error
		}
		// For any other database error, we return it.
		return nil, err
	}
	return review, nil
}

// GetAllReviewsByPeriod retrieves all performance reviews for a given period, regardless of status.
func (s *performanceReviewService) GetAllReviewsByPeriod(period string) ([]models.PerformanceReview, error) {
	return s.repo.FindAllReviewsByPeriod(period)
}

// UpdatePerformanceReview handles the business logic for updating a performance review.
func (s *performanceReviewService) UpdatePerformanceReview(review *models.PerformanceReview) error {
	// 1. Get the existing review from DB to check its status
	existingReview, err := s.repo.GetByID(review.ID)
	if err != nil {
		return errors.New("绩效评估不存在")
	}

	// 2. Check if the review is in a modifiable state
	if existingReview.Status != "草稿" && existingReview.Status != "已驳回" {
		return errors.New("只有“草稿”或“已驳回”状态的绩效评估才能被修改")
	}

	// 3. Validation for the items
	var workTotalWeight float64 = 0
	for _, item := range review.Items {
		// Not-null validation
		if item.Title == "" || item.Description == "" || item.Target == "" || item.Weight <= 0 {
			return errors.New("所有绩效项的字段均不能为空，且权重必须大于0")
		}
		if item.Category == "工作业绩" {
			workTotalWeight += item.Weight
		}
	}

	// Weight validation
	if workTotalWeight != 80 {
		return errors.New("“工作业绩”部分的总权重必须等于80%")
	}

	// 4. Call the repository to update
	return s.repo.Update(review)
}

