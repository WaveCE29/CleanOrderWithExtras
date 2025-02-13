package services

import (
	"testing"

	"github.com/WaveCE29/custom-film-cleaner/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCleanOrderWithComplementaryItems(t *testing.T) {
	input := models.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         50,
		TotalPrice:        100,
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Qty:        2,
			UnitPrice:  50,
			TotalPrice: 100,
		},
		{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         3,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	result, err := CleanOrder(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCleanOrderWithSpecialChars(t *testing.T) {

	input := models.InputOrder{
		No:                1,
		PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         50,
		TotalPrice:        100,
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Qty:        2,
			UnitPrice:  50,
			TotalPrice: 100,
		},
		{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         3,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	cleanedOrders, err := CleanOrder(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, cleanedOrders)
}

func TestCleanOrderWithMultiplier(t *testing.T) {
	input := models.InputOrder{
		No:                1,
		PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
		Qty:               1,
		UnitPrice:         90,
		TotalPrice:        90,
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
			MaterialId: "FG0A-MATTE",
			ModelId:    "IPHONE16PROMAX",
			Qty:        3,
			UnitPrice:  30,
			TotalPrice: 90,
		},
		{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         3,
			ProductId:  "MATTE-CLEANNER",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	result, err := CleanOrder(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCleanOrderWithSpecialCharsAndSplitProducts(t *testing.T) {
	input := models.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
		Qty:               1,
		UnitPrice:         80,
		TotalPrice:        80,
	}
	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         2,
			ProductId:  "FG0A-CLEAR-OPPOA3-B",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3-B",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         3,
			ProductId:  "WIPING-CLOTH",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         4,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	cleanedOrders, err := CleanOrder(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, cleanedOrders)
}

func TestCleanOrderWithMultipleProductsAndSplitFormats(t *testing.T) {
	input := models.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
		Qty:               1,
		UnitPrice:         120,
		TotalPrice:        120,
	}

	expected := []models.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         2,
			ProductId:  "FG0A-CLEAR-OPPOA3-B",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3-B",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         3,
			ProductId:  "FG0A-MATTE-OPPOA3",
			MaterialId: "FG0A-MATTE",
			ModelId:    "OPPOA3",
			Qty:        1,
			UnitPrice:  40,
			TotalPrice: 40,
		},
		{
			No:         4,
			ProductId:  "WIPING-CLOTH",
			Qty:        3,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         5,
			ProductId:  "CLEAR-CLEANNER",
			Qty:        2,
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         6,
			ProductId:  "MATTE-CLEANNER",
			Qty:        1,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	result, err := CleanOrder(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
