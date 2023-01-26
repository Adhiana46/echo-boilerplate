package utils

import (
	"strings"

	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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
			return nil, errors.NewBadRequestError("malformed sortBy query parameter, should be field.orderdirection")
		}

		field, order := chunks[0], chunks[1]
		order = strings.ToLower(order)

		if order != "asc" && order != "desc" {
			return nil, errors.NewBadRequestError("malformed orderdirection in sortBy query parameter, should be asc or desc")
		}

		result[field] = order
	}

	return result, nil
}

func ValidationErrors(validationErrs validator.ValidationErrors, trans *ut.Translator) map[string][]string {
	errorFields := map[string][]string{}
	for _, e := range validationErrs {
		if trans != nil {
			errorFields[e.Field()] = append(errorFields[e.Field()], e.Translate(*trans))
		} else {
			errorFields[e.Field()] = append(errorFields[e.Field()], e.Tag())
		}
	}

	return errorFields
}
