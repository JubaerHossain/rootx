package utilQuery

import (
	"bytes"

	// "crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	coreEntity "github.com/JubaerHossain/rootx/pkg/core/entity"
	"github.com/JubaerHossain/rootx/pkg/utils"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func Pagination(req *http.Request, app *app.App, table string) (pagination coreEntity.Pagination, returnLimit int, returnOffset int, err error) {
	ctx := req.Context()

	// Count total items
	var totalItems int
	if err := app.DB.QueryRow(ctx, "SELECT COUNT(*) FROM "+table).Scan(&totalItems); err != nil {
		return coreEntity.Pagination{}, 0, 0, err
	}

	// Extract limit and offset from query parameters
	queryValues := req.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit // Default offset (starting from 0 for database query)

	var nextPage, previousPage *int
	if page > 1 {
		prevPage := page - 1
		previousPage = &prevPage
	}
	// You may need to adjust the condition based on your total items count
	if offset+limit < totalItems {
		nextPageValue := page + 1
		nextPage = &nextPageValue
	}

	// Calculate total pages
	totalPages := totalItems / limit
	if totalItems%limit != 0 {
		totalPages++
	}

	// Prepare pagination struct
	return coreEntity.Pagination{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		FirstPage:    1,
		LastPage:     totalPages,
	}, limit, offset, nil
}

func Paginate(req *http.Request, app *app.App, baseQuery, filterQuery string) (coreEntity.Pagination, int, int, error) {
	ctx := req.Context()

	// Count total items with filters applied
	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s%s) AS filtered", baseQuery, filterQuery)	
	if app.Config.DBType == "mysql" {
		if err := app.MDB.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
			return coreEntity.Pagination{}, 0, 0, err
		}
	} else {
		if err := app.DB.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
			return coreEntity.Pagination{}, 0, 0, err
		}
	}

	// Extract limit and offset from query parameters
	queryValues := req.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var nextPage, previousPage *int
	if page > 1 {
		prevPage := page - 1
		previousPage = &prevPage
	}
	if offset+limit < totalItems {
		nextPageValue := page + 1
		nextPage = &nextPageValue
	}

	// Calculate total pages
	totalPages := totalItems / limit
	if totalItems%limit != 0 {
		totalPages++
	}

	// Prepare pagination struct
	return coreEntity.Pagination{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		FirstPage:    1,
		LastPage:     totalPages,
	}, limit, offset, nil
}

func RawPagination(sqlQuery string, queryValues map[string][]string) string {
	q := url.Values(queryValues)
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	switch {
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	sqlQuery += " LIMIT " + strconv.Itoa(pageSize) + " OFFSET " + strconv.Itoa(offset)

	return sqlQuery
}

func HashPassword(password string) (string, error) {
	bp := []byte(password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func BodyParse(s interface{}, w http.ResponseWriter, r *http.Request, isValidation bool) error {
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return err
	}

	if isValidation {
		validate := validator.New()
		validateErr := validate.Struct(s)
		if validateErr != nil {
			utils.WriteJSONEValidation(w, http.StatusBadRequest, validateErr.(validator.ValidationErrors))
			return validateErr
		}
	}
	return nil
}
func BodyParseValidation(s interface{}, w http.ResponseWriter, r *http.Request, isValidation bool) error {
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return err
	}

	if isValidation {
		validate := validator.New()
		validateErr := validate.Struct(s)
		if validateErr != nil {
			return validateErr.(validator.ValidationErrors)
		}
	}
	return nil
}

func Round(num float64, places int) float64 {
	if places < 0 {
		panic("places cannot be negative")
	}
	pow := math.Pow(10, float64(places))
	rounded := math.Round(num*pow) / pow
	return rounded
}

func OrderBy(queryValues map[string][]string) string {
	q := url.Values(queryValues)
	orderBy := "created_at"
	if conditions, ok := q["orderBy"]; ok && len(conditions) > 0 {
		orderBy = conditions[0]
	}

	sortOrder := "asc"
	if conditions, ok := q["sortBy"]; ok && len(conditions) > 0 {
		sortOrder = conditions[0]
	}

	orderBy = orderBy + " " + sortOrder

	return orderBy
}

func GenerateUniqueNumber(length int) (string, error) {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		b := rand.Intn(10)                       // Generate a random digit (0-9)
		buffer.WriteString(fmt.Sprintf("%d", b)) // Convert digit to string and append
	}
	return buffer.String(), nil
}
