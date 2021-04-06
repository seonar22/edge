package edge

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

// Container stores a single container
type Container struct {
	ID    primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Users []primitive.ObjectID `json:"users"`
	Image string               `json:"image" validate:"required"`
	Ports []string             `json:"ports" validate:"required,portmap"`
	Env   map[string]string    `json:"env"`
}

// RegisterValidators registers custom container validation handlers with the validator
func RegisterValidators(validate *validator.Validate) error {
	for name, function := range map[string]func(fl validator.FieldLevel) bool{
		"portmap": func(fl validator.FieldLevel) bool {
			for _, port := range fl.Field().Interface().([]string) {
				portPair := strings.Split(port, ":")
				if len(portPair) != 2 {
					return false
				}
				hostPort := portPair[0]
				containerPort := portPair[1]
				if _, err := strconv.Atoi(hostPort); err != nil {
					return false
				}
				if _, err := strconv.Atoi(containerPort); err != nil {
					return false
				}
			}
			return true
		},
	} {
		err := validate.RegisterValidation(name, function)
		if err != nil {
			return err
		}
	}

	return nil // nil error
}
