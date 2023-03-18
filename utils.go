package notifier

func ArrayContains(arr []string, s string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}

	return false
}
