package apidoc

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/pkg/model"
	"strings"
)

// Repository interface defines all database operations for API documentation
type Repository interface {
	// Collection operations
	CreateCollection(collection *Collection) error
	GetCollectionByID(id int64, companyID int64) (*Collection, error)
	GetCollections(companyID int64, limit, offset int) ([]*Collection, int64, error)
	GetCollectionWithStats(id int64, companyID int64) (*CollectionWithStats, error)
	UpdateCollection(collection *Collection) error
	DeleteCollection(id int64, companyID int64) error

	// Folder operations
	CreateFolder(folder *Folder) error
	GetFoldersByCollectionID(collectionID int64) ([]*Folder, error)
	GetFolderByID(id int64) (*Folder, error)
	GetFoldersHierarchy(collectionID int64) ([]*FolderWithChildren, error)
	UpdateFolder(folder *Folder) error
	DeleteFolder(id int64) error

	// Endpoint operations
	CreateEndpoint(endpoint *Endpoint) error
	GetEndpointsByCollectionID(collectionID int64, limit, offset int) ([]*Endpoint, int64, error)
	GetEndpointByID(id int64) (*Endpoint, error)
	GetEndpointWithDetails(id int64) (*EndpointWithDetails, error)
	UpdateEndpoint(endpoint *Endpoint) error
	DeleteEndpoint(id int64) error

	// Environment operations
	CreateEnvironment(environment *Environment) error
	GetEnvironmentsByCollectionID(collectionID int64) ([]*Environment, error)
	GetEnvironmentByID(id int64) (*Environment, error)
	GetEnvironmentWithVariables(id int64) (*EnvironmentWithVariables, error)
	UpdateEnvironment(environment *Environment) error
	DeleteEnvironment(id int64) error

	// Environment Variable operations
	CreateEnvironmentVariable(variable *EnvironmentVariable) error
	GetEnvironmentVariables(environmentID int64) ([]*EnvironmentVariable, error)
	GetEnvironmentVariableByID(id int64) (*EnvironmentVariable, error)
	UpdateEnvironmentVariable(variable *EnvironmentVariable) error
	DeleteEnvironmentVariable(id int64) error

	// Header operations
	CreateHeader(header *Header) error
	GetHeadersByEndpointID(endpointID int64) ([]*Header, error)
	UpdateHeader(header *Header) error
	DeleteHeader(id int64) error

	// Parameter operations
	CreateParameter(parameter *Parameter) error
	GetParametersByEndpointID(endpointID int64) ([]*Parameter, error)
	UpdateParameter(parameter *Parameter) error
	DeleteParameter(id int64) error

	// Request Body operations
	CreateRequestBody(requestBody *RequestBody) error
	GetRequestBodyByEndpointID(endpointID int64) (*RequestBody, error)
	UpdateRequestBody(requestBody *RequestBody) error
	DeleteRequestBody(id int64) error

	// Response operations
	CreateResponse(response *Response) error
	GetResponsesByEndpointID(endpointID int64) ([]*Response, error)
	UpdateResponse(response *Response) error
	DeleteResponse(id int64) error

	// Tag operations
	GetAllTags() ([]*Tag, error)
	CreateTag(tag *Tag) error
	GetTagsByEndpointID(endpointID int64) ([]*Tag, error)
	AddTagToEndpoint(endpointID, tagID int64) error
	RemoveTagFromEndpoint(endpointID, tagID int64) error

	// Export operations
	GetCollectionForExport(collectionID int64, companyID int64) (*CollectionExport, error)
}

