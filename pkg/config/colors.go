package config

import "github.com/fatih/color"

var (
	// Titles and main elements
	DarkBlue = color.New(color.FgHiBlue).SprintFunc()

	// Success and confirmations
	LightGreen      = color.New(color.FgHiGreen).SprintFunc()
	LightGreenPrint = color.New(color.FgHiGreen).PrintfFunc()

	// Warnings
	Yellow      = color.New(color.FgHiYellow).SprintFunc()
	YellowPrint = color.New(color.FgHiYellow).PrintfFunc()

	// Errors
	Red      = color.New(color.FgHiRed).SprintFunc()
	RedPrint = color.New(color.FgHiRed).PrintfFunc()

	// Information and highlights
	Cyan = color.New(color.FgHiCyan).SprintFunc()

	// Interactive elements
	Purple = color.New(color.FgHiMagenta).SprintFunc()

	// Secondary text and details
	DarkGray = color.New(color.FgHiBlack).Add(color.BgBlack).SprintFunc()

	// Main text on dark backgrounds
	White = color.New(color.FgHiWhite).SprintFunc()

	// Additional background examples
	BgDarkBlue   = color.New(color.BgBlue).SprintFunc()
	BgLightGreen = color.New(color.BgGreen).SprintFunc()
)
