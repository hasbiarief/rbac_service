package docs

// Package docs berisi anotasi Swagger untuk modul [MODULE_NAME].
// File ini terpisah dari route.go untuk menjaga kebersihan kode.
//
// INSTRUKSI PENGGUNAAN:
// 1. Salin file ini ke internal/modules/[module_name]/docs/swagger.go
// 2. Ganti [MODULE_NAME] dengan nama modul Anda
// 3. Ganti [ModuleTag] dengan tag yang sesuai (contoh: Users, Companies, Products)
// 4. Ganti [Entity] dengan nama entitas Anda (contoh: User, Company, Product)
// 5. Sesuaikan path, parameter, dan response sesuai kebutuhan
// 6. Hapus komentar instruksi ini setelah selesai
// 7. Jalankan: make swagger-gen

// ============================================================================
// CONTOH 1: GET - Mendapatkan daftar entitas dengan pagination
// ============================================================================

// @Summary      Get list of [entities]
// @Description  Retrieves a paginated list of [entities] with optional filters
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        limit     query     int     false  "Items per page"  default(10)
// @Param        search    query     string  false  "Search term"
// @Param        sort      query     string  false  "Sort field"  default(created_at)
// @Param        order     query     string  false  "Sort order"  Enums(asc, desc)  default(desc)
// @Success      200       {object}  response.Response{data=[]module.[Entity]}  "List retrieved successfully"
// @Failure      400       {object}  response.Response  "Invalid query parameters"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name] [get]

// ============================================================================
// CONTOH 2: GET - Mendapatkan entitas berdasarkan ID
// ============================================================================

// @Summary      Get [entity] by ID
// @Description  Retrieves a single [entity] by its unique identifier
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "[Entity] ID"
// @Success      200  {object}  response.Response{data=module.[Entity]}  "[Entity] found"
// @Failure      400  {object}  response.Response  "Invalid ID format"
// @Failure      404  {object}  response.Response  "[Entity] not found"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id} [get]

// ============================================================================
// CONTOH 3: POST - Membuat entitas baru
// ============================================================================

// @Summary      Create new [entity]
// @Description  Creates a new [entity] with the provided information
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        [entity]  body      module.Create[Entity]Request  true  "[Entity] information"
// @Success      201       {object}  response.Response{data=module.[Entity]}  "[Entity] created successfully"
// @Failure      400       {object}  response.Response  "Invalid request body or validation failed"
// @Failure      409       {object}  response.Response  "[Entity] already exists"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name] [post]

// ============================================================================
// CONTOH 4: PUT - Memperbarui entitas
// ============================================================================

// @Summary      Update [entity]
// @Description  Updates an existing [entity] with new information
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id        path      int                           true  "[Entity] ID"
// @Param        [entity]  body      module.Update[Entity]Request  true  "Updated [entity] information"
// @Success      200       {object}  response.Response{data=module.[Entity]}  "[Entity] updated successfully"
// @Failure      400       {object}  response.Response  "Invalid request body or validation failed"
// @Failure      404       {object}  response.Response  "[Entity] not found"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id} [put]

// ============================================================================
// CONTOH 5: PATCH - Memperbarui sebagian entitas
// ============================================================================

// @Summary      Partially update [entity]
// @Description  Updates specific fields of an existing [entity]
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id        path      int                          true  "[Entity] ID"
// @Param        [entity]  body      module.Patch[Entity]Request  true  "Fields to update"
// @Success      200       {object}  response.Response{data=module.[Entity]}  "[Entity] updated successfully"
// @Failure      400       {object}  response.Response  "Invalid request body"
// @Failure      404       {object}  response.Response  "[Entity] not found"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id} [patch]

// ============================================================================
// CONTOH 6: DELETE - Menghapus entitas
// ============================================================================

// @Summary      Delete [entity]
// @Description  Deletes an [entity] by its ID
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "[Entity] ID"
// @Success      200  {object}  response.Response  "[Entity] deleted successfully"
// @Failure      400  {object}  response.Response  "Invalid ID format"
// @Failure      404  {object}  response.Response  "[Entity] not found"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id} [delete]

