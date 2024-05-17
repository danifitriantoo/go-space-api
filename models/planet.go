package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Planet struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Diameter int                `json:"diameter"` // in scale km
	Moons    int                `json:"moons"`
	Distance int                `json:"distance"` // distance from the sun
}
