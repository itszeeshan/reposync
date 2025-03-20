package helpers

import (
	"log"
	"strconv"

	colors "github.com/itszeeshan/reposync/constants/colors"
)

/*
parseStringToInt safely converts group ID string to integer.
Provides user-friendly error handling for invalid numeric inputs,
ensuring valid API requests with properly formatted group IDs.
*/
func ParseStringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf(colors.Red+"Invalid group ID: %s"+colors.Reset, s)
	}
	return n
}