// repository implements Repository interface using raw SQL
type repository struct {
	*model.Repository
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) Repository {
	return &repository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// Collection operations

func (r *repository) CreateCollection(collection *Collection) error {
	query := `
		INSERT INTO api_collections (name, description, version, base_url, schema_version, created_by, company_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		collection.Name,
		collection.Description,
		collection.Version,
		collection.BaseURL,
		collection.SchemaVersion,
		collection.CreatedBy,
		collection.CompanyID,
		collection.IsActive,
	).Scan(&collection.ID, &collection.CreatedAt, &collection.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

func (r *repository) GetCollectionByID(id int64, companyID int64) (*Collection, error) {
	collection := &Collection{}
	query := `
		SELECT id, name, description, version, base_url, schema_version, created_by, company_id, is_active, created_at, updated_at
		FROM api_collections
		WHERE id = $1 AND company_id = $2 AND is_active = true
	`

	err := r.db.QueryRow(query, id, companyID).Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.Version,
		&collection.BaseURL,
		&collection.SchemaVersion,
		&collection.CreatedBy,
		&collection.CompanyID,
		&collection.IsActive,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return collection, nil
}

func (r *repository) GetCollections(companyID int64, limit, offset int) ([]*Collection, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM api_collections WHERE company_id = $1 AND is_active = true`
	err := r.db.QueryRow(countQuery, companyID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collections count: %w", err)
	}

	// Get collections
	query := `
		SELECT id, name, description, version, base_url, schema_version, created_by, company_id, is_active, created_at, updated_at
		FROM api_collections
		WHERE company_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, companyID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get collections: %w", err)
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		collection := &Collection{}
		err := rows.Scan(
			&collection.ID,
			&collection.Name,
			&collection.Description,
			&collection.Version,
			&collection.BaseURL,
			&collection.SchemaVersion,
			&collection.CreatedBy,
			&collection.CompanyID,
			&collection.IsActive,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, collection)
	}

	return collections, total, nil
}

func (r *repository) GetCollectionWithStats(id int64, companyID int64) (*CollectionWithStats, error) {
	stats := &CollectionWithStats{}
	query := `
		SELECT 
			c.id, c.name, c.description, c.version, c.base_url, c.schema_version, 
			c.created_by, c.company_id, c.is_active, c.created_at, c.updated_at,
			COALESCE(s.total_folders, 0) as total_folders,
			COALESCE(s.total_endpoints, 0) as total_endpoints,
			COALESCE(s.total_environments, 0) as total_environments,
			COALESCE(s.get_endpoints, 0) as get_endpoints,
			COALESCE(s.post_endpoints, 0) as post_endpoints,
			COALESCE(s.put_endpoints, 0) as put_endpoints,
			COALESCE(s.delete_endpoints, 0) as delete_endpoints
		FROM api_collections c
		LEFT JOIN api_collection_stats s ON c.id = s.id
		WHERE c.id = $1 AND c.company_id = $2 AND c.is_active = true
	`

	err := r.db.QueryRow(query, id, companyID).Scan(
		&stats.ID,
		&stats.Name,
		&stats.Description,
		&stats.Version,
		&stats.BaseURL,
		&stats.SchemaVersion,
		&stats.CreatedBy,
		&stats.CompanyID,
		&stats.IsActive,
		&stats.CreatedAt,
		&stats.UpdatedAt,
		&stats.TotalFolders,
		&stats.TotalEndpoints,
		&stats.TotalEnvironments,
		&stats.GetEndpoints,
		&stats.PostEndpoints,
		&stats.PutEndpoints,
		&stats.DeleteEndpoints,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, fmt.Errorf("failed to get collection with stats: %w", err)
	}

	return stats, nil
}

func (r *repository) UpdateCollection(collection *Collection) error {
	query := `
		UPDATE api_collections 
		SET name = $1, description = $2, version = $3, base_url = $4, schema_version = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND company_id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		collection.Name,
		collection.Description,
		collection.Version,
		collection.BaseURL,
		collection.SchemaVersion,
		collection.IsActive,
		collection.ID,
		collection.CompanyID,
	).Scan(&collection.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("collection not found or access denied")
		}
		return fmt.Errorf("failed to update collection: %w", err)
	}

	return nil
}

func (r *repository) DeleteCollection(id int64, companyID int64) error {
	query := `DELETE FROM api_collections WHERE id = $1 AND company_id = $2`

	result, err := r.db.Exec(query, id, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found or access denied")
	}

	return nil
}

// Folder operations

