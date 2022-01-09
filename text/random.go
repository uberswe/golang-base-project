package text

import "math/rand"

// RandomString generates a random string of n length. Based on https://stackoverflow.com/a/22892986/1260548
func RandomString(n int) string {
	// remove vowels to make it less likely to generate something offensive
	var letters = []rune("bcdfghjklmnpqrstvwxzBCDFGHJKLMNPQRSTVWXZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
