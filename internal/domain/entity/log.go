package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog represents a log entry for a user or system action.
type AuditLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp time.Time          `bson:"timestamp"`
	Username  string             `bson:"username"` // Can be "system" for system events
	IPAddress string             `bson:"ipAddress"`
	Action    string             `bson:"action"`   // E.g., "Login", "CreateAsset", "UpdateConfig"
	Endpoint  string             `bson:"endpoint"` // The API endpoint that was hit
	Details   string             `bson:"details"`  // E.g., "User logged in successfully" or "Asset 'Laptop-123' created"
}

// PermissionLog tracks changes to DAC permissions.
type PermissionLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp  time.Time          `bson:"timestamp"`
	ActorID    primitive.ObjectID `bson:"actorId"`    // Who made the change
	TargetID   primitive.ObjectID `bson:"targetId"`   // The user whose permissions were changed
	ResourceID primitive.ObjectID `bson:"resourceId"` // The asset/resource being affected
	Action     string             `bson:"action"`     // "Grant", "Revoke"
	Permission string             `bson:"permission"` // The permission that was changed (e.g., "read", "write")
}
