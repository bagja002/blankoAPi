package tools

import "strings"

func RemoverSpaci(input string) string {

	remover := strings.ReplaceAll(input, " ", "")

	return remover
}
