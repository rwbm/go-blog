package exception

import "errors"

// Internal error defitions
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Standar error codes
const (
	CodeInternalServerError = "internal_server_error"
	CodeBadRequest          = "bad_request"
	CodeInvalidPage         = "invalid_page"
	CodeInvalidPageSize     = "invalid_page_size"
)

var (
	// map with text messages associated
	messages = map[string]string{
		CodeInternalServerError: "internal server error ocurred",
		CodeBadRequest:          "one or more parameters are missing or wrong",
		CodeInvalidPage:         "invalid page value",
		CodeInvalidPageSize:     "invalid page size value",
	}
)

// GetErrorMap returns a map with the provided error code and associated message;
// it's useful for building HTTP error responses
func GetErrorMap(code, msg string) (m map[string]interface{}) {
	if code != "" || msg != "" {
		m = make(map[string]interface{})

		if code != "" {
			m["code"] = code
		} else {
			m["code"] = CodeInternalServerError
		}

		if msg != "" {
			m["message"] = msg
		} else if code != "" {
			if message, ok := messages[code]; ok {
				m["message"] = message
			}
		}
	}
	return
}

// GetErrorMapWithFields returns a map with the provided error code and associated message;
// it also has the chance to set `fields`
func GetErrorMapWithFields(code, msg, fields string) (m map[string]interface{}) {
	if code != "" || msg != "" {
		m = make(map[string]interface{})

		if code != "" {
			m["code"] = code
		} else {
			m["code"] = CodeInternalServerError
		}

		if msg != "" {
			m["message"] = msg
		} else if code != "" {
			if message, ok := messages[code]; ok {
				m["message"] = message
			}
		}

		if fields != "" {
			m["fields"] = fields
		}
	}
	return
}
