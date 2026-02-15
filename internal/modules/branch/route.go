package branch

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Handler methods

// @Summary      Get all branches
// @Description  Mendapatkan daftar semua branch dengan filter opsional dan pagination. Mendukung nested hierarchy dengan query parameter nested=true
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        limit       query     int     false  "Limit jumlah data"
// @Param        offset      query     int     false  "Offset data"
// @Param        search      query     string  false  "Search by name atau code"
// @Param        company_id  query     int     false  "Filter by company ID"
// @Param        is_active   query     bool    false  "Filter by active status"
// @Param        nested      query     string  false  "Return nested hierarchy (true/false)"
// @Success      200         {object}  response.Response{data=branch.BranchListResponse}  "Branches berhasil diambil"
// @Failure      400         {object}  response.Response  "Bad request"
// @Failure      500         {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches [get]
// @Security     BearerAuth
func (h *Handler) GetBranches(c *gin.Context) {
	var req BranchListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameter", err.Error())
		return
	}

	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.service.GetBranchesNested(&req)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.service.GetBranches(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved, result)
}

// @Summary      Get branch by ID
// @Description  Mendapatkan detail branch berdasarkan ID
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Branch ID"
// @Success      200  {object}  response.Response{data=branch.BranchResponse}  "Branch berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid branch ID"
// @Failure      404  {object}  response.Response  "Branch tidak ditemukan"
// @Router       /api/v1/branches/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetBranchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.service.GetBranchByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgBranchNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchRetrieved, result)
}

// @Summary      Create new branch
// @Description  Membuat branch baru dengan support untuk hierarchical structure
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        branch  body      branch.CreateBranchRequest  true  "Branch data"
// @Success      201     {object}  response.Response{data=branch.BranchResponse}  "Branch berhasil dibuat"
// @Failure      400     {object}  response.Response  "Bad request - validation failed"
// @Failure      409     {object}  response.Response  "Conflict - branch code sudah ada"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches [post]
// @Security     BearerAuth
func (h *Handler) CreateBranch(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.service.CreateBranch(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create branch", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgBranchCreated, result)
}

// @Summary      Update branch
// @Description  Memperbarui informasi branch
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        id      path      int                         true  "Branch ID"
// @Param        branch  body      branch.UpdateBranchRequest  true  "Branch data yang akan diupdate"
// @Success      200     {object}  response.Response{data=branch.BranchResponse}  "Branch berhasil diupdate"
// @Failure      400     {object}  response.Response  "Bad request - Invalid branch ID atau validation failed"
// @Failure      404     {object}  response.Response  "Branch tidak ditemukan"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.service.UpdateBranch(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchUpdated, result)
}

// @Summary      Delete branch
// @Description  Menghapus branch berdasarkan ID
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Branch ID"
// @Success      200  {object}  response.Response  "Branch berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid branch ID"
// @Failure      404  {object}  response.Response  "Branch tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	if err := h.service.DeleteBranch(id); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to delete branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchDeleted, nil)
}

// @Summary      Get company branches
// @Description  Mendapatkan semua branches dari company tertentu. Mendukung nested hierarchy dengan query parameter nested=true
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        companyId          path      int     true   "Company ID"
// @Param        include_hierarchy  query     string  false  "Include hierarchy information (true/false)"
// @Param        nested             query     string  false  "Return nested hierarchy (true/false)"
// @Success      200                {object}  response.Response{data=[]branch.BranchResponse}  "Company branches berhasil diambil"
// @Failure      400                {object}  response.Response  "Bad request - Invalid company ID"
// @Failure      404                {object}  response.Response  "Company tidak ditemukan"
// @Failure      500                {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/company/{companyId} [get]
// @Security     BearerAuth
func (h *Handler) GetCompanyBranches(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Company ID is invalid", "Company ID must be a valid number")
		return
	}

	includeHierarchy := c.DefaultQuery("include_hierarchy", "true") == "true"
	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.service.GetCompanyBranchesNested(companyID)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.service.GetCompanyBranches(companyID, includeHierarchy)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved, result)
}

// @Summary      Get branch children
// @Description  Mendapatkan daftar child branches dari branch tertentu
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Branch ID"
// @Success      200  {object}  response.Response{data=[]branch.BranchResponse}  "Branch children berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid branch ID"
// @Failure      404  {object}  response.Response  "Branch tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/{id}/children [get]
// @Security     BearerAuth
func (h *Handler) GetBranchChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.service.GetBranchChildren(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branch children", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchChildrenRetrieved, result)
}

// @Summary      Get branch hierarchy
// @Description  Mendapatkan hierarchy branch berdasarkan ID. Mendukung nested format dengan query parameter nested=true
// @Tags         Branches
// @Accept       json
// @Produce      json
// @Param        id      path      int     true   "Branch ID"
// @Param        nested  query     string  false  "Return nested hierarchy (true/false)"
// @Success      200     {object}  response.Response{data=branch.BranchResponse}  "Branch hierarchy berhasil diambil"
// @Failure      400     {object}  response.Response  "Bad request - Invalid branch ID"
// @Failure      404     {object}  response.Response  "Branch tidak ditemukan"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/{id}/hierarchy [get]
// @Security     BearerAuth
func (h *Handler) GetBranchHierarchy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	nested := c.DefaultQuery("nested", "false") == "true"

	result, err := h.service.GetBranchHierarchyByID(id, nested)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branch hierarchy", err.Error())
		return
	}

	message := constants.MsgBranchHierarchyRetrieved
	if nested {
		message += " (nested)"
	}

	response.Success(c, http.StatusOK, message, result)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	branches := api.Group("/branches")
	{
		// GET /api/v1/branches - Get all branches with optional filters
		branches.GET("", handler.GetBranches)

		// GET /api/v1/branches/:id - Get branch by ID
		branches.GET("/:id", handler.GetBranchByID)

		// GET /api/v1/branches/:id/hierarchy - Get branch hierarchy by ID
		branches.GET("/:id/hierarchy", handler.GetBranchHierarchy)

		// GET /api/v1/branches/:id/children - Get branch children by ID
		branches.GET("/:id/children", handler.GetBranchChildren)

		// POST /api/v1/branches - Create new branch
		branches.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateBranchRequest{},
			}),
			handler.CreateBranch,
		)

		// PUT /api/v1/branches/:id - Update branch by ID
		branches.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateBranchRequest{},
			}),
			handler.UpdateBranch,
		)

		// DELETE /api/v1/branches/:id - Delete branch by ID
		branches.DELETE("/:id", handler.DeleteBranch)

		// GET /api/v1/branches/company/:companyId - Get branches by company ID
		branches.GET("/company/:companyId", handler.GetCompanyBranches)
	}
}
