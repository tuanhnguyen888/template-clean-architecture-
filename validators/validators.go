package validators

import (
	"net"

	"github.com/go-playground/validator/v10"
)

func ValidateIpv4(field validator.FieldLevel) bool {
	return net.ParseIP(field.Field().String()) != nil
}
