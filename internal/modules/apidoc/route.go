package apidoc

import (
	"database/sql"
	"fmt"
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

// Collection handlers

func (h *Handler) GetCollections(c *gin.Context) {
	var req CollectionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	userID := c.GetInt64("current_user_id")
	result, err := h.service.GetCollectionsByUser(&req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get collections", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCollectionsRetrieved, result)
}

func (h *Handler) GetCollectionByID(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetCollectionByIDByUser(collectionID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgCollectionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCollectionRetrieved, result)
}

func (h *Handler) GetCollectionWithStats(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetCollectionWithStatsByUser(collectionID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgCollectionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCollectionRetrieved, result)
}

func (h *Handler) CreateCollection(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateCollectionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	userID := c.GetInt64("current_user_id")
	result, err := h.service.CreateCollectionByUser(req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create collection", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgCollectionCreated, result)
}

func (h *Handler) UpdateCollection(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateCollectionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateCollectionByUser(collectionID, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCollectionUpdated, result)
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	if err := h.service.DeleteCollectionByUser(collectionID, userID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCollectionDeleted, nil)
}

// Folder handlers

func (h *Handler) GetFolders(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetFoldersByUser(collectionID, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get folders", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgFoldersRetrieved, result)
}

func (h *Handler) GetFolderByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid folder ID")
		return
	}

	userID := c.GetInt64("current_user_id")
	result, err := h.service.GetFolderByIDByUser(id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgFolderNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgFolderRetrieved, result)
}

func (h *Handler) CreateFolder(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateFolderRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateFolderByUser(collectionID, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create folder", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgFolderCreated, result)
}

func (h *Handler) UpdateFolder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid folder ID")
		return
	}

	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateFolderRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateFolderByUser(id, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgFolderUpdated, result)
}

func (h *Handler) DeleteFolder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid folder ID")
		return
	}

	userID := c.GetInt64("current_user_id")

	if err := h.service.DeleteFolderByUser(id, userID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgFolderDeleted, nil)
}

func (h *Handler) ReorderFolders(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	var folderOrders []FolderOrderRequest
	if err := c.ShouldBindJSON(&folderOrders); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.service.ReorderFoldersByUser(collectionID, userID, folderOrders); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Folders reordered successfully", nil)
}

// Endpoint handlers

