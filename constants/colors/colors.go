package constants

// ANSI escape codes for terminal text coloring
// These constants provide consistent color formatting for different message types:
// - Reset returns to default terminal colors
// - Colors are used for success (Green), warnings (Yellow), errors (Red), and information (Blue/Cyan)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)
