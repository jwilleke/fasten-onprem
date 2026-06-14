package handler

import (
	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetCurrentUser returns the current user's profile information
func GetCurrentUser(c *gin.Context) {
	logger := c.MustGet(pkg.ContextKeyTypeLogger).(*logrus.Entry)
	databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)

	currentUser, err := databaseRepo.GetCurrentUser(c)
	if err != nil {
		logger.Errorf("Failed to get current user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get current user"})
		return
	}

	// Create a sanitized user object (without password)
	sanitizedUser := gin.H{
		"id":        currentUser.ID,
		"username":  currentUser.Username,
		"full_name": currentUser.FullName,
		"email":     currentUser.Email,
		"picture":   currentUser.Picture,
		"role":      currentUser.Role,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sanitizedUser,
	})
}

// ChangePassword updates the current user's password after verifying their current one.
func ChangePassword(c *gin.Context) {
	logger := c.MustGet(pkg.ContextKeyTypeLogger).(*logrus.Entry)
	databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request"})
		return
	}

	currentUser, err := databaseRepo.GetCurrentUser(c)
	if err != nil {
		logger.Errorf("Failed to get current user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to change password"})
		return
	}

	// Verify the current password before allowing a change.
	if err := currentUser.CheckPassword(req.CurrentPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "current password is incorrect"})
		return
	}

	// HashPassword also rejects an empty new password.
	if err := currentUser.HashPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := databaseRepo.UpdateUserPassword(c, currentUser.Password); err != nil {
		logger.Errorf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UX: this is a secure endpoint, and should only be called after a double confirmation
func DeleteAccount(c *gin.Context) {
	logger := c.MustGet(pkg.ContextKeyTypeLogger).(*logrus.Entry)
	databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)

	err := databaseRepo.DeleteCurrentUser(c)

	if err != nil {
		logger.Errorln("An error occurred while deleting current user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
