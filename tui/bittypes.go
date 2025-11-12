package tui

// pixelCoord represents a coordinate on the character grid, with support for half-pixels
type pixelCoord struct {
	x      float64 // Use float64 for half-pixel precision
	y      int
	isHalf bool // Flag to indicate if this is a half-pixel position
}

// DescenderInfo holds information about a character's descender properties
type DescenderInfo struct {
	HasDescender    bool
	BaselineHeight  int // Height of the main character body (excluding descender)
	DescenderHeight int // Height of the descender part
	TotalHeight     int // Total character height
	VerticalOffset  int // How much to offset this character vertically
}

// ShadowStyleOption represents shadow style options
type ShadowStyleOption struct {
	Name string
	Char rune
	Hex  string
}

// Default shadow style options
var shadowStyleOptions = []ShadowStyleOption{
	{"Light Shade", '░', ""},  // U+2591 LIGHT SHADE - Uses main text color
	{"Medium Shade", '▒', ""}, // U+2592 MEDIUM SHADE - Uses main text color
	{"Dark Shade", '▓', ""},   // U+2593 DARK SHADE - Uses main text color
}
