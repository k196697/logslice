package parser

// SupportedFormats lists the log formats that logslice can parse.
var SupportedFormats = []string{
	"json",
	"csv",
	"logfmt",
}

// FormatDescription returns a human-readable description for a given format
// name, or an empty string if the format is unknown.
func FormatDescription(format string) string {
	switch format {
	case "json":
		return "Newline-delimited JSON (one object per line)"
	case "csv":
		return "Comma-separated values with a header row"
	case "logfmt":
		return "Heroku-style logfmt (key=value pairs)"
	default:
		return ""
	}
}

// IsSupported reports whether format is in SupportedFormats.
func IsSupported(format string) bool {
	for _, f := range SupportedFormats {
		if f == format {
			return true
		}
	}
	return false
}
