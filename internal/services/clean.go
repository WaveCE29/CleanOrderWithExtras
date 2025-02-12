package services

import (
	"regexp"
	"strings"

	"strconv"

	"github.com/WaveCE29/custom-film-cleaner/internal/models"
	"github.com/pkg/errors"
)

// CleanOrder processes the input order and returns the cleaned orders with the correct quantities.
func CleanOrder(input models.InputOrder) ([]models.CleanedOrder, error) {
	// Split the PlatformProductId by the separator '/%20x'
	products := strings.Split(input.PlatformProductId, "/%20x")
	var cleanedOrders []models.CleanedOrder
	var totalQty int
	var textureId string

	// Loop through the split products to handle each one
	for i, product := range products {
		// Extract the product details using the regex
		productRegex := regexp.MustCompile(`([A-Z0-9\-]+-[A-Z]+)-([A-Z0-9]+)(\*([0-9]+))?`)
		matches := productRegex.FindStringSubmatch(product)

		if len(matches) == 0 {
			return nil, errors.New("invalid PlatformProductId format")
		}

		// Extract the base product ID and quantity multiplier
		productId := matches[1] + "-" + matches[2] // Correctly build the product ID (do not append modelId again)
		multiplier := 1
		if len(matches) > 4 && matches[4] != "" {
			// If multiplier is present (e.g., *3), use it
			multiplier, _ = strconv.Atoi(matches[4])
		}

		// Split the texture part (e.g., 'FG0A-MATTE' becomes 'MATTE')
		textureParts := strings.Split(matches[1], "-")
		textureId = textureParts[len(textureParts)-1] // Last part is the texture (e.g., 'MATTE')

		// Adjust quantities and prices based on the multiplier
		adjustedQty := input.Qty * multiplier             // The base quantity is multiplied by the multiplier
		adjustedUnitPrice := input.UnitPrice / multiplier // The price should be halved for each product
		adjustedTotalPrice := input.TotalPrice            // The total price should be halved for each product

		// Create cleaned order for the current product
		cleanedOrder := models.CleanedOrder{
			No:         input.No + i,
			ProductId:  productId,
			MaterialId: matches[1],         // e.g., 'FG0A-MATTE'
			ModelId:    matches[2],         // e.g., 'OPPOA3'
			Qty:        adjustedQty,        // Adjusted quantity
			UnitPrice:  adjustedUnitPrice,  // Adjusted price (halved for second product)
			TotalPrice: adjustedTotalPrice, // Adjusted total price (halved for second product)
		}

		// Add the main product to cleanedOrders
		cleanedOrders = append(cleanedOrders, cleanedOrder)

		// Update the total quantity for complementary products
		totalQty += adjustedQty
	}

	// Now, add the complementary products with the adjusted quantity
	complementaryOrders := []models.CleanedOrder{
		{
			No:         input.No + len(products),
			ProductId:  "WIPING-CLOTH",
			Qty:        totalQty, // Adjusted quantity based on total qty of main products
			UnitPrice:  0,
			TotalPrice: 0,
		},
		{
			No:         input.No + len(products) + 1,
			ProductId:  textureId + "-CLEANNER",
			Qty:        totalQty, // Adjusted quantity based on total qty of main products
			UnitPrice:  0,
			TotalPrice: 0,
		},
	}

	// Add complementary orders
	cleanedOrders = append(cleanedOrders, complementaryOrders...)

	// Return all cleaned orders
	return cleanedOrders, nil
}
