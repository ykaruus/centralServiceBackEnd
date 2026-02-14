package response

var (
	ErrBadRequest = ErrorBody{
		Code:    "BAD_REQUEST",
		Message: "Requisição inválida",
	}

	ErrUnauthorized = ErrorBody{
		Code:    "UNAUTHORIZED",
		Message: "Não autorizado",
	}

	ErrForbidden = ErrorBody{
		Code:    "FORBIDDEN",
		Message: "Acesso negado",
	}

	ErrNotFound = ErrorBody{
		Code:    "NOT_FOUND",
		Message: "Recurso não encontrado",
	}

	ErrInternal = ErrorBody{
		Code:    "INTERNAL_ERROR",
		Message: "Erro interno do servidor",
	}
)
