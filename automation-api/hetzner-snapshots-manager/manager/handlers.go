package manager

import (
	"automation-api/common/apiserver"
	"automation-api/hetzner-snapshots-manager/hetzner"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Additional handlers for api server.
func getAllRoutes(snapshots *hetzner.Snapshots) []apiserver.Route {
	return []apiserver.Route{
		{
			Path: "/hetzner/snapshots",
			Handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					server := c.Query("server")
					if server == "" {
						c.JSON(http.StatusNotFound, gin.H{
							"error":  "`server` parameter required",
							"status": "ERROR",
						})
					}
					lastSnapshotInfo, err := snapshots.GetLastForServer(server)
					if err != nil {
						if errors.Is(err, hetzner.ErrSnapshotNotFound) {
							c.JSON(http.StatusNotFound, gin.H{
								"error":  err.Error(),
								"status": "ERROR",
							})

							return
						}
						c.JSON(http.StatusInternalServerError, gin.H{
							"error":  err.Error(),
							"status": "ERROR",
						})

						return
					}

					c.JSON(http.StatusOK, gin.H{
						"body":   lastSnapshotInfo,
						"status": "OK",
					})
				}
			}(),
		},
	}
}
