package main

func generateValues(begin, end, step int) []int {
	values := make([]int, 0, end-begin+1)
	for val := begin; val <= end; val += step {
		values = append(values, val)
	}

	return values
}
