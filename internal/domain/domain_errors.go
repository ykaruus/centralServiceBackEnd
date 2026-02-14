package domain

import "errors"

var (
	ErrEquipmentNotAvailable = errors.New("domain: the equipment is not available")
	ErrUserEmailOutDomain    = errors.New("domain: the user email is out of internal domain")
	ErrUserNotHasTargetRole  = errors.New("domain: the user not has the specific role")
)
