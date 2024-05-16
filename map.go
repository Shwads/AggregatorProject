package main

func Map[T any, U any](array []T, function func(T) U) []U {
	output := make([]U, len(array))

	for index, value := range array {
		output[index] = function(value)
	}

	return output
}
