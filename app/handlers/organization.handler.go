package handlers

import (
	"database-ms/app/models"
	organizationSrv "database-ms/app/services/organization"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type OrganizationHandler struct {
	organization organizationSrv.OrganizationServiceInterface
}

func NewOrganizationAPI(organizationService organizationSrv.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{
		organization: organizationService,
	}
}

func (handler *OrganizationHandler) Create(c *gin.Context) {
	var newOrganization models.Organization
	c.BindJSON(&newOrganization)
	result := make(map[string]interface{})
	newOrganization.ID = bson.NewObjectId()
	err := handler.organization.Create(c.Request.Context(), &newOrganization)
	var status int
	if err == nil {
		res := &createOrganizationRes{
			ID: newOrganization.ID,
		}
		result = utils.SuccessPayload(res, "Successfully created organization")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

// func (handler *OrganizationHandler) FindByOrganizationId(c *gin.Context) {
// 	result := make(map[string]interface{})
// 	organization, err := handler.organization.FindByOrganizationId(c.Request.Context(), c.Param("organizationId"))
// 	if err == nil {
// 		result = utils.SuccessPayload(organization, "Successfully retrieved organization")
// 		utils.Response(c, http.StatusOK, result)
// 	} else {
// 		result = utils.NewHTTPError(utils.OrganizationNotFound)
// 		utils.Response(c, http.StatusBadRequest, result)
// 	}
// }

// func (handler *OrganizationHandler) Update(c *gin.Context) {
// 	var updateOrganization models.OrganizationUpdate
// 	c.BindJSON(&updateOrganization)
// 	result := make(map[string]interface{})
// 	err := handler.organization.Update(c.Request.Context(), c.Param("organizationId"), &updateOrganization)
// 	if err != nil {
// 		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
// 		utils.Response(c, http.StatusBadRequest, result)
// 		return
// 	}

// 	result = utils.SuccessPayload(nil, "Successfully updated")
// 	utils.Response(c, http.StatusOK, result)
// }

func (handler *OrganizationHandler) Delete(c *gin.Context) {
	result := make(map[string]interface{})
	err := handler.organization.Delete(c.Request.Context(), c.Param("organizationId"))
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(c, http.StatusOK, result)
}
