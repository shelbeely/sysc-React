package tui

// ANSIColorMap provides mappings from ANSI color codes to hex values
// This is the canonical color mapping used throughout the library
var ANSIColorMap = map[string]string{
	"30": "#000000", // Black
	"31": "#CD3131", // Red
	"32": "#0DBC79", // Green
	"33": "#E5E510", // Yellow
	"34": "#2472C8", // Blue
	"35": "#BC3FBC", // Magenta
	"36": "#11A8CD", // Cyan
	"37": "#E5E5E5", // White
	"90": "#808080", // Gray
	"91": "#FF9999", // Bright Red
	"92": "#99FF99", // Bright Green
	"93": "#FFFF99", // Bright Yellow
	"94": "#66BBFF", // Bright Blue
	"95": "#FF99FF", // Bright Magenta
	"96": "#99FFFF", // Bright Cyan
	"97": "#FFFFFF", // Bright White
}
