package parser

// FormatDescription returns a human-readable description of a supported log format.
func FormatDescription(format string) string {
	switch format {
	case "json":
		return "Newline-delimited JSON (one JSON object per line)"
	case "csv":
		return "Comma-separated values with a header row"
	case "logfmt":
		return "Logfmt key=value pairs (e.g. level=info msg=started)"
	default:
		return "Unknown format"
	}
}

// IsSupported reports whether the given format name is supported by logslice.
func IsSupported(format string) bool {
	switch format {
	case "json", "csv", "logfmt":
		return true
	default:
		return false
	}
}

// SupportedFormats returns the list of all supported format names.
func SupportedFormats() []string {
	return []string{"json", "csv", "logfmt"}
}
