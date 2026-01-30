package repository

import (
	"golang-backend/entity"
	"golang-backend/utils"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Paginate(filters map[string]interface{}, page, perPage int) (*utils.PaginationResult, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Paginate(filters map[string]interface{}, page, perPage int) (*utils.PaginationResult, error) {
	var users []entity.User
	var total int64
	query := r.db.Model(&entity.User{})

	// Smart Search
	if search, ok := filters["search"].(string); ok && search != "" {
		r.applySmartSearch(query, search)
	}

	// Filter by is_verified
	if isVerified, ok := filters["is_verified"].(string); ok && isVerified != "" {
		if isVerified == "true" {
			query.Where("is_verified = ?", true)
		} else if isVerified == "false" {
			query.Where("is_verified = ?", false)
		}
	}

	// Sorting
	sortBy := "created_at"
	if s, ok := filters["sort_by"].(string); ok && s != "" {
		sortBy = s
	}
	sortOrder := "desc"
	if s, ok := filters["sort_order"].(string); ok && s != "" {
		sortOrder = s
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply Sorting (unless using relevance from full text, usually handled in applyFullTextSearch but simplicity here)
	// If full text search was applied, it might have its own order, but GORM chain order matters.
	// For now, simple standard sort.
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	// Only apply default sort if no specific order clause exists (e.g. from full text rank)
	// Checking statement.Clauses is hard in GORM v2 without reflection or private fields sometimes.
	// We'll just append order.
	query.Order(sortBy + " " + sortOrder)

	// Pagination
	offset := (page - 1) * perPage
	err := query.Limit(perPage).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}

	pagination := utils.CalculatePagination(total, page, perPage)

	return &utils.PaginationResult{
		Items:      users,
		Pagination: pagination,
	}, nil
}

// Smart Search Logic

func (r *userRepository) applySmartSearch(query *gorm.DB, searchTerm string) {
	if r.shouldUseFullText(searchTerm) {
		r.applyFullTextSearch(query, searchTerm)
	} else {
		r.applyLikeSearch(query, searchTerm)
	}
}

func (r *userRepository) shouldUseFullText(searchTerm string) bool {
	// Simple heuristic: length >= 3 and not purely numeric
	return len(searchTerm) >= 3 && !regexp.MustCompile(`^\d+$`).MatchString(searchTerm)
}

func (r *userRepository) applyFullTextSearch(query *gorm.DB, searchTerm string) {
	tsQuery := r.preparePrefixSearchTerm(searchTerm)
	// Postgres specific
	query.Where("to_tsvector('indonesian', name || ' ' || email) @@ to_tsquery('indonesian', ?)", tsQuery)
	// Optional: Select relevance and sort by it? For now just filter.
}

func (r *userRepository) preparePrefixSearchTerm(searchTerm string) string {
	words := strings.Fields(searchTerm)
	var tsQueryParts []string
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	for _, word := range words {
		clean := re.ReplaceAllString(word, "")
		if clean != "" {
			tsQueryParts = append(tsQueryParts, clean+":*")
		}
	}
	return strings.Join(tsQueryParts, " & ")
}

func (r *userRepository) applyLikeSearch(query *gorm.DB, searchTerm string) {
	// Fallback for short terms or non-Postgres (though query above assumes Postgres for fulltext)
	search := "%" + searchTerm + "%"
	query.Where("name ILIKE ? OR email ILIKE ?", search, search)
}
