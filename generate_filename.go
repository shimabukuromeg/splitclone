package split

func genFileName(n int) string {
	var filename string
	// 分割したい数
	val := n
	val2 := val / 26
	if val2 < 1 {
		filename = "xa" + string('a'+rune(val-1))
	} else {
		filename = "x" + string('a'+rune(val2)) + string('a'+rune((val%26)))
	}
	return filename
}