func (r *repository) CreateFolder(folder *Folder) error {
	query := `
		INSERT INTO api_folders (collection_id, parent_id, name, description, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		folder.CollectionID,
		folder.ParentID,
		folder.Name,
		folder.Description,
		folder.SortOrder,
	).Scan(&folder.ID, &folder.CreatedAt, &folder.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	return nil
}

func (r *repository) GetFoldersByCollectionID(collectionID int64) ([]*Folder, error) {
	query := `
		SELECT id, collection_id, parent_id, name, description, sort_order, created_at, updated_at
		FROM api_folders
		WHERE collection_id = $1
		ORDER BY sort_order, name
	`

	rows, err := r.db.Query(query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}
	defer rows.Close()

	var folders []*Folder
	for rows.Next() {
		folder := &Folder{}
		err := rows.Scan(
			&folder.ID,
			&folder.CollectionID,
			&folder.ParentID,
			&folder.Name,
			&folder.Description,
			&folder.SortOrder,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *repository) GetFolderByID(id int64) (*Folder, error) {
	folder := &Folder{}
	query := `
		SELECT id, collection_id, parent_id, name, description, sort_order, created_at, updated_at
		FROM api_folders
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&folder.ID,
		&folder.CollectionID,
		&folder.ParentID,
		&folder.Name,
		&folder.Description,
		&folder.SortOrder,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return folder, nil
}

func (r *repository) GetFoldersHierarchy(collectionID int64) ([]*FolderWithChildren, error) {
	folders, err := r.GetFoldersByCollectionID(collectionID)
	if err != nil {
		return nil, err
	}

	// Build hierarchy
	folderMap := make(map[int64]*FolderWithChildren)
	var rootFolders []*FolderWithChildren

	// First pass: create all folder nodes
	for _, folder := range folders {
		folderWithChildren := &FolderWithChildren{
			Folder:   *folder,
			Children: make([]*FolderWithChildren, 0),
		}
		folderMap[folder.ID] = folderWithChildren
	}

	// Second pass: build hierarchy
	for _, folder := range folders {
		if folder.ParentID.Valid {
			// This is a child folder
			if parent, exists := folderMap[folder.ParentID.Int64]; exists {
				parent.Children = append(parent.Children, folderMap[folder.ID])
			}
		} else {
			// This is a root folder
			rootFolders = append(rootFolders, folderMap[folder.ID])
		}
	}

	return rootFolders, nil
}

func (r *repository) UpdateFolder(folder *Folder) error {
	query := `
		UPDATE api_folders 
		SET parent_id = $1, name = $2, description = $3, sort_order = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		folder.ParentID,
		folder.Name,
		folder.Description,
		folder.SortOrder,
		folder.ID,
	).Scan(&folder.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("folder not found")
		}
		return fmt.Errorf("failed to update folder: %w", err)
	}

	return nil
}

func (r *repository) DeleteFolder(id int64) error {
	query := `DELETE FROM api_folders WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("folder not found")
	}

	return nil
}

// Endpoint operations

func (r *repository) CreateEndpoint(endpoint *Endpoint) error {
	query := `
		INSERT INTO api_endpoints (collection_id, folder_id, name, description, method, url, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		endpoint.CollectionID,
		endpoint.FolderID,
		endpoint.Name,
		endpoint.Description,
		endpoint.Method,
		endpoint.URL,
		endpoint.SortOrder,
		endpoint.IsActive,
	).Scan(&endpoint.ID, &endpoint.CreatedAt, &endpoint.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	return nil
}

func (r *repository) GetEndpointsByCollectionID(collectionID int64, limit, offset int) ([]*Endpoint, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM api_endpoints WHERE collection_id = $1 AND is_active = true`
	err := r.db.QueryRow(countQuery, collectionID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get endpoints count: %w", err)
	}

	// Get endpoints
	query := `
		SELECT id, collection_id, folder_id, name, description, method, url, sort_order, is_active, created_at, updated_at
		FROM api_endpoints
		WHERE collection_id = $1 AND is_active = true
		ORDER BY sort_order, name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, collectionID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get endpoints: %w", err)
	}
	defer rows.Close()

	var endpoints []*Endpoint
	for rows.Next() {
		endpoint := &Endpoint{}
		err := rows.Scan(
			&endpoint.ID,
			&endpoint.CollectionID,
			&endpoint.FolderID,
			&endpoint.Name,
			&endpoint.Description,
			&endpoint.Method,
			&endpoint.URL,
			&endpoint.SortOrder,
			&endpoint.IsActive,
			&endpoint.CreatedAt,
			&endpoint.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan endpoint: %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, total, nil
}

func (r *repository) GetEndpointByID(id int64) (*Endpoint, error) {
	endpoint := &Endpoint{}
	query := `
		SELECT id, collection_id, folder_id, name, description, method, url, sort_order, is_active, created_at, updated_at
		FROM api_endpoints
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&endpoint.ID,
		&endpoint.CollectionID,
		&endpoint.FolderID,
		&endpoint.Name,
		&endpoint.Description,
		&endpoint.Method,
		&endpoint.URL,
		&endpoint.SortOrder,
		&endpoint.IsActive,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("endpoint not found")
		}
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	return endpoint, nil
}

