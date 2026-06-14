package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	mock_database "github.com/fastenhealth/fasten-onprem/backend/pkg/database/mock"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/models"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/web/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func changePasswordContext(t *testing.T, mockDB *mock_database.MockDatabaseRepository, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(pkg.ContextKeyTypeLogger, logrus.WithField("test", "account"))
	c.Set(pkg.ContextKeyTypeDatabase, mockDB)
	jsonData, _ := json.Marshal(body)
	c.Request, _ = http.NewRequest(http.MethodPost, "/account/password", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

// userWithPassword builds a User whose stored hash matches the given plaintext.
func userWithPassword(t *testing.T, plaintext string) *models.User {
	t.Helper()
	u := &models.User{Username: "testuser"}
	assert.NoError(t, u.HashPassword(plaintext))
	return u
}

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("succeeds with the correct current password", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mock_database.NewMockDatabaseRepository(mockCtrl)

		mockDB.EXPECT().GetCurrentUser(gomock.Any()).Return(userWithPassword(t, "oldpass"), nil)
		mockDB.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Return(nil)

		c, w := changePasswordContext(t, mockDB, gin.H{"current_password": "oldpass", "new_password": "newpass123"})
		handler.ChangePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("rejects an incorrect current password (no DB write)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mock_database.NewMockDatabaseRepository(mockCtrl)

		mockDB.EXPECT().GetCurrentUser(gomock.Any()).Return(userWithPassword(t, "oldpass"), nil)
		// UpdateUserPassword must NOT be called — gomock fails the test if it is.

		c, w := changePasswordContext(t, mockDB, gin.H{"current_password": "wrongpass", "new_password": "newpass123"})
		handler.ChangePassword(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("rejects an empty new password", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mock_database.NewMockDatabaseRepository(mockCtrl)

		mockDB.EXPECT().GetCurrentUser(gomock.Any()).Return(userWithPassword(t, "oldpass"), nil)

		c, w := changePasswordContext(t, mockDB, gin.H{"current_password": "oldpass", "new_password": "   "})
		handler.ChangePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
