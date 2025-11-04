package utils

import "fmt"

func GenerateProductPK(category string) string {
	return fmt.Sprintf("CATEGORY#%s", category)
}

func GenerateProductSK(productID string) string {
	return fmt.Sprintf("PRODUCT#%s", productID)
}
