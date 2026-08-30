package mappings

import "net/http"

var (
	RequestBodyParsingError = ErrorDetails{
		InternalCode: "common:request:body-parsing-error",
		StatusCode:   http.StatusBadRequest,
		Message:      "invalid request body",
	}

	InternalServerError = ErrorDetails{
		InternalCode: "common:internal-server-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "internal server error",
	}

	UnauthorizedError = ErrorDetails{
		InternalCode: "common:unauthorized",
		StatusCode:   http.StatusUnauthorized,
		Message:      "unauthorized",
	}

	ForbiddenError = ErrorDetails{
		InternalCode: "common:forbidden",
		StatusCode:   http.StatusForbidden,
		Message:      "forbidden",
	}

	NotFoundError = ErrorDetails{
		InternalCode: "common:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "resource not found",
	}
)
