package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/condition"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/database"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetConditionsClassified returns the classified Condition list for the authenticated user: a
// stateless compute-on-request derivation that synthesizes Condition.category and a display state,
// separating real health problems from social/administrative "Personal Health Conditions". Never
// materialized — see pkg/condition and docs/your-phr-dashboard/phase-1-condition-classifier-spec.md.
func GetConditionsClassified(c *gin.Context) {
	logger := c.MustGet(pkg.ContextKeyTypeLogger).(*logrus.Entry)
	databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)

	resources, err := databaseRepo.ListResources(c, models.ListResourceQueryOptions{SourceResourceType: "Condition"})
	if err != nil {
		logger.Errorf("error listing Condition resources for classification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	inputs := make([]condition.InputResource, 0, len(resources))
	for i := range resources {
		inputs = append(inputs, condition.InputResource{
			SourceResourceType: resources[i].SourceResourceType,
			SourceResourceID:   resources[i].SourceResourceID,
			SourceID:           resources[i].SourceID.String(),
			Raw:                json.RawMessage(resources[i].ResourceRaw),
		})
	}

	classified := condition.Classify(inputs, time.Now().UTC())
	c.JSON(http.StatusOK, gin.H{"success": true, "data": classified})
}
