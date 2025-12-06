package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// BranchCreateRequest represents the request payload for creating a branch
// Supports both old format (country, state, district, city as strings) and new format (country_id, state_id, etc. as integers)
type BranchCreateRequest struct {
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email,omitempty"`
	CoordinatorName string   `json:"coordinator_name,omitempty"`
	ContactNumber   string   `json:"contact_number" binding:"required"`
	EstablishedOn   *string  `json:"established_on,omitempty"`
	AashramArea     *float64 `json:"aashram_area,omitempty"`

	// Support both old format (strings) and new format (integers)
	Country    interface{} `json:"country,omitempty"`  // Can be string or number
	State      interface{} `json:"state,omitempty"`    // Can be string or number
	District   interface{} `json:"district,omitempty"` // Can be string or number
	City       interface{} `json:"city,omitempty"`     // Can be string or number
	CountryID  *uint       `json:"country_id,omitempty"`
	StateID    *uint       `json:"state_id,omitempty"`
	DistrictID *uint       `json:"district_id,omitempty"`
	CityID     *uint       `json:"city_id,omitempty"`

	Address        string `json:"address,omitempty"`
	Pincode        string `json:"pincode,omitempty"`
	PostOffice     string `json:"post_office,omitempty"`
	PoliceStation  string `json:"police_station,omitempty"`
	OpenDays       string `json:"open_days,omitempty"`
	DailyStartTime string `json:"daily_start_time,omitempty"`
	DailyEndTime   string `json:"daily_end_time,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	UpdatedBy      string `json:"updated_by,omitempty"`
	ParentBranch   interface{} `json:"parent_branch_id,omitempty"`

	Infrastructure []InfrastructureEntry `json:"infrastructure,omitempty"`
	ChildBranches  []ChildBranchEntry    `json:"child_branches,omitempty"`
	BranchMembers  []uint                `json:"branch_members,omitempty"`
}

type InfrastructureEntry struct {
	Type  string      `json:"type"`
	Count interface{} `json:"count"`
}

type ChildBranchEntry struct {
	BranchID string `json:"branchId,omitempty"`
	Address  string `json:"address,omitempty"`
}

// ToBranch converts the request to a Branch model
func (r *BranchCreateRequest) ToBranch() (*models.Branch, error) {
	branch := &models.Branch{
		Name:            r.Name,
		Email:           r.Email,
		CoordinatorName: r.CoordinatorName,
		ContactNumber:   r.ContactNumber,
		AashramArea:     0,
		Address:         r.Address,
		Pincode:         r.Pincode,
		PostOffice:      r.PostOffice,
		PoliceStation:   r.PoliceStation,
		OpenDays:        r.OpenDays,
		DailyStartTime:  r.DailyStartTime,
		DailyEndTime:    r.DailyEndTime,
		CreatedBy:       r.CreatedBy,
		UpdatedBy:       r.UpdatedBy,
	}

	if r.AashramArea != nil {
		branch.AashramArea = *r.AashramArea
	}

	// Parse EstablishedOn if provided
	if r.EstablishedOn != nil && *r.EstablishedOn != "" {
		establishedOn, err := parseTime(*r.EstablishedOn)
		if err == nil {
			branch.EstablishedOn = &establishedOn
		}
	}

	// Handle country - support both old format (string) and new format (number)
	if r.CountryID != nil {
		branch.CountryID = r.CountryID
	} else if r.Country != nil {
		countryID, err := parseID(r.Country)
		if err == nil && countryID > 0 {
			id := uint(countryID)
			branch.CountryID = &id
		}
	}

	// Handle state - support both old format (string) and new format (number)
	if r.StateID != nil {
		branch.StateID = r.StateID
	} else if r.State != nil {
		stateID, err := parseID(r.State)
		if err == nil && stateID > 0 {
			id := uint(stateID)
			branch.StateID = &id
		}
	}

	// Handle district - support both old format (string) and new format (number)
	if r.DistrictID != nil {
		branch.DistrictID = r.DistrictID
	} else if r.District != nil {
		districtID, err := parseID(r.District)
		if err == nil && districtID > 0 {
			id := uint(districtID)
			branch.DistrictID = &id
		}
	}

	// Handle city - support both old format (string) and new format (number)
	if r.CityID != nil {
		branch.CityID = r.CityID
	} else if r.City != nil {
		cityID, err := parseID(r.City)
		if err == nil && cityID > 0 {
			id := uint(cityID)
			branch.CityID = &id
		}
	}

	// Handle parent branch id if provided
	if r.ParentBranch != nil {
		if pb, err := parseID(r.ParentBranch); err == nil && pb > 0 {
			id := uint(pb)
			branch.ParentBranchID = &id
		}
	}

	return branch, nil
}

// parseID converts various types to uint
func parseID(value interface{}) (uint, error) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return 0, nil
		}
		parsed, err := strconv.ParseUint(v, 10, 32)
		return uint(parsed), err
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, nil
	}
}

// parseTime parses time string
func parseTime(timeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil // Return zero time and nil error if parsing fails (optional field)
}

// CreateBranchHandler godoc
// @Summary Create a new branch
// @Tags Branches
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param branch body BranchCreateRequest true "Branch payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/branches [post]
func CreateBranchHandler(c *gin.Context) {
	var req BranchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to Branch model
	branch, err := req.ToBranch()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Validate branch input
	if err := validators.ValidateBranchInput(branch.Name, branch.Email, branch.ContactNumber, branch.CoordinatorName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateBranch(branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Persist infrastructure entries from payload (expects canonical 'type'/'count')
	for _, infra := range req.Infrastructure {
		rt := infra.Type
		if rt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "infrastructure.type is required"})
			return
		}

		// parse count which may be number or string
		num := 0
		switch v := infra.Count.(type) {
		case string:
			if v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					num = n
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "infrastructure.count must be numeric"})
					return
				}
			}
		case float64:
			num = int(v)
		case int:
			num = v
		case nil:
			num = 0
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "infrastructure.count must be a number or numeric string"})
			return
		}

		infraModel := models.BranchInfrastructure{
			BranchID:  branch.ID,
			Type:      rt,
			Count:     num,
			CreatedBy: branch.CreatedBy,
		}
		if err := services.CreateBranchInfrastructure(&infraModel); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Link child branches by id (child_branches[].branchId)
	for _, child := range req.ChildBranches {
		if child.BranchID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "child_branches entries must include branchId to link existing branches"})
			return
		}
		cid, err := strconv.ParseUint(child.BranchID, 10, 64)
		if err != nil || cid == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branchId"})
			return
		}
		// Update child branch to set parent_branch_id to created branch
		updateData := map[string]interface{}{"parent_branch_id": branch.ID}
		if err := services.UpdateBranch(uint(cid), updateData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Link existing branch members by IDs (branch_members)
	for _, memberID := range req.BranchMembers {
		if memberID == 0 {
			continue
		}
		if err := services.UpdateBranchMember(memberID, map[string]interface{}{"branch_id": branch.ID}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

// GetBranchHandler godoc
// @Summary Get a branch by ID
// @Tags Branches
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {object} models.Branch
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/branches/{id} [get]
func GetBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	branchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	branch, err := services.GetBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branch)
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

	// Bind into a generic map so we can accept nested keys (infrastructure, child_branches, branch_members)
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract nested collections and remove them from update map before updating branch table
	infraRaw, hasInfra := payload["infrastructure"]
	if hasInfra {
		delete(payload, "infrastructure")
	}
	childRaw, hasChildren := payload["child_branches"]
	if hasChildren {
		delete(payload, "child_branches")
	}
	membersRaw, hasMembers := payload["branch_members"]
	if hasMembers {
		delete(payload, "branch_members")
	}

	// Validate remaining branch update fields
	if err := validators.ValidateBranchUpdateFields(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update branch table
	if err := services.UpdateBranch(uint(branchID), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process infrastructure: replace existing infra with provided list (if provided)
	if hasInfra {
		// Delete existing infra for branch
		if existing, err := services.GetInfrastructureByBranch(uint(branchID)); err == nil {
			for _, e := range existing {
				_ = services.DeleteBranchInfrastructure(e.ID)
			}
		}

		// Create new infra entries
		if arr, ok := infraRaw.([]interface{}); ok {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					infraType := ""
					if v, ok := m["type"]; ok {
						if s, ok := v.(string); ok {
							infraType = s
						}
					}

					number := 0
					if v, ok := m["count"]; ok {
						switch n := v.(type) {
						case string:
							if n != "" {
								if val, err := strconv.Atoi(n); err == nil {
									number = val
								}
							}
						case float64:
							number = int(n)
						}
					}

					infraModel := models.BranchInfrastructure{
						BranchID:  uint(branchID),
						Type:      infraType,
						Count:     number,
						CreatedBy: "",
					}
					_ = services.CreateBranchInfrastructure(&infraModel)
				}
			}
		}
	}

	// Process child branches: link provided branch IDs to this branch
	if hasChildren {
		if arr, ok := childRaw.([]interface{}); ok {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					if v, ok := m["branchId"]; ok {
						var cid uint64
						switch x := v.(type) {
						case string:
							cid, _ = strconv.ParseUint(x, 10, 64)
						case float64:
							cid = uint64(x)
						}
						if cid > 0 {
							_ = services.UpdateBranch(uint(cid), map[string]interface{}{"parent_branch_id": uint(branchID)})
						}
					}
				}
			}
		}
	}

	// Process branch members: link existing member IDs to this branch
	if hasMembers {
		if arr, ok := membersRaw.([]interface{}); ok {
			for _, item := range arr {
				switch v := item.(type) {
				case float64:
					mid := uint(v)
					_ = services.UpdateBranchMember(mid, map[string]interface{}{"branch_id": uint(branchID)})
				case int:
					mid := uint(v)
					_ = services.UpdateBranchMember(mid, map[string]interface{}{"branch_id": uint(branchID)})
				}
			}
		}
	}

	// Return the updated branch object (with relations preloaded)
	branch, err := services.GetBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branch)
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

	// Validate update fields
	if err := validators.ValidateBranchInfrastructureUpdateFields(updateData); err != nil {
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
