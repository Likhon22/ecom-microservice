package utils

import "fmt"

func CreateKey(email string) string {
	key := fmt.Sprintf("cart:%s", email)
	return key

}
