package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fastenhealth/fasten-onprem/backend/pkg"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/database"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/document"
	"github.com/fastenhealth/fasten-onprem/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetDocumentsClassified returns the patient's DocumentReference list with a synthesized category: a
// stateless compute-on-request derivation that separates genuine clinical documents (C-CDA/HTML) from
// the flood of wearable/lifestyle "Notes" that sources such as FollowMyHealth dump into
// DocumentReference with no category. Never materialized — see pkg/document and
// docs/your-phr-dashboard/classification-and-display-architecture.md.
func GetDocumentsClassified(c *gin.Context) {
	logger := c.MustGet(pkg.ContextKeyTypeLogger).(*logrus.Entry)
	databaseRepo := c.MustGet(pkg.ContextKeyTypeDatabase).(database.DatabaseRepository)

	resources, err := databaseRepo.ListResources(c, models.ListResourceQueryOptions{SourceResourceType: "DocumentReference"})
	if err != nil {
		logger.Errorf("error listing DocumentReference resources for classification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	inputs := make([]document.InputResource, 0, len(resources))
	for i := range resources {
		inputs = append(inputs, document.InputResource{
			SourceResourceType: resources[i].SourceResourceType,
			SourceResourceID:   resources[i].SourceResourceID,
			SourceID:           resources[i].SourceID.String(),
			Raw:                json.RawMessage(resources[i].ResourceRaw),
		})
	}

	classified := document.Classify(inputs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": classified})
}