func (r *repository) GetEndpointWithDetails(id int64) (*EndpointWithDetails, error) {
	// Get the endpoint first
	endpoint, err := r.GetEndpointByID(id)
	if err != nil {
		return nil, err
	}

	details := &EndpointWithDetails{
		Endpoint: *endpoint,
	}

	// Get headers
	details.Headers, err = r.GetHeadersByEndpointID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get headers: %w", err)
	}

	// Get parameters
	details.Parameters, err = r.GetParametersByEndpointID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Get request body
	details.RequestBody, err = r.GetRequestBodyByEndpointID(id)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, fmt.Errorf("failed to get request body: %w", err)
	}

	// Get responses
	details.Responses, err = r.GetResponsesByEndpointID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	// Get tags
	details.Tags, err = r.GetTagsByEndpointID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return details, nil
}

func (r *repository) UpdateEndpoint(endpoint *Endpoint) error {
	query := `
		UPDATE api_endpoints 
		SET folder_id = $1, name = $2, description = $3, method = $4, url = $5, sort_order = $6, is_active = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		endpoint.FolderID,
		endpoint.Name,
		endpoint.Description,
		endpoint.Method,
		endpoint.URL,
		endpoint.SortOrder,
		endpoint.IsActive,
		endpoint.ID,
	).Scan(&endpoint.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("endpoint not found")
		}
		return fmt.Errorf("failed to update endpoint: %w", err)
	}

	return nil
}

func (r *repository) DeleteEndpoint(id int64) error {
	query := `DELETE FROM api_endpoints WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("endpoint not found")
	}

	return nil
}

// Environment operations

func (r *repository) CreateEnvironment(environment *Environment) error {
	query := `
		INSERT INTO api_environments (collection_id, name, description, is_default)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		environment.CollectionID,
		environment.Name,
		environment.Description,
		environment.IsDefault,
	).Scan(&environment.ID, &environment.CreatedAt, &environment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	return nil
}

func (r *repository) GetEnvironmentsByCollectionID(collectionID int64) ([]*Environment, error) {
	query := `
		SELECT id, collection_id, name, description, is_default, created_at, updated_at
		FROM api_environments
		WHERE collection_id = $1
		ORDER BY is_default DESC, name
	`

	rows, err := r.db.Query(query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments: %w", err)
	}
	defer rows.Close()

	var environments []*Environment
	for rows.Next() {
		environment := &Environment{}
		err := rows.Scan(
			&environment.ID,
			&environment.CollectionID,
			&environment.Name,
			&environment.Description,
			&environment.IsDefault,
			&environment.CreatedAt,
			&environment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan environment: %w", err)
		}
		environments = append(environments, environment)
	}

	return environments, nil
}

func (r *repository) GetEnvironmentByID(id int64) (*Environment, error) {
	environment := &Environment{}
	query := `
		SELECT id, collection_id, name, description, is_default, created_at, updated_at
		FROM api_environments
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&environment.ID,
		&environment.CollectionID,
		&environment.Name,
		&environment.Description,
		&environment.IsDefault,
		&environment.CreatedAt,
		&environment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("environment not found")
		}
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	return environment, nil
}

func (r *repository) GetEnvironmentWithVariables(id int64) (*EnvironmentWithVariables, error) {
	environment, err := r.GetEnvironmentByID(id)
	if err != nil {
		return nil, err
	}

	variables, err := r.GetEnvironmentVariables(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %w", err)
	}

	return &EnvironmentWithVariables{
		Environment: *environment,
		Variables:   variables,
	}, nil
}

