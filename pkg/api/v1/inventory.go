package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
)

// InventoryController is the API controller for interacting with Inventories.
type InventoryController struct {
	HostService inventory.HostService
}

// NewInventoryController initialises an InventoryController.
func NewInventoryController(h inventory.HostService) *InventoryController {
	return &InventoryController{
		HostService: h,
	}
}

// GetInventory creates an inventory file for an environment's hosts.
func (ctrl *InventoryController) GetInventory(c *gin.Context) {
	env := getEnvironment(c)

	h, err := ctrl.HostService.GetHosts(env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	file := inventory.NewInventoryFile(h)
	c.JSON(http.StatusOK, file)
}
