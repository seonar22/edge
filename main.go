package edge

import "go.mongodb.org/mongo-driver/bson/primitive"

// Container stores a single container
type Container struct {
	ID    primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Users []primitive.ObjectID `json:"users"`
	Image string               `json:"image" validate:"required"`
	Ports []string             `json:"ports" validate:"required"`
	Env   map[string]string    `json:"env"`
}
