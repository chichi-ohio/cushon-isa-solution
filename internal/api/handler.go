package api

import (
	"fmt"
	"net/http"
	"strconv"

	"cushion-isa/internal/db"
	"cushion-isa/internal/logger"
	"cushion-isa/internal/models"
	"cushion-isa/internal/queue"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db    *db.Client
	queue queue.QueueClient
}

func NewHandler(db *db.Client, queue queue.QueueClient) *Handler {
	return &Handler{
		db:    db,
		queue: queue,
	}
}

func (h *Handler) CreateInvestment(c *gin.Context) {
	var inv models.Investment
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set initial status
	inv.Status = "pending"

	// Save to database
	if err := h.db.CreateInvestment(&inv); err != nil {
		logger.Error("Failed to create investment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create investment"})
		return
	}

	// Publish to queue
	if err := h.queue.PublishInvestment(&inv); err != nil {
		logger.Error("Failed to publish investment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish investment"})
		return
	}

	c.JSON(http.StatusCreated, inv)
}

func (h *Handler) GetInvestment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid investment ID"})
		return
	}

	inv, err := h.db.GetInvestment(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

func (h *Handler) ListInvestments(c *gin.Context) {
	investments, err := h.db.ListInvestments()
	if err != nil {
		logger.Error("Failed to list investments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list investments"})
		return
	}

	c.JSON(http.StatusOK, investments)
}

func (h *Handler) ProcessInvestment(inv *models.Investment) error {
	logger.WithFields(logrus.Fields{
		"investment_id": inv.ID,
		"status":        "processing",
	}).Info("Processing investment")

	// Update status to processing
	if err := h.db.UpdateInvestmentStatus(inv.ID, "processing"); err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Simulate some processing work
	// TODO: Add actual investment processing logic

	// Update status to completed
	if err := h.db.UpdateInvestmentStatus(inv.ID, "completed"); err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"investment_id": inv.ID,
		"status":        "completed",
	}).Info("Investment processed successfully")

	return nil
}
