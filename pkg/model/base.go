package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// JSONB type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return json.Unmarshal(bytes, j)
}

// NullString handles nullable string fields
type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// NullInt64 handles nullable int64 fields
type NullInt64 struct {
	sql.NullInt64
}

func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// NullTime handles nullable time fields
type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		nt.Time = *t
	} else {
		nt.Valid = false
	}
	return nil
}

// Repository provides common database operations
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetTableName extracts table name from struct tag or struct name
func (r *Repository) GetTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Convert struct name to snake_case
	name := t.Name()
	return toSnakeCase(name)
}

// GetColumns extracts column names from struct tags
func (r *Repository) GetColumns(model interface{}) []string {
	var columns []string
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("db"); tag != "" && tag != "-" {
			columns = append(columns, tag)
		}
	}

	return columns
}

// BuildInsertQuery builds INSERT query for a model
func (r *Repository) BuildInsertQuery(model interface{}) (string, []interface{}) {
	tableName := r.GetTableName(model)
	columns := r.GetColumns(model)

	// Remove id, created_at, updated_at from insert columns
	var insertColumns []string
	var placeholders []string
	var values []interface{}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for _, col := range columns {
		if col == "id" || col == "created_at" || col == "updated_at" {
			continue
		}

		insertColumns = append(insertColumns, col)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))

		// Get field value
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			if field.Tag.Get("db") == col {
				values = append(values, v.Field(j).Interface())
				break
			}
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id, created_at, updated_at",
		tableName,
		strings.Join(insertColumns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values
}

// BuildUpdateQuery builds UPDATE query for a model
func (r *Repository) BuildUpdateQuery(model interface{}, id int64) (string, []interface{}) {
	tableName := r.GetTableName(model)
	columns := r.GetColumns(model)

	var setClauses []string
	var values []interface{}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for _, col := range columns {
		if col == "id" || col == "created_at" {
			continue
		}

		if col == "updated_at" {
			setClauses = append(setClauses, fmt.Sprintf("%s = CURRENT_TIMESTAMP", col))
			continue
		}

		// Get field value
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			if field.Tag.Get("db") == col {
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, len(values)+1))
				values = append(values, v.Field(j).Interface())
				break
			}
		}
	}

	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d RETURNING updated_at",
		tableName,
		strings.Join(setClauses, ", "),
		len(values),
	)

	return query, values
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
