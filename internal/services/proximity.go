package services

import (
	"server/internal/models"
)

const ProximityThreshold = 50.0 

func CheckProximity(p1, p2 models.Position) bool {
    dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= ProximityThreshold*ProximityThreshold
}