func (r *repository) UpdateEnvironment(environment *Environment) error {
	query := `
		UPDATE api_environments 
		SET name = $1, description = $2, is_default = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		environment.Name,
		environment.Description,
		environment.IsDefault,
		environment.ID,
	).Scan(&environment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("environment not found")
		}
		return fmt.Errorf("failed to update environment: %w", err)
	}

	return nil
}

func (r *repository) DeleteEnvironment(id int64) error {
	query := `DELETE FROM api_environments WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("environment not found")
	}

	return nil
}

// Environment Variable operations

func (r *repository) CreateEnvironmentVariable(variable *EnvironmentVariable) error {
	query := `
		INSERT INTO api_environment_variables (environment_id, key_name, value, description, is_secret)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		variable.EnvironmentID,
		variable.KeyName,
		variable.Value,
		variable.Description,
		variable.IsSecret,
	).Scan(&variable.ID, &variable.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create environment variable: %w", err)
	}

	return nil
}

func (r *repository) GetEnvironmentVariables(environmentID int64) ([]*EnvironmentVariable, error) {
	query := `
		SELECT id, environment_id, key_name, value, description, is_secret, created_at
		FROM api_environment_variables
		WHERE environment_id = $1
		ORDER BY key_name
	`

	rows, err := r.db.Query(query, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %w", err)
	}
	defer rows.Close()

	var variables []*EnvironmentVariable
	for rows.Next() {
		variable := &EnvironmentVariable{}
		err := rows.Scan(
			&variable.ID,
			&variable.EnvironmentID,
			&variable.KeyName,
			&variable.Value,
			&variable.Description,
			&variable.IsSecret,
			&variable.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan environment variable: %w", err)
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

func (r *repository) GetEnvironmentVariableByID(id int64) (*EnvironmentVariable, error) {
	variable := &EnvironmentVariable{}
	query := `
		SELECT id, environment_id, key_name, value, description, is_secret, created_at
		FROM api_environment_variables
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&variable.ID,
		&variable.EnvironmentID,
		&variable.KeyName,
		&variable.Value,
		&variable.Description,
		&variable.IsSecret,
		&variable.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("environment variable not found")
		}
		return nil, fmt.Errorf("failed to get environment variable: %w", err)
	}

	return variable, nil
}

func (r *repository) UpdateEnvironmentVariable(variable *EnvironmentVariable) error {
	query := `
		UPDATE api_environment_variables 
		SET key_name = $1, value = $2, description = $3, is_secret = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		variable.KeyName,
		variable.Value,
		variable.Description,
		variable.IsSecret,
		variable.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update environment variable: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("environment variable not found")
	}

	return nil
}

func (r *repository) DeleteEnvironmentVariable(id int64) error {
	query := `DELETE FROM api_environment_variables WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("environment variable not found")
	}

	return nil
}

// Header operations

func (r *repository) CreateHeader(header *Header) error {
	query := `
		INSERT INTO api_headers (endpoint_id, key_name, value, description, is_required, header_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		header.EndpointID,
		header.KeyName,
		header.Value,
		header.Description,
		header.IsRequired,
		header.HeaderType,
	).Scan(&header.ID, &header.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create header: %w", err)
	}

	return nil
}

func (r *repository) GetHeadersByEndpointID(endpointID int64) ([]*Header, error) {
	query := `
		SELECT id, endpoint_id, key_name, value, description, is_required, header_type, created_at
		FROM api_headers
		WHERE endpoint_id = $1
		ORDER BY header_type, key_name
	`

	rows, err := r.db.Query(query, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get headers: %w", err)
	}
	defer rows.Close()

	var headers []*Header
	for rows.Next() {
		header := &Header{}
		err := rows.Scan(
			&header.ID,
			&header.EndpointID,
			&header.KeyName,
			&header.Value,
			&header.Description,
			&header.IsRequired,
			&header.HeaderType,
			&header.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan header: %w", err)
		}
		headers = append(headers, header)
	}

	return headers, nil
}

func (r *repository) UpdateHeader(header *Header) error {
	query := `
		UPDATE api_headers 
		SET key_name = $1, value = $2, description = $3, is_required = $4, header_type = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		query,
		header.KeyName,
		header.Value,
		header.Description,
		header.IsRequired,
		header.HeaderType,
		header.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update header: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("header not found")
	}

	return nil
}

