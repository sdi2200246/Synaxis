package apperr

import "errors"

var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrConflict      = errors.New("conflict")
    ErrBadInput      = errors.New("bad input")
    ErrInternal      = errors.New("internal server error")
)


var (
    ErrUsernameConflict  = errors.New("username already taken")
    ErrEmailConflict     = errors.New("email already taken")
    ErrPendingApproval   = errors.New("registration pending admin approval")
    ErrRejected          = errors.New("registration was rejected by admin")
    ErrInvalidPassword   = errors.New("invalid username or password")
)

var (
    ErrCannotPublishWithoutTickets = errors.New("cannot publish event without tickets")
    ErrCannotPublishPastEvent      = errors.New("cannot publish an event in the past")
    ErrInvalidEventStatus          = errors.New("invalid event status")
)