func (h *Handler) GetEndpoints(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	var req EndpointListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetEndpointsByUser(collectionID, userID, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get endpoints", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEndpointsRetrieved, result)
}

func (h *Handler) GetEndpointByID(c *gin.Context) {
	endpointID := c.GetInt64("endpoint_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEndpointByIDByUser(endpointID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgEndpointNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEndpointRetrieved, result)
}

func (h *Handler) GetEndpointWithDetails(c *gin.Context) {
	endpointID := c.GetInt64("endpoint_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEndpointWithDetailsByUser(endpointID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgEndpointNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEndpointRetrieved, result)
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateEndpointRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateEndpointByUser(collectionID, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create endpoint", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgEndpointCreated, result)
}

func (h *Handler) UpdateEndpoint(c *gin.Context) {
	endpointID := c.GetInt64("endpoint_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateEndpointRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateEndpointByUser(endpointID, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEndpointUpdated, result)
}

func (h *Handler) DeleteEndpoint(c *gin.Context) {
	endpointID := c.GetInt64("endpoint_id")
	userID := c.GetInt64("current_user_id")

	if err := h.service.DeleteEndpointByUser(endpointID, userID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEndpointDeleted, nil)
}

// Bulk endpoint operations

func (h *Handler) BulkCreateEndpoints(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*BulkCreateEndpointsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.BulkCreateEndpointsByUser(collectionID, userID, req.Endpoints)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create endpoints", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Endpoints created successfully", result)
}

func (h *Handler) BulkUpdateEndpoints(c *gin.Context) {
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*BulkUpdateEndpointsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.BulkUpdateEndpointsByUser(userID, req.Updates)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update endpoints", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Endpoints updated successfully", result)
}

func (h *Handler) BulkDeleteEndpoints(c *gin.Context) {
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*BulkDeleteEndpointsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.BulkDeleteEndpointsByUser(userID, req.EndpointIDs); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to delete endpoints", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Endpoints deleted successfully", nil)
}

func (h *Handler) BulkMoveEndpoints(c *gin.Context) {
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*BulkMoveEndpointsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.BulkMoveEndpointsByUser(userID, req.EndpointIDs, req.TargetFolderID); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to move endpoints", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Endpoints moved successfully", nil)
}

// Environment handlers

func (h *Handler) GetEnvironments(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEnvironmentsByUser(collectionID, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get environments", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEnvironmentsRetrieved, result)
}

func (h *Handler) GetEnvironmentByID(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEnvironmentByIDByUser(environmentID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgEnvironmentNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEnvironmentRetrieved, result)
}

func (h *Handler) GetEnvironmentWithVariables(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEnvironmentWithVariablesByUser(environmentID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgEnvironmentNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEnvironmentRetrieved, result)
}

func (h *Handler) CreateEnvironment(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateEnvironmentRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateEnvironmentByUser(collectionID, userID, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create environment", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgEnvironmentCreated, result)
}

func (h *Handler) UpdateEnvironment(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateEnvironmentRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateEnvironmentByUser(environmentID, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEnvironmentUpdated, result)
}

func (h *Handler) DeleteEnvironment(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	if err := h.service.DeleteEnvironmentByUser(environmentID, userID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgEnvironmentDeleted, nil)
}

// Environment Variable handlers

func (h *Handler) GetEnvironmentVariables(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	result, err := h.service.GetEnvironmentVariablesByUser(environmentID, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get environment variables", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgVariablesRetrieved, result)
}

func (h *Handler) CreateEnvironmentVariable(c *gin.Context) {
	environmentID := c.GetInt64("environment_id")
	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateEnvironmentVariableRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateEnvironmentVariableByUser(environmentID, userID, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create environment variable", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgVariableCreated, result)
}

func (h *Handler) UpdateEnvironmentVariable(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("variable_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid variable ID")
		return
	}

	userID := c.GetInt64("current_user_id")

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateEnvironmentVariableRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateEnvironmentVariableByUser(id, req, userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgVariableUpdated, result)
}

func (h *Handler) DeleteEnvironmentVariable(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("variable_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid variable ID")
		return
	}

	userID := c.GetInt64("current_user_id")

	if err := h.service.DeleteEnvironmentVariableByUser(id, userID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgVariableDeleted, nil)
}

// Export handlers

func (h *Handler) ExportCollection(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	var req ExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	// Validate format
	if req.Format == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "format parameter is required")
		return
	}

	// Perform export
	result, err := h.service.ExportCollectionByUser(collectionID, userID, req.Format, req.EnvironmentID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Export failed", err.Error())
		return
	}

	// Set appropriate headers for file download
	filename := result.Filename
	contentType := result.ContentType

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", result.Size))

	// Return the content
	if contentType == "application/x-yaml" {
		c.String(http.StatusOK, result.Content.(string))
	} else {
		c.String(http.StatusOK, result.Content.(string))
	}
}

func (h *Handler) GetSupportedExportFormats(c *gin.Context) {
	formats := h.service.GetSupportedExportFormats()

	response.Success(c, http.StatusOK, "Supported export formats retrieved", gin.H{
		"formats": formats,
	})
}

// Cache and Performance monitoring handlers

func (h *Handler) GetExportCacheStats(c *gin.Context) {
	stats, err := h.service.GetExportCacheStats()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get cache stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Cache statistics retrieved", stats)
}

func (h *Handler) GetExportPerformanceStats(c *gin.Context) {
	stats := h.service.GetExportPerformanceStats()

	response.Success(c, http.StatusOK, "Performance statistics retrieved", stats)
}

func (h *Handler) GetExportProgress(c *gin.Context) {
	progressID := c.Param("progress_id")
	if progressID == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "Progress ID is required")
		return
	}

	progress, exists := h.service.GetExportProgress(progressID)
	if !exists {
		response.Error(c, http.StatusNotFound, "Progress not found", "Export progress not found")
		return
	}

	response.Success(c, http.StatusOK, "Export progress retrieved", progress)
}

