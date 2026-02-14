package authservice

import "centralService/internal/service/errors"

var (
	ErrUserInvalidId       = errors.New("INVALID_USER_ID", "O ID de usuário está invalido")
	ErrUserAlreadyExist    = errors.New("USER_ALREADY_EXIST", "O usuário já possui um cadastro")
	ErrUserNotFound        = errors.New("USER_NOT_FOUND", "O usuário não possui perfil no sistema")
	ErrUserOutOfDomain     = errors.New("USER_EMAIL_NOT_INTERNAL", "O email do usuário não pertence ao email comporativo.")
	ErrUserNoHasTargetRole = errors.New("USER_NOT_AUTHORIZED", "O usuário não tem a permissão especifica para realizar esta operação")
)

var (
	ErrTokenMissing = errors.New(
		"TOKEN_MISSING",
		"Token de autenticação não informado.",
	)

	ErrTokenMalformed = errors.New(
		"TOKEN_MALFORMED",
		"O token de autenticação está malformado.",
	)

	ErrTokenInvalidSignature = errors.New(
		"TOKEN_INVALID_SIGNATURE",
		"O token de autenticação é inválido.",
	)

	ErrTokenExpired = errors.New(
		"TOKEN_EXPIRED",
		"O token de autenticação expirou.",
	)

	ErrTokenNotValidYet = errors.New(
		"TOKEN_NOT_VALID_YET",
		"O token de autenticação ainda não é válido.",
	)

	ErrTokenInvalidIssuer = errors.New(
		"TOKEN_INVALID_ISSUER",
		"O emissor do token de autenticação é inválido.",
	)

	ErrTokenInvalidAudience = errors.New(
		"TOKEN_INVALID_AUDIENCE",
		"O token de autenticação não foi emitido para este serviço.",
	)

	ErrTokenInvalidAlgorithm = errors.New(
		"TOKEN_INVALID_ALGORITHM",
		"O algoritmo do token de autenticação não é permitido.",
	)

	ErrTokenInvalidKey = errors.New(
		"TOKEN_INVALID_KEY",
		"Falha interna ao validar o token de autenticação.",
	)

	ErrTokenClaimsInvalid = errors.New(
		"TOKEN_INVALID_CLAIMS",
		"Os dados do token de autenticação são inválidos.",
	)

	ErrTokenUnknown = errors.New(
		"TOKEN_UNKNOWN_ERROR",
		"Não foi possível validar o token de autenticação.",
	)
)
