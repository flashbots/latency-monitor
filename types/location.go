package types

type Location [36]byte

func LocationSize() int {
	return 36
}

func (l Location) String() string {
	n := 0
	for n < len(l) {
		if l[n] == 0 {
			break
		}
		n += 1
	}
	return string(l[:n])
}