// ============================================================================
// CONTOH 7: POST - Endpoint dengan multiple parameters
// ============================================================================

// @Summary      Search [entities]
// @Description  Advanced search for [entities] with multiple criteria
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        query     body      module.Search[Entity]Request  true  "Search criteria"
// @Param        page      query     int                           false "Page number"  default(1)
// @Param        limit     query     int                           false "Items per page"  default(10)
// @Success      200       {object}  response.Response{data=module.Search[Entity]Response}  "Search results"
// @Failure      400       {object}  response.Response  "Invalid search criteria"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/search [post]

// ============================================================================
// CONTOH 8: GET - Endpoint dengan relasi (nested resource)
// ============================================================================

// @Summary      Get [entity] details with relations
// @Description  Retrieves [entity] with related data (e.g., user, company)
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "[Entity] ID"
// @Success      200  {object}  response.Response{data=module.[Entity]WithDetails}  "[Entity] with details"
// @Failure      404  {object}  response.Response  "[Entity] not found"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id}/details [get]

// ============================================================================
// CONTOH 9: POST - File upload
// ============================================================================

// @Summary      Upload [entity] file
// @Description  Uploads a file associated with the [entity]
// @Tags         [ModuleTag]
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      int   true  "[Entity] ID"
// @Param        file  formData  file  true  "File to upload"
// @Success      200   {object}  response.Response{data=module.FileUploadResponse}  "File uploaded successfully"
// @Failure      400   {object}  response.Response  "Invalid file or file too large"
// @Failure      404   {object}  response.Response  "[Entity] not found"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id}/upload [post]

// ============================================================================
// CONTOH 10: GET - Export data
// ============================================================================

// @Summary      Export [entities]
// @Description  Exports [entities] data in specified format
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      application/octet-stream
// @Param        format  query  string  false  "Export format"  Enums(csv, xlsx, pdf)  default(csv)
// @Param        filter  query  string  false  "Filter criteria"
// @Success      200     {file}  binary  "[Entities] exported successfully"
// @Failure      400     {object}  response.Response  "Invalid export format"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/export [get]

// ============================================================================
// CONTOH 11: POST - Bulk operations
// ============================================================================

// @Summary      Bulk create [entities]
// @Description  Creates multiple [entities] in a single request
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        [entities]  body      module.BulkCreate[Entity]Request  true  "List of [entities] to create"
// @Success      201         {object}  response.Response{data=module.BulkCreate[Entity]Response}  "[Entities] created successfully"
// @Failure      400         {object}  response.Response  "Invalid request or validation failed"
// @Failure      500         {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/bulk [post]

// ============================================================================
// CONTOH 12: GET - Endpoint tanpa autentikasi (public)
// ============================================================================

// @Summary      Get public [entity] information
// @Description  Retrieves publicly available [entity] information (no authentication required)
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "[Entity] ID"
// @Success      200  {object}  response.Response{data=module.Public[Entity]}  "Public information retrieved"
// @Failure      404  {object}  response.Response  "[Entity] not found"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/[module_name]/{id}/public [get]

// ============================================================================
// CONTOH 13: GET - Endpoint dengan enum parameter
// ============================================================================

// @Summary      Get [entities] by status
// @Description  Retrieves [entities] filtered by status
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        status  path      string  true  "[Entity] status"  Enums(active, inactive, pending, archived)
// @Success      200     {object}  response.Response{data=[]module.[Entity]}  "[Entities] retrieved"
// @Failure      400     {object}  response.Response  "Invalid status value"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/status/{status} [get]

// ============================================================================
// CONTOH 14: POST - Endpoint dengan custom header
// ============================================================================

