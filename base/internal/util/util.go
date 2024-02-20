package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type (
	Param struct {
		Search           Search `json:"search"`
		Limit            int    // holds the calculated limit
		Offset           int    // holds the offset number
		PaginationString string // holds the limit/offset tring
		Sort             string // holds the calculated sort string
		ColumnMapping    map[string]string
	}

	Search struct {
		Filters    []Filter   `json:"filters"`
		Pagination Pagination `json:"pagination"`
		Sort       string     `json:"sort"` // comma separated string, use a '-' before column name to sort DESC i.e.: id,-name => "SORT BY id ASC, name DESC"
	}

	Filter struct {
		Column  string      `json:"column"`
		Compare string      `json:"compare"`
		Value   interface{} `json:"value"`
	}

	Pagination struct {
		PageLimit  int `json:"page_limit"`
		PageNumber int `json:"page_number"`
	}
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomString(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // get desired length
}

func GetTypeCount(i interface{}) int {
	switch reflect.ValueOf(i).Kind() {
	case reflect.Map:
		return reflect.ValueOf(i).Len()
	case reflect.Array:
		return reflect.ValueOf(i).Len()
	case reflect.Slice:
		return reflect.ValueOf(i).Len()
	default:
		return 1
	}
}

func ValidJson(jsonValue json.RawMessage) bool {
	bValue, err := jsonValue.MarshalJSON()
	if err != nil {
		return false
	}
	check := make(map[string]interface{}, 0)
	if errCheck := json.Unmarshal(bValue, &check); errCheck != nil {
		return false
	}
	return true
}

func (p *Param) CalculateParam(primarySort string, availableSort map[string]string) (err error) {
	p.ColumnMapping = availableSort
	// calculate the limit
	if p.Search.Pagination.PageLimit > 0 {
		if p.Search.Pagination.PageNumber == 0 {
			// should not be empty, default to first page
			p.Search.Pagination.PageNumber = 1
		}
		p.Limit = p.Search.Pagination.PageLimit
		p.Offset = p.Search.Pagination.PageNumber - 1
		p.Offset *= p.Search.Pagination.PageLimit
	}
	// calculate the sort
	if primarySort == "" {
		return
	}
	if p.Search.Sort == "" {
		p.Search.Sort = primarySort
	}
	sorted := []string{}
	sortParts := strings.Split(p.Search.Sort, ",")
	for _, s := range sortParts {
		direction := "ASC"
		name := s
		if string(name[0]) == "-" {
			direction = "DESC"
			name = string(name[1:])
		}
		if _, ok := availableSort[name]; !ok {
			// if the name is not in the available sort list, you could return and error here
			continue
		}
		sorted = append(sorted, fmt.Sprintf("%s %s", availableSort[name], direction))
	}
	p.Sort = strings.Join(sorted, ", ")
	return
}

// for each element in 'compare', if NOT in 'src', it will be added to the resulting 'diff'
// if you want to add from an existing, compared to another list, src => existing; compare => new list
// if you wnat to delete from to an existing, compared to another list, src => new list; compare => existing
func ArrayDiff(src, compare []string) (diff []string) {
	m := make(map[string]struct{})
	for _, i := range src {
		m[i] = struct{}{}
	}
	for _, i := range compare {
		if _, ok := m[i]; !ok {
			diff = append(diff, i)
		}
	}
	return
}
