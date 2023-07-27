package nostr

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func Colorize(color, s string) string {
	return color + s + Reset
}

func inRed(s string) string {
	return Colorize(Red, s)
}

func inGreen(s string) string {
	return Colorize(Green, s)
}

func inYellow(s string) string {
	return Colorize(Yellow, s)
}

func inBlue(s string) string {
	return Colorize(Blue, s)
}

func inPurple(s string) string {
	return Colorize(Purple, s)
}

func inCyan(s string) string {
	return Colorize(Cyan, s)
}

func inGray(s string) string {
	return Colorize(Gray, s)
}

func inWhite(s string) string {
	return Colorize(White, s)
}
