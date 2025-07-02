package services

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/WaveCE29/custom-film-cleaner/internal/models"
)

func CleanOrder(inputOrders []models.InputOrder) ([]models.CleanedOrder, error) {
	var cleanedOrders []models.CleanedOrder
	orderNo := 1

	// Track materials for complementary items
	materialCounts := make(map[string]int)

	// Process each input order
	for _, order := range inputOrders {
		// Parse the platform product ID
		products := parseProductId(order.PlatformProductId)

		if len(products) == 0 {
			continue
		}

		// Calculate total quantity of all products in this order
		totalProductQty := 0
		for _, product := range products {
			totalProductQty += product.Qty
		}

		// Calculate unit price per product
		var unitPrice float64
		if order.UnitPrice > 0 {
			// Use the provided unit price divided by total quantity
			unitPrice = order.UnitPrice / float64(totalProductQty)
		} else if order.TotalPrice > 0 {
			// Calculate from total price divided by total quantity
			unitPrice = order.TotalPrice / float64(totalProductQty)
		}
		// Add main products
		for _, product := range products {
			materialId, modelId := parseProductComponents(product.Id)

			productQty := product.Qty * order.Qty
			productTotal := unitPrice * float64(productQty)

			cleanedOrder := models.CleanedOrder{
				No:         orderNo,
				ProductId:  product.Id,
				MaterialId: materialId,
				ModelId:    modelId,
				Qty:        productQty,
				UnitPrice:  unitPrice,
				TotalPrice: productTotal,
			}

			cleanedOrders = append(cleanedOrders, cleanedOrder)
			orderNo++

			// Track material counts for complementary items
			if materialId != "" {
				materialCounts[materialId] += productQty
			}
		}
	}

	// Add complementary items
	totalWipingCloth := 0

	// Calculate wiping cloth quantity (sum of all main products)
	for _, order := range cleanedOrders {
		if order.MaterialId != "" {
			totalWipingCloth += order.Qty
		}
	}

	// Add wiping cloth
	if totalWipingCloth > 0 {
		cleanedOrders = append(cleanedOrders, models.CleanedOrder{
			No:         orderNo,
			ProductId:  "WIPING-CLOTH",
			Qty:        totalWipingCloth,
			UnitPrice:  0,
			TotalPrice: 0,
		})
		orderNo++
	}

	// Add cleaners for each material type
	for material, qty := range materialCounts {
		// Extract just the material name (remove FG0A- prefix)
		materialName := strings.Replace(material, "FG0A-", "", 1)
		cleanerName := materialName + "-CLEANNER"

		cleanedOrders = append(cleanedOrders, models.CleanedOrder{
			No:         orderNo,
			ProductId:  cleanerName,
			Qty:        qty,
			UnitPrice:  0,
			TotalPrice: 0,
		})
		orderNo++
	}

	return cleanedOrders, nil
}

// Product represents a parsed product with quantity
type Product struct {
	Id  string
	Qty int
}

// parseProductId parses the platform product ID and returns individual products
func parseProductId(productId string) []Product {
	var products []Product

	// Remove leading special characters and clean the string
	cleanId := regexp.MustCompile(`^[^A-Z]*`).ReplaceAllString(productId, "")

	// Handle different splitting patterns
	// Split by: /, /%20x, /...x patterns
	parts := regexp.MustCompile(`/\s*%20\s*x|/%20x|/`).Split(cleanId, -1)

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Clean each part
		part = strings.TrimSpace(part)
		part = regexp.MustCompile(`^[^A-Z]*`).ReplaceAllString(part, "")

		if part == "" {
			continue
		}

		// Check for multiplier pattern (*number)
		qty := 1
		multiplierRegex := regexp.MustCompile(`\*(\d+)$`)
		if matches := multiplierRegex.FindStringSubmatch(part); len(matches) > 1 {
			if parsedQty, err := strconv.Atoi(matches[1]); err == nil {
				qty = parsedQty
			}
			// Remove the multiplier from the product ID
			part = multiplierRegex.ReplaceAllString(part, "")
		}

		// Only add if it starts with a capital letter (valid product format)
		if part != "" && regexp.MustCompile(`^[A-Z]`).MatchString(part) {
			products = append(products, Product{
				Id:  part,
				Qty: qty,
			})
		}
	}

	return products
}

// parseProductComponents extracts material and model from product ID
func parseProductComponents(productId string) (string, string) {
	// Pattern: FG0A-MATERIAL-MODEL
	parts := strings.Split(productId, "-")

	if len(parts) >= 3 && parts[0] == "FG0A" {
		material := parts[0] + "-" + parts[1] // FG0A-CLEAR, FG0A-MATTE, etc.
		model := strings.Join(parts[2:], "-") // IPHONE16PROMAX, OPPOA3, etc.
		return material, model
	}

	return "", ""
}