func (h *Handler) GetAllExportProgress(c *gin.Context) {
	progress := h.service.GetAllExportProgress()

	response.Success(c, http.StatusOK, "All export progress retrieved", progress)
}

func (h *Handler) InvalidateExportCache(c *gin.Context) {
	collectionID := c.GetInt64("collection_id")
	userID := c.GetInt64("current_user_id")

	// Check permissions
	if err := h.service.checkExportPermission(userID, "write"); err != nil {
		response.ErrorWithAutoStatus(c, "Permission denied", err.Error())
		return
	}

	// Check resource ownership
	if err := h.service.checkResourceOwnership(userID, collectionID); err != nil {
		response.ErrorWithAutoStatus(c, "Access denied", err.Error())
		return
	}

	if err := h.service.InvalidateExportCache(collectionID); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to invalidate cache", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Export cache invalidated successfully", nil)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler, db *sql.DB) {
	// API Documentation routes
	apiDocs := api.Group("/api-docs")
	apiDocs.Use(middleware.APIDocAuthMiddleware(db))
	apiDocs.Use(middleware.RequireResourceOwnership())
	{
		// Collection routes
		collections := apiDocs.Group("/collections")
		collections.Use(middleware.APIDocCollectionPermission("read")) // Base read permission for collections
		{
			// GET /api/v1/api-docs/collections - Get all collections
			collections.GET("", handler.GetCollections)

			// POST /api/v1/api-docs/collections - Create new collection
			collections.POST("",
				middleware.APIDocCollectionPermission("write"),
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &CreateCollectionRequest{},
				}),
				handler.CreateCollection,
			)

			// Collection-specific routes
			collection := collections.Group("/:collection_id")
			collection.Use(middleware.ValidateCollectionID())
			{
				// GET /api/v1/api-docs/collections/:id - Get collection by ID
				collection.GET("", handler.GetCollectionByID)

				// GET /api/v1/api-docs/collections/:id/stats - Get collection with stats
				collection.GET("/stats", handler.GetCollectionWithStats)

				// PUT /api/v1/api-docs/collections/:id - Update collection
				collection.PUT("",
					middleware.APIDocCollectionPermission("write"),
					middleware.ValidateRequest(middleware.ValidationRules{
						Body: &UpdateCollectionRequest{},
					}),
					handler.UpdateCollection,
				)

				// DELETE /api/v1/api-docs/collections/:id - Delete collection
				collection.DELETE("",
					middleware.APIDocCollectionPermission("delete"),
					handler.DeleteCollection,
				)

				// Folder routes within collection
				folders := collection.Group("/folders")
				folders.Use(middleware.APIDocCollectionPermission("read")) // Base read permission
				{
					// GET /api/v1/api-docs/collections/:id/folders - Get folders
					folders.GET("", handler.GetFolders)

					// POST /api/v1/api-docs/collections/:id/folders - Create folder
					folders.POST("",
						middleware.APIDocCollectionPermission("write"),
						middleware.ValidateRequest(middleware.ValidationRules{
							Body: &CreateFolderRequest{},
						}),
						handler.CreateFolder,
					)

					// POST /api/v1/api-docs/collections/:id/folders/reorder - Reorder folders
					folders.POST("/reorder",
						middleware.APIDocCollectionPermission("write"),
						handler.ReorderFolders,
					)
				}

				// Endpoint routes within collection
				endpoints := collection.Group("/endpoints")
				endpoints.Use(middleware.APIDocEndpointPermission("read")) // Base read permission
				{
					// GET /api/v1/api-docs/collections/:id/endpoints - Get endpoints
					endpoints.GET("", handler.GetEndpoints)

					// POST /api/v1/api-docs/collections/:id/endpoints - Create endpoint
					endpoints.POST("",
						middleware.APIDocEndpointPermission("write"),
						middleware.ValidateRequest(middleware.ValidationRules{
							Body: &CreateEndpointRequest{},
						}),
						handler.CreateEndpoint,
					)

					// POST /api/v1/api-docs/collections/:id/endpoints/bulk - Bulk create endpoints
					endpoints.POST("/bulk",
						middleware.APIDocEndpointPermission("write"),
						middleware.ValidateRequest(middleware.ValidationRules{
							Body: &BulkCreateEndpointsRequest{},
						}),
						handler.BulkCreateEndpoints,
					)
				}

				// Environment routes within collection
				environments := collection.Group("/environments")
				environments.Use(middleware.APIDocEnvironmentPermission("read")) // Base read permission
				{
					// GET /api/v1/api-docs/collections/:id/environments - Get environments
					environments.GET("", handler.GetEnvironments)

					// POST /api/v1/api-docs/collections/:id/environments - Create environment
					environments.POST("",
						middleware.APIDocEnvironmentPermission("write"),
						middleware.ValidateRequest(middleware.ValidationRules{
							Body: &CreateEnvironmentRequest{},
						}),
						handler.CreateEnvironment,
					)
				}

				// Export routes
				exports := collection.Group("/export")
				exports.Use(middleware.APIDocExportPermission("read"))
				{
					// GET /api/v1/api-docs/collections/:id/export - Export collection
					exports.GET("", handler.ExportCollection)

					// GET /api/v1/api-docs/collections/:id/export/formats - Get supported formats
					exports.GET("/formats", handler.GetSupportedExportFormats)

					// Cache management routes
					exports.POST("/cache/invalidate",
						middleware.APIDocExportPermission("write"),
						handler.InvalidateExportCache,
					)
				}
			}
		}

		// Export monitoring routes (admin only)
		monitoring := apiDocs.Group("/monitoring")
		monitoring.Use(middleware.APIDocCollectionPermission("read"))
		{
			// GET /api/v1/api-docs/monitoring/cache/stats - Get cache statistics
			monitoring.GET("/cache/stats", handler.GetExportCacheStats)

			// GET /api/v1/api-docs/monitoring/performance/stats - Get performance statistics
			monitoring.GET("/performance/stats", handler.GetExportPerformanceStats)

			// GET /api/v1/api-docs/monitoring/progress - Get all export progress
			monitoring.GET("/progress", handler.GetAllExportProgress)

			// GET /api/v1/api-docs/monitoring/progress/:progress_id - Get specific export progress
			monitoring.GET("/progress/:progress_id", handler.GetExportProgress)
		}

		// Individual folder routes
		folders := apiDocs.Group("/folders")
		folders.Use(middleware.APIDocCollectionPermission("read"))
		{
			// GET /api/v1/api-docs/folders/:id - Get folder by ID
			folders.GET("/:id", handler.GetFolderByID)

			// PUT /api/v1/api-docs/folders/:id - Update folder
			folders.PUT("/:id",
				middleware.APIDocCollectionPermission("write"),
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &UpdateFolderRequest{},
				}),
				handler.UpdateFolder,
			)

			// DELETE /api/v1/api-docs/folders/:id - Delete folder
			folders.DELETE("/:id",
				middleware.APIDocCollectionPermission("delete"),
				handler.DeleteFolder,
			)
		}

		// Individual endpoint routes
		endpoints := apiDocs.Group("/endpoints")
		endpoints.Use(middleware.APIDocEndpointPermission("read"))
		{
			// GET /api/v1/api-docs/endpoints/:id - Get endpoint by ID
			endpoints.GET("/:endpoint_id",
				middleware.ValidateEndpointID(),
				handler.GetEndpointByID,
			)

			// GET /api/v1/api-docs/endpoints/:id/details - Get endpoint with details
			endpoints.GET("/:endpoint_id/details",
				middleware.ValidateEndpointID(),
				handler.GetEndpointWithDetails,
			)

			// PUT /api/v1/api-docs/endpoints/:id - Update endpoint
			endpoints.PUT("/:endpoint_id",
				middleware.APIDocEndpointPermission("write"),
				middleware.ValidateEndpointID(),
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &UpdateEndpointRequest{},
				}),
				handler.UpdateEndpoint,
			)

			// DELETE /api/v1/api-docs/endpoints/:id - Delete endpoint
			endpoints.DELETE("/:endpoint_id",
				middleware.APIDocEndpointPermission("delete"),
				middleware.ValidateEndpointID(),
				handler.DeleteEndpoint,
			)
		}

		// Bulk endpoint operations
		bulkEndpoints := apiDocs.Group("/endpoints/bulk")
		bulkEndpoints.Use(middleware.APIDocEndpointPermission("write"))
		{
			// PUT /api/v1/api-docs/endpoints/bulk - Bulk update endpoints
			bulkEndpoints.PUT("",
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &BulkUpdateEndpointsRequest{},
				}),
				handler.BulkUpdateEndpoints,
			)

			// DELETE /api/v1/api-docs/endpoints/bulk - Bulk delete endpoints
			bulkEndpoints.DELETE("",
				middleware.APIDocEndpointPermission("delete"),
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &BulkDeleteEndpointsRequest{},
				}),
				handler.BulkDeleteEndpoints,
			)

			// POST /api/v1/api-docs/endpoints/bulk/move - Bulk move endpoints
			bulkEndpoints.POST("/move",
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &BulkMoveEndpointsRequest{},
				}),
				handler.BulkMoveEndpoints,
			)
		}

		// Individual environment routes
		environments := apiDocs.Group("/environments")
		environments.Use(middleware.APIDocEnvironmentPermission("read"))
		{
			// GET /api/v1/api-docs/environments/:id - Get environment by ID
			environments.GET("/:environment_id",
				middleware.ValidateEnvironmentID(),
				handler.GetEnvironmentByID,
			)

			// GET /api/v1/api-docs/environments/:id/with-variables - Get environment with variables
			environments.GET("/:environment_id/with-variables",
				middleware.ValidateEnvironmentID(),
				handler.GetEnvironmentWithVariables,
			)

			// PUT /api/v1/api-docs/environments/:id - Update environment
			environments.PUT("/:environment_id",
				middleware.APIDocEnvironmentPermission("write"),
				middleware.ValidateEnvironmentID(),
				middleware.ValidateRequest(middleware.ValidationRules{
					Body: &UpdateEnvironmentRequest{},
				}),
				handler.UpdateEnvironment,
			)

			// DELETE /api/v1/api-docs/environments/:id - Delete environment
			environments.DELETE("/:environment_id",
				middleware.APIDocEnvironmentPermission("delete"),
				middleware.ValidateEnvironmentID(),
				handler.DeleteEnvironment,
			)

			// Environment variable routes
			envVars := environments.Group("/:environment_id/variables")
			envVars.Use(middleware.ValidateEnvironmentID())
			{
				// GET /api/v1/api-docs/environments/:id/variables - Get variables
				envVars.GET("", handler.GetEnvironmentVariables)

				// POST /api/v1/api-docs/environments/:id/variables - Create variable
				envVars.POST("",
					middleware.APIDocEnvironmentPermission("write"),
					middleware.ValidateRequest(middleware.ValidationRules{
						Body: &CreateEnvironmentVariableRequest{},
					}),
					handler.CreateEnvironmentVariable,
				)

				// PUT /api/v1/api-docs/environments/:id/variables/:variable_id - Update variable
				envVars.PUT("/:variable_id",
					middleware.APIDocEnvironmentPermission("write"),
					middleware.ValidateRequest(middleware.ValidationRules{
						Body: &UpdateEnvironmentVariableRequest{},
					}),
					handler.UpdateEnvironmentVariable,
				)

				// DELETE /api/v1/api-docs/environments/:id/variables/:variable_id - Delete variable
				envVars.DELETE("/:variable_id",
					middleware.APIDocEnvironmentPermission("delete"),
					handler.DeleteEnvironmentVariable,
				)
			}
		}
	}
}
