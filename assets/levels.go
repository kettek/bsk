package assets

func GetLevelSource(s string) ([]byte, error) {
	return fs.ReadFile(s + ".lvl")
}
