package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateBranchHandler godoc
// @Summary Create a new branch
// @Tags Branches
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param branch body models.Branch true "Branch payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branches [post]
func CreateBranchHandler(c *gin.Context) {
	var branch models.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate branch input
	if err := validators.ValidateBranchInput(branch.Name, branch.Email, branch.ContactNumber, branch.CoordinatorName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateBranch(&branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Branch created successfully",
		"branch":  branch,
	})
}

// GetAllBranchesHandler godoc
// @Summary Get all branches
// @Tags Branches
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Branch
// @Failure 500 {object} map[string]string
// @Router /api/branches [get]
func GetAllBranchesHandler(c *gin.Context) {
	branches, err := services.GetAllBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branches)
}

// GetBranchSearchHandler godoc
// @Summary Get branches by name or coordinator (or all if none provided)
// @Description Retrieve branches by name and/or coordinator name, or list all if no filters.
// @Tags Branches
// @Security ApiKeyAuth
// @Produce json
// @Param name query string false "Branch Name"
// @Param coordinator query string false "Coordinator Name"
// @Success 200 {array} models.Branch
// @Failure 404 {object} map[string]string
// @Router /api/branches/search [get]
func GetBranchSearchHandler(c *gin.Context) {
	name := c.Query("name")
	coordinator := c.Query("coordinator")

	branches, err := services.GetBranchSearch(name, coordinator)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// UpdateBranchHandler godoc
// @Summary Update a branch
// @Tags Branches
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Param branch body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branches/{id} [put]
func UpdateBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	branchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate update fields
	if err := validators.ValidateBranchUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateBranch(uint(branchID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch updated successfully"})
}

// DeleteBranchHandler godoc
// @Summary Delete a branch
// @Tags Branches
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branches/{id} [delete]
func DeleteBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	branchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	if err := services.DeleteBranch(uint(branchID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
}

// *************************************** Branch Infrastructure ****************************************************** //

// CreateBranchInfrastructureHandler godoc
// @Summary Create a new branch infrastructure record
// @Description Adds a new infrastructure entry for a branch
// @Tags BranchInfrastructure
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param infra body models.BranchInfrastructure true "Infrastructure payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-infra [post]
func CreateBranchInfrastructureHandler(c *gin.Context) {
	var infra models.BranchInfrastructure
	if err := c.ShouldBindJSON(&infra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate infrastructure input
	if err := validators.ValidateBranchInfrastructure(infra.BranchID, infra.Type, infra.Count); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateBranchInfrastructure(&infra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Infrastructure created successfully",
		"data":    infra,
	})
}

// GetAllBranchInfrastructureHandler godoc
// @Summary Get all branch infrastructure records
// @Description Fetch all branch infrastructure entries
// @Tags BranchInfrastructure
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.BranchInfrastructure
// @Failure 500 {object} map[string]string
// @Router /api/branch-infra [get]
func GetAllBranchInfrastructureHandler(c *gin.Context) {
	infra, err := services.GetAllBranchInfrastructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, infra)
}

// GetInfrastructureByBranchHandler godoc
// @Summary Get infrastructure records for a branch
// @Description Fetch all infrastructure entries for a specific branch
// @Tags BranchInfrastructure
// @Security ApiKeyAuth
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {array} models.BranchInfrastructure
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/branch-infra/branch/{branch_id} [get]
func GetInfrastructureByBranchHandler(c *gin.Context) {
	branchIDParam := c.Param("branch_id")
	branchID, err := strconv.ParseUint(branchIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch ID"})
		return
	}

	infra, err := services.GetInfrastructureByBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infra)
}

// UpdateBranchInfrastructureHandler godoc
// @Summary Update a branch infrastructure record
// @Description Update existing infrastructure entry by ID
// @Tags BranchInfrastructure
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Infrastructure ID"
// @Param infra body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-infra/{id} [put]
func UpdateBranchInfrastructureHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateBranchInfrastructure(uint(id), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Infrastructure updated successfully"})
}

// DeleteBranchInfrastructureHandler godoc
// @Summary Delete a branch infrastructure record
// @Description Delete an infrastructure entry by ID
// @Tags BranchInfrastructure
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Infrastructure ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-infra/{id} [delete]
func DeleteBranchInfrastructureHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := services.DeleteBranchInfrastructure(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Infrastructure deleted successfully"})
}

// *************************************** Branch Infrastructure ****************************************************** //

// CreateBranchMemberHandler godoc
// @Summary Create a new branch member
// @Tags BranchMember
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param member body models.BranchMember true "Branch Member payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-member [post]
func CreateBranchMemberHandler(c *gin.Context) {
	var member models.BranchMember
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate branch member input
	if err := validators.ValidateBranchMember(member.Name, member.MemberType, member.BranchID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateBranchMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Branch member created successfully",
		"data":    member,
	})
}

// GetAllBranchMembersHandler godoc
// @Summary Get all branch members
// @Tags BranchMember
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.BranchMember
// @Failure 500 {object} map[string]string
// @Router /api/branch-member [get]
func GetAllBranchMembersHandler(c *gin.Context) {
	members, err := services.GetAllBranchMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

// GetMembersByBranchHandler godoc
// @Summary Get members of a branch
// @Tags BranchMember
// @Security ApiKeyAuth
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {array} models.BranchMember
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/branch-member/branch/{branch_id} [get]
func GetMembersByBranchHandler(c *gin.Context) {
	branchIDParam := c.Param("branch_id")
	branchID, err := strconv.ParseUint(branchIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch ID"})
		return
	}

	members, err := services.GetMembersByBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// UpdateBranchMemberHandler godoc
// @Summary Update a branch member
// @Tags BranchMember
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Param member body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-member/{id} [put]
func UpdateBranchMemberHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate update fields
	if err := validators.ValidateBranchMemberUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateBranchMember(uint(id), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch member updated successfully"})
}

// DeleteBranchMemberHandler godoc
// @Summary Delete a branch member
// @Tags BranchMember
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branch-member/{id} [delete]
func DeleteBranchMemberHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := services.DeleteBranchMember(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch member deleted successfully"})
}
