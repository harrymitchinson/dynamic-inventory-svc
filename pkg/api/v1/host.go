package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
)

// HostController is the API controller for interacting with Hosts.
type HostController struct {
	HostService inventory.HostService
}

// NewHostController initialises a HostController.
func NewHostController(h inventory.HostService) *HostController {
	return &HostController{
		HostService: h,
	}
}

// GetHosts gets the Hosts for the Environment.
func (ctrl *HostController) GetHosts(c *gin.Context) {
	env := getEnvironment(c)

	h, err := ctrl.HostService.GetHosts(env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, h)
}

// CreateHost creates a Host in the Environment.
func (ctrl *HostController) CreateHost(c *gin.Context) {
	env := getEnvironment(c)

	var h inventory.Host
	if err := c.BindJSON(&h); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := Validator.Struct(&h); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.HostService.SetHost(env, &h); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
