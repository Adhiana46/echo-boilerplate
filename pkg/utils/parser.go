package utils

import (
	"errors"
	"strings"
)

// parse sortBy=name.asc,updated_at.desc -> map[string]string
func QuerySortToMap(sortBy string) (map[string]string, error) {
	if sortBy == "" {
		return map[string]string{}, nil
	}

	result := map[string]string{}
	raws := strings.Split(sortBy, ",")
	for _, raw := range raws {
		chunks := strings.Split(raw, ".")

		if len(chunks) != 2 {
			return nil, errors.New("malformed sortBy query parameter, should be field.orderdirection")
		}

		field, order := chunks[0], chunks[1]
		order = strings.ToLower(order)

		if order != "asc" && order != "desc" {
			return nil, errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
		}

		result[field] = order
	}

	return result, nil
}
