package tools

import "strconv"



func IntToString(x int)string{
	Number := strconv.Itoa(x)
	return Number
}

func StringToInt(x string)int {
	Number,_ := strconv.Atoi(x)
	return Number
} 