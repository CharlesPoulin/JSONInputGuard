package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validateOnce sync.Once
	v            *validator.Validate
)

// V returns a fast, singleton validator instance.
func V() *validator.Validate {
	validateOnce.Do(func() {
		v = validator.New(validator.WithRequiredStructEnabled())
	})
	return v
}
