package service

import (
	"fmt"
	"math"
	"sort"

	"pack-calculator/internal/domain/model"
)

type PackCalculator struct{}

func NewPackCalculator() *PackCalculator {
	return &PackCalculator{}
}

func (pc *PackCalculator) Calculate(
	packSizes []int,
	orderQuantity int,
) (model.PackDistribution, error) {
	if len(packSizes) == 0 {
		return nil, model.ErrEmptyPackSizes
	}
	if orderQuantity <= 0 {
		return nil, model.ErrInvalidOrderQuantity
	}
	for _, size := range packSizes {
		if size <= 0 {
			return nil, model.ErrInvalidPackSize
		}
	}

	// Sort descending for consistency
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	type state struct {
		overage      int
		packCount    int
		distribution model.PackDistribution
	}

	memo := make(map[int]state)

	var dfs func(int) state
	dfs = func(remaining int) state {
		if val, ok := memo[remaining]; ok {
			return val
		}

		if remaining <= 0 {
			return state{
				overage:      -remaining,
				packCount:    0,
				distribution: make(model.PackDistribution),
			}
		}

		best := state{
			overage:   math.MaxInt32,
			packCount: math.MaxInt32,
		}

		for _, size := range packSizes {
			sub := dfs(remaining - size)

			newDist := make(model.PackDistribution)
			for k, v := range sub.distribution {
				newDist[k] = v
			}
			newDist[size]++

			newPackCount := sub.packCount + 1
			newOverage := sub.overage

			if newOverage < best.overage ||
				(newOverage == best.overage && newPackCount < best.packCount) {
				best = state{
					overage:      newOverage,
					packCount:    newPackCount,
					distribution: newDist,
				}
			}
		}

		memo[remaining] = best
		return best
	}

	result := dfs(orderQuantity)

	if !result.distribution.CanFulfill(orderQuantity) {
		return nil, fmt.Errorf("unable to fulfill order")
	}

	return result.distribution, nil
}
