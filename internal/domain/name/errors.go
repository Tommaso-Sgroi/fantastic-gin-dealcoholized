package name

import "github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"

// Domain-level aliases keep usage ergonomic while reusing shared app errors.
var (
	ErrNotFound     = errs.ErrNotFound
	ErrInvalidInput = errs.ErrInvalidInput
	ErrConflict     = errs.ErrConflict
)
