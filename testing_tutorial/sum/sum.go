package sum

func Ints(i ...int) int {
	return ints(i)
}

func ints(i []int) int {
	if len(i) == 0 {
		return 0
	}
	return ints(i[1:]) + i[0]
}

func Double(i int) int{
return i*2
}