func (r *repository) DeleteHeader(id int64) error {
	query := `DELETE FROM api_headers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete header: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("header not found")
	}

	return nil
}

// Parameter operations

func (r *repository) CreateParameter(parameter *Parameter) error {
	query := `
		INSERT INTO api_parameters (endpoint_id, name, type, data_type, description, default_value, example_value, is_required)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		parameter.EndpointID,
		parameter.Name,
		parameter.Type,
		parameter.DataType,
		parameter.Description,
		parameter.DefaultValue,
		parameter.ExampleValue,
		parameter.IsRequired,
	).Scan(&parameter.ID, &parameter.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create parameter: %w", err)
	}

	return nil
}

func (r *repository) GetParametersByEndpointID(endpointID int64) ([]*Parameter, error) {
	query := `
		SELECT id, endpoint_id, name, type, data_type, description, default_value, example_value, is_required, created_at
		FROM api_parameters
		WHERE endpoint_id = $1
		ORDER BY type, name
	`

	rows, err := r.db.Query(query, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}
	defer rows.Close()

	var parameters []*Parameter
	for rows.Next() {
		parameter := &Parameter{}
		err := rows.Scan(
			&parameter.ID,
			&parameter.EndpointID,
			&parameter.Name,
			&parameter.Type,
			&parameter.DataType,
			&parameter.Description,
			&parameter.DefaultValue,
			&parameter.ExampleValue,
			&parameter.IsRequired,
			&parameter.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parameter: %w", err)
		}
		parameters = append(parameters, parameter)
	}

	return parameters, nil
}

func (r *repository) UpdateParameter(parameter *Parameter) error {
	query := `
		UPDATE api_parameters 
		SET name = $1, type = $2, data_type = $3, description = $4, default_value = $5, example_value = $6, is_required = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(
		query,
		parameter.Name,
		parameter.Type,
		parameter.DataType,
		parameter.Description,
		parameter.DefaultValue,
		parameter.ExampleValue,
		parameter.IsRequired,
		parameter.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update parameter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parameter not found")
	}

	return nil
}

func (r *repository) DeleteParameter(id int64) error {
	query := `DELETE FROM api_parameters WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete parameter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parameter not found")
	}

	return nil
}

// Request Body operations

func (r *repository) CreateRequestBody(requestBody *RequestBody) error {
	query := `
		INSERT INTO api_request_bodies (endpoint_id, content_type, body_content, description, schema_definition)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		requestBody.EndpointID,
		requestBody.ContentType,
		requestBody.BodyContent,
		requestBody.Description,
		requestBody.SchemaDefinition,
	).Scan(&requestBody.ID, &requestBody.CreatedAt, &requestBody.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create request body: %w", err)
	}

	return nil
}

func (r *repository) GetRequestBodyByEndpointID(endpointID int64) (*RequestBody, error) {
	requestBody := &RequestBody{}
	query := `
		SELECT id, endpoint_id, content_type, body_content, description, schema_definition, created_at, updated_at
		FROM api_request_bodies
		WHERE endpoint_id = $1
	`

	err := r.db.QueryRow(query, endpointID).Scan(
		&requestBody.ID,
		&requestBody.EndpointID,
		&requestBody.ContentType,
		&requestBody.BodyContent,
		&requestBody.Description,
		&requestBody.SchemaDefinition,
		&requestBody.CreatedAt,
		&requestBody.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("request body not found")
		}
		return nil, fmt.Errorf("failed to get request body: %w", err)
	}

	return requestBody, nil
}

func (r *repository) UpdateRequestBody(requestBody *RequestBody) error {
	query := `
		UPDATE api_request_bodies 
		SET content_type = $1, body_content = $2, description = $3, schema_definition = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		requestBody.ContentType,
		requestBody.BodyContent,
		requestBody.Description,
		requestBody.SchemaDefinition,
		requestBody.ID,
	).Scan(&requestBody.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("request body not found")
		}
		return fmt.Errorf("failed to update request body: %w", err)
	}

	return nil
}

