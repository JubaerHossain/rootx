package utilQuery

import (
	"bytes"
	"context"

	// "crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	coreEntity "github.com/JubaerHossain/rootx/pkg/core/entity"
	"github.com/JubaerHossain/rootx/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Pagination(db *pgxpool.Pool, queryValues map[string][]string, query string, args ...interface{}) (string, []interface{}) {
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

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	return query, args
}

func Paginate(ctx context.Context, db *pgxpool.Pool, queryValues map[string][]string, totalItems int, query string, args ...interface{}) ([]pgx.Row, coreEntity.Pagination, error) {
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

    var nextPage, previousPage *int
    if page > 1 {
        prevPage := page - 1
        previousPage = &prevPage
    }
    // You may need to adjust the condition based on your total items count
    if offset+pageSize < totalItems {
        nextPageValue := page + 1
        nextPage = &nextPageValue
    }

    pagination := coreEntity.Pagination{
        TotalItems:   totalItems,
        TotalPages:   int(math.Ceil(float64(totalItems) / float64(pageSize))),
        CurrentPage:  page,
        NextPage:     nextPage,
        PreviousPage: previousPage,
        FirstPage:    1,
        LastPage:     int(math.Ceil(float64(totalItems) / float64(pageSize))),
    }

    query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

    rows, err := db.Query(ctx, query, args...)
    if err != nil {
        return nil, pagination, err
    }

    var result []pgx.Row
    for rows.Next() {
        result = append(result, rows)
    }

    return result, pagination, nil
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
