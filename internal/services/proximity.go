package services

import (
	"math"
	"server/internal/models"
)

func CheckProximity(p1, p2 models.Position, rangeLimit float64) bool {
    squaredDistance := math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
    return squaredDistance <= math.Pow(rangeLimit, 2)
}