func (r *repository) DeleteRequestBody(id int64) error {
	query := `DELETE FROM api_request_bodies WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete request body: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("request body not found")
	}

	return nil
}

// Response operations

func (r *repository) CreateResponse(response *Response) error {
	query := `
		INSERT INTO api_responses (endpoint_id, status_code, status_text, content_type, response_body, description, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		response.EndpointID,
		response.StatusCode,
		response.StatusText,
		response.ContentType,
		response.ResponseBody,
		response.Description,
		response.IsDefault,
	).Scan(&response.ID, &response.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	return nil
}

func (r *repository) GetResponsesByEndpointID(endpointID int64) ([]*Response, error) {
	query := `
		SELECT id, endpoint_id, status_code, status_text, content_type, response_body, description, is_default, created_at
		FROM api_responses
		WHERE endpoint_id = $1
		ORDER BY status_code
	`

	rows, err := r.db.Query(query, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}
	defer rows.Close()

	var responses []*Response
	for rows.Next() {
		response := &Response{}
		err := rows.Scan(
			&response.ID,
			&response.EndpointID,
			&response.StatusCode,
			&response.StatusText,
			&response.ContentType,
			&response.ResponseBody,
			&response.Description,
			&response.IsDefault,
			&response.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response: %w", err)
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (r *repository) UpdateResponse(response *Response) error {
	query := `
		UPDATE api_responses 
		SET status_code = $1, status_text = $2, content_type = $3, response_body = $4, description = $5, is_default = $6
		WHERE id = $7
	`

	result, err := r.db.Exec(
		query,
		response.StatusCode,
		response.StatusText,
		response.ContentType,
		response.ResponseBody,
		response.Description,
		response.IsDefault,
		response.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update response: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("response not found")
	}

	return nil
}

func (r *repository) DeleteResponse(id int64) error {
	query := `DELETE FROM api_responses WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete response: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("response not found")
	}

	return nil
}

// Tag operations

func (r *repository) GetAllTags() ([]*Tag, error) {
	query := `
		SELECT id, name, color, description, created_at
		FROM api_tags
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := &Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *repository) CreateTag(tag *Tag) error {
	query := `
		INSERT INTO api_tags (name, color, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		tag.Name,
		tag.Color,
		tag.Description,
	).Scan(&tag.ID, &tag.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

func (r *repository) GetTagsByEndpointID(endpointID int64) ([]*Tag, error) {
	query := `
		SELECT t.id, t.name, t.color, t.description, t.created_at
		FROM api_tags t
		JOIN api_endpoint_tags et ON t.id = et.tag_id
		WHERE et.endpoint_id = $1
		ORDER BY t.name
	`

	rows, err := r.db.Query(query, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := &Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *repository) AddTagToEndpoint(endpointID, tagID int64) error {
	query := `
		INSERT INTO api_endpoint_tags (endpoint_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT (endpoint_id, tag_id) DO NOTHING
	`

	_, err := r.db.Exec(query, endpointID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to endpoint: %w", err)
	}

	return nil
}

func (r *repository) RemoveTagFromEndpoint(endpointID, tagID int64) error {
	query := `DELETE FROM api_endpoint_tags WHERE endpoint_id = $1 AND tag_id = $2`

	result, err := r.db.Exec(query, endpointID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from endpoint: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag relationship not found")
	}

	return nil
}

// Export operations

func (r *repository) GetCollectionForExport(collectionID int64, companyID int64) (*CollectionExport, error) {
	// Get collection
	collection, err := r.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, err
	}

	export := &CollectionExport{
		Collection: *collection,
	}

	// Get folders hierarchy
	export.Folders, err = r.GetFoldersHierarchy(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	// Get all endpoints with details
	endpoints, _, err := r.GetEndpointsByCollectionID(collectionID, 1000, 0) // Get all endpoints
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	for _, endpoint := range endpoints {
		details, err := r.GetEndpointWithDetails(endpoint.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get endpoint details: %w", err)
		}
		export.Endpoints = append(export.Endpoints, details)
	}

	// Get environments with variables
	environments, err := r.GetEnvironmentsByCollectionID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments: %w", err)
	}

	for _, env := range environments {
		envWithVars, err := r.GetEnvironmentWithVariables(env.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get environment variables: %w", err)
		}
		export.Environments = append(export.Environments, envWithVars)
	}

	// Get all tags
	export.Tags, err = r.GetAllTags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return export, nil
}
