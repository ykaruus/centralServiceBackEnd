package models

import (
	"centralService/internal/domain/enums"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	Name             string                 `bson:"name"`
	Email            string                 `bson:"email"`
	Role             []enums.UserPermission `bson:"role"`
	Picture          string                 `bson:"picture"`
	ID               bson.ObjectID          `bson:"_id,omitempty"`
	IsFirstLogin     bool                   `bson:"first_login"`
	AssignedToRegion enums.RegionFlag       `bson:"region_flag"`
	CreatedAt        time.Time              `bson:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at"`
	V                int                    `bson:"__v" json:"__v"`
}
