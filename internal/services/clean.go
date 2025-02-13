package services

import (
	"regexp"
	"strings"

	"strconv"

	"github.com/WaveCE29/custom-film-cleaner/internal/models"
	"github.com/pkg/errors"
)

func CleanOrder(input models.InputOrder) ([]models.CleanedOrder, error) {
	normalizedProductId := strings.Trim(input.PlatformProductId, "- %")

	productList := regexp.MustCompile(`/[%0-9A-Za-z]*x|/`)
	products := productList.Split(normalizedProductId, -1)

	var cleanedOrders []models.CleanedOrder
	var complementaryOrders []models.CleanedOrder
	var totalQty int

	productRegex := regexp.MustCompile(`([A-Z0-9]+-[A-Z]+)-([A-Z0-9\-]+)(\*([0-9]+))?`)

	for i, product := range products {
		matches := productRegex.FindStringSubmatch(product)
		if len(matches) == 0 {
			return nil, errors.New("invalid PlatformProductId format")
		}

		materialId := matches[1]
		modelId := matches[2]
		multiplier := 1
		if len(matches) > 4 && matches[4] != "" {
			multiplier, _ = strconv.Atoi(matches[4])
		}

		adjustedQty := input.Qty * multiplier
		adjustedUnitPrice := input.UnitPrice / len(products) / multiplier
		adjustedTotalPrice := input.TotalPrice / len(products)

		cleanedOrder := models.CleanedOrder{
			No:         input.No + i,
			ProductId:  materialId + "-" + modelId,
			MaterialId: materialId,
			ModelId:    modelId,
			Qty:        adjustedQty,
			UnitPrice:  adjustedUnitPrice,
			TotalPrice: adjustedTotalPrice,
		}

		cleanedOrders = append(cleanedOrders, cleanedOrder)
		totalQty += adjustedQty
	}

	complementaryOrders = []models.CleanedOrder{
		{
			No:         input.No + len(products),
			ProductId:  "WIPING-CLOTH",
			Qty:        totalQty,
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	textureCount := make(map[string]int)

	for _, order := range cleanedOrders {
		if order.MaterialId != "" {
			textureParts := strings.Split(order.MaterialId, "-")
			textureId := textureParts[len(textureParts)-1]
			textureCount[textureId]++
		}
	}

	if len(textureCount) > 1 {
		for textureId, count := range textureCount {
			complementaryOrders = append(complementaryOrders, models.CleanedOrder{
				No:         input.No + len(products) + len(complementaryOrders),
				ProductId:  textureId + "-CLEANNER",
				Qty:        count,
				UnitPrice:  0,
				TotalPrice: 0,
			})
		}
	} else {
		textureParts := strings.Split(cleanedOrders[0].MaterialId, "-")
		textureId := textureParts[len(textureParts)-1]
		complementaryOrders = append(complementaryOrders, models.CleanedOrder{
			No:         input.No + len(products) + 1,
			ProductId:  textureId + "-CLEANNER",
			Qty:        totalQty,
			UnitPrice:  0,
			TotalPrice: 0,
		})
	}

	cleanedOrders = append(cleanedOrders, complementaryOrders...)

	return cleanedOrders, nil
}
