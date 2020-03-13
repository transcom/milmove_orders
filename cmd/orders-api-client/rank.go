package main

import (
	"github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
)

// CategorizedRank combines officer status and paygrade
type CategorizedRank struct {
	officer  bool
	paygrade ordersmessages.Rank
}

var paygradeToRank = map[string]CategorizedRank{
	"E01": {officer: false, paygrade: ordersmessages.RankE1},
	"E02": {officer: false, paygrade: ordersmessages.RankE2},
	"E03": {officer: false, paygrade: ordersmessages.RankE3},
	"E04": {officer: false, paygrade: ordersmessages.RankE4},
	"E05": {officer: false, paygrade: ordersmessages.RankE5},
	"E06": {officer: false, paygrade: ordersmessages.RankE6},
	"E07": {officer: false, paygrade: ordersmessages.RankE7},
	"E08": {officer: false, paygrade: ordersmessages.RankE8},
	"E09": {officer: false, paygrade: ordersmessages.RankE9},
	"O01": {officer: true, paygrade: ordersmessages.RankO1},
	"O02": {officer: true, paygrade: ordersmessages.RankO2},
	"O03": {officer: true, paygrade: ordersmessages.RankO3},
	"O04": {officer: true, paygrade: ordersmessages.RankO4},
	"O05": {officer: true, paygrade: ordersmessages.RankO5},
	"O06": {officer: true, paygrade: ordersmessages.RankO6},
	"O07": {officer: true, paygrade: ordersmessages.RankO7},
	"O08": {officer: true, paygrade: ordersmessages.RankO8},
	"O09": {officer: true, paygrade: ordersmessages.RankO9},
	"O10": {officer: true, paygrade: ordersmessages.RankO10},
	"W01": {officer: true, paygrade: ordersmessages.RankW1},
	"W02": {officer: true, paygrade: ordersmessages.RankW2},
	"W03": {officer: true, paygrade: ordersmessages.RankW3},
	"W04": {officer: true, paygrade: ordersmessages.RankW4},
	"W05": {officer: true, paygrade: ordersmessages.RankW5},
}