// @Summary      Process [entity] with custom options
// @Description  Processes [entity] with options specified in custom header
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        id              path      int                         true  "[Entity] ID"
// @Param        X-Process-Mode  header    string                      false "Processing mode"  Enums(sync, async)  default(sync)
// @Param        options         body      module.Process[Entity]Request true  "Processing options"
// @Success      200             {object}  response.Response{data=module.Process[Entity]Response}  "Processing completed"
// @Failure      400             {object}  response.Response  "Invalid request"
// @Failure      404             {object}  response.Response  "[Entity] not found"
// @Failure      500             {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/{id}/process [post]

// ============================================================================
// CONTOH 15: GET - Endpoint dengan multiple response types
// ============================================================================

// @Summary      Get [entity] statistics
// @Description  Retrieves statistical information about [entities]
// @Tags         [ModuleTag]
// @Accept       json
// @Produce      json
// @Param        period  query     string  false  "Time period"  Enums(daily, weekly, monthly, yearly)  default(monthly)
// @Param        from    query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string  false  "End date (YYYY-MM-DD)"
// @Success      200     {object}  response.Response{data=module.[Entity]Statistics}  "Statistics retrieved"
// @Failure      400     {object}  response.Response  "Invalid date format or period"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/[module_name]/statistics [get]

// ============================================================================
// TIPS DAN BEST PRACTICES
// ============================================================================
//
// 1. GUNAKAN TAG YANG KONSISTEN
//    - Gunakan PascalCase untuk tag (Users, Companies, Products)
//    - Kelompokkan endpoint yang terkait dengan tag yang sama
//
// 2. DESKRIPSI YANG JELAS
//    - @Summary: Deskripsi singkat (1 baris)
//    - @Description: Deskripsi detail tentang fungsi endpoint
//
// 3. DOKUMENTASIKAN SEMUA PARAMETER
//    - Sertakan tipe data yang tepat
//    - Tandai parameter required/optional dengan benar
//    - Berikan deskripsi yang jelas untuk setiap parameter
//
// 4. DOKUMENTASIKAN SEMUA RESPONSE
//    - Sertakan success response (200, 201, dll)
//    - Sertakan semua kemungkinan error response (400, 401, 404, 500, dll)
//    - Berikan pesan yang deskriptif untuk setiap response
//
// 5. GUNAKAN TYPE REFERENCE YANG TEPAT
//    - Gunakan nama struct yang lengkap: module.EntityName
//    - Untuk response dengan data: response.Response{data=module.Entity}
//    - Untuk array: []module.Entity
//
// 6. SECURITY ANNOTATION
//    - Tambahkan @Security BearerAuth untuk endpoint yang memerlukan autentikasi
//    - Jangan tambahkan untuk endpoint public
//
// 7. PARAMETER TYPES
//    - path: Parameter di URL path (/users/{id})
//    - query: Query string parameter (?page=1&limit=10)
//    - body: Request body (JSON)
//    - header: HTTP header
//    - formData: Form data (untuk file upload)
//
// 8. HTTP METHODS
//    - GET: Mengambil data
//    - POST: Membuat data baru
//    - PUT: Update data lengkap
//    - PATCH: Update sebagian data
//    - DELETE: Menghapus data
//
// 9. STATUS CODES
//    - 200: Success (GET, PUT, PATCH, DELETE)
//    - 201: Created (POST)
//    - 400: Bad Request (validasi gagal)
//    - 401: Unauthorized (autentikasi gagal)
//    - 403: Forbidden (tidak ada akses)
//    - 404: Not Found (resource tidak ditemukan)
//    - 409: Conflict (duplikasi data)
//    - 500: Internal Server Error
//
// 10. REGENERASI DOKUMENTASI
//     - Setelah menambah/mengubah anotasi, jalankan: make swagger-gen
//     - Untuk validasi: make swagger-validate
//     - Untuk watch mode: make swagger-watch
//
// ============================================================================
// REFERENSI
// ============================================================================
//
// - Swagger Guide: docs/SWAGGER_GUIDE.md
// - Swaggo Documentation: https://github.com/swaggo/swag
// - OpenAPI Specification: https://swagger.io/specification/
// - Contoh implementasi: internal/modules/auth/docs/swagger.go
//
// ============================================================================
