package api

import (
	"net/http"
	"strconv"

	"cepm-backend/models"
	"cepm-backend/services"

	"github.com/gin-gonic/gin"
)

type PerformanceReviewHandler struct {
	service services.PerformanceReviewService
}

func NewPerformanceReviewHandler(service services.PerformanceReviewService) *PerformanceReviewHandler {
	return &PerformanceReviewHandler{service: service}
}

// CreatePerformanceReview handles the HTTP request to create a performance review.
func (h *PerformanceReviewHandler) CreatePerformanceReview(c *gin.Context) {
	var review models.PerformanceReview

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.service.CreatePerformanceReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create performance review: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetPerformanceReview handles the HTTP request to get a single performance review by its ID.
func (h *PerformanceReviewHandler) GetPerformanceReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.service.GetPerformanceReview(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// ListUserReviews handles the HTTP request to list all reviews for a user.
func (h *PerformanceReviewHandler) ListUserReviews(c *gin.Context) {
	// TODO: In a real app, get userID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	reviews, err := h.service.ListUserReviews(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// ListTeamReviews handles the HTTP request to list all reviews for a manager's team.
func (h *PerformanceReviewHandler) ListTeamReviews(c *gin.Context) {
	// TODO: In a real app, get managerID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	managerIdStr := c.Query("managerId")
	if managerIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "managerId query parameter is required"})
		return
	}
	managerID, err := strconv.ParseUint(managerIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid managerId"})
		return
	}

	reviews, err := h.service.ListTeamReviews(uint(managerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// SubmitPerformanceReview handles the HTTP request to submit a performance review for approval.
func (h *PerformanceReviewHandler) SubmitPerformanceReview(c *gin.Context) {
	// 1. Get review ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// TODO: In a real app, get userID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	// 3. Call the service to submit the review
	if err := h.service.SubmitPerformanceReview(uint(id), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Return success response
	c.JSON(http.StatusOK, gin.H{"message": "绩效评估已成功提交审批"})
}

// ApprovePerformanceReview handles the HTTP request to approve a performance review.
func (h *PerformanceReviewHandler) ApprovePerformanceReview(c *gin.Context) {
	// 1. Get review ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// TODO: In a real app, get approverID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	approverIdStr := c.Query("approverId")
	if approverIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approverId query parameter is required"})
		return
	}
	approverID, err := strconv.ParseUint(approverIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approverId"})
		return
	}

	// Optional: Get comment from request body
	var input struct { Comment string `json:"comment"` }
	c.ShouldBindJSON(&input) // No error check needed, comment is optional

	// Call the service to approve the review
	if err := h.service.ApprovePerformanceReview(uint(id), uint(approverID), input.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "绩效评估已成功批准"})
}

// RejectPerformanceReview handles the HTTP request to reject a performance review.
func (h *PerformanceReviewHandler) RejectPerformanceReview(c *gin.Context) {
	// 1. Get review ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// TODO: In a real app, get approverID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	approverIdStr := c.Query("approverId")
	if approverIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approverId query parameter is required"})
		return
	}
	approverID, err := strconv.ParseUint(approverIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approverId"})
		return
	}

	// Optional: Get comment from request body
	var input struct { Comment string `json:"comment"` }
	c.ShouldBindJSON(&input) // No error check needed, comment is optional

	// Call the service to reject the review
	if err := h.service.RejectPerformanceReview(uint(id), uint(approverID), input.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "绩效评估已成功驳回"})
}

// ScorePerformanceReview handles the HTTP request to score a performance review.
func (h *PerformanceReviewHandler) ScorePerformanceReview(c *gin.Context) {
	// 1. Get review ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// 2. Bind the JSON request body to the ScoreInput struct
	var input services.ScoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 3. Call the service to perform the scoring
	if err := h.service.ScorePerformanceReview(uint(id), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Performance review scored successfully"})
}

// GetPerformanceReviewByPeriod handles the HTTP request to get a single performance review by user ID and period.
func (h *PerformanceReviewHandler) GetPerformanceReviewByPeriod(c *gin.Context) {
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	period := c.Query("period")
	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period query parameter is required"})
		return
	}

	review, err := h.service.GetPerformanceReviewByPeriod(uint(userID), period)
	if err != nil {
		// If there's an actual error from the service (not just record not found),
		// return an internal server error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if review == nil {
		// If review is nil, it means no record was found for the period, which is a valid scenario.
		// Return 200 OK with an empty JSON object.
		c.JSON(http.StatusOK, gin.H{})
		return
	}

			c.JSON(http.StatusOK, review)
}

// ListAllReviewsByPeriod handles the HTTP request to list all reviews for a given period, regardless of status.
func (h *PerformanceReviewHandler) ListAllReviewsByPeriod(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period query parameter is required"})
		return
	}

	reviews, err := h.service.GetAllReviewsByPeriod(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// ListAllSubmittedReviews handles the HTTP request to list all submitted reviews for HR role.
func (h *PerformanceReviewHandler) ListAllSubmittedReviews(c *gin.Context) {
	// TODO: In a real app, get userID from JWT token.
	// For now, we'll get it from a query parameter for testing.
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	reviews, err := h.service.ListAllSubmittedReviews(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdatePerformanceReview handles the HTTP request to update a performance review.
func (h *PerformanceReviewHandler) UpdatePerformanceReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.PerformanceReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Ensure the ID from the URL is used, not from the body if present
	review.ID = uint(id)

	if err := h.service.UpdatePerformanceReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update performance review: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}