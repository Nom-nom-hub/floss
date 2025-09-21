package styles

import (
	"fmt"
	"image/color"
	"log"
	"testing"
)

func TestFlossTheme(t *testing.T) {
	// Create and set the FLOSS theme
	themeManager := NewManager()
	flossTheme := NewFlossTheme()
	themeManager.Register(flossTheme)
	themeManager.SetTheme("floss")
	SetDefaultManager(themeManager)

	// Get the current theme
	theme := CurrentTheme()

	// Verify that the FLOSS theme colors are set correctly
	fmt.Println("FLOSS Theme Verification:")
	fmt.Println("========================")

	// Check primary color (should be Citron, not Mauve)
	r1, g1, b1, a1 := theme.Primary.RGBA()
	r2, g2, b2, a2 := ParseHex("#ebbcba").RGBA()
	fmt.Printf("Primary color RGBA: (%d, %d, %d, %d)\n", r1, g1, b1, a1)
	fmt.Printf("Expected Citron RGBA: (%d, %d, %d, %d)\n", r2, g2, b2, a2)
	if isSameColor(theme.Primary, ParseHex("#ebbcba")) {
		fmt.Println("✓ Primary color is Citron (#ebbcba)")
	} else {
		fmt.Println("✗ Primary color is not Citron")
		log.Fatalf("Expected Citron (#ebbcba), got %v", theme.Primary)
	}

	// Check secondary color (should be Guac, not Squid)
	r1, g1, b1, a1 = theme.Secondary.RGBA()
	r2, g2, b2, a2 = ParseHex("#31748f").RGBA()
	fmt.Printf("Secondary color RGBA: (%d, %d, %d, %d)\n", r1, g1, b1, a1)
	fmt.Printf("Expected Guac RGBA: (%d, %d, %d, %d)\n", r2, g2, b2, a2)
	if isSameColor(theme.Secondary, ParseHex("#31748f")) {
		fmt.Println("✓ Secondary color is Guac (#31748f)")
	} else {
		fmt.Println("✗ Secondary color is not Guac")
		log.Fatalf("Expected Guac (#31748f), got %v", theme.Secondary)
	}

	// Check that the theme name is correct
	if theme.Name == "floss" {
		fmt.Println("✓ Theme name is 'floss'")
	} else {
		fmt.Println("✗ Theme name is not 'floss'")
		log.Fatalf("Expected theme name 'floss', got %s", theme.Name)
	}

	// Check background colors
	r1, g1, b1, a1 = theme.BgSubtle.RGBA()
	r2, g2, b2, a2 = ParseHex("#221f30").RGBA()
	fmt.Printf("BgSubtle color RGBA: (%d, %d, %d, %d)\n", r1, g1, b1, a1)
	fmt.Printf("Expected warmed Charcoal RGBA: (%d, %d, %d, %d)\n", r2, g2, b2, a2)
	if isSameColor(theme.BgSubtle, ParseHex("#221f30")) {
		fmt.Println("✓ BgSubtle is warmed Charcoal (#221f30)")
	} else {
		fmt.Println("✗ BgSubtle is not warmed Charcoal")
		log.Fatalf("Expected warmed Charcoal (#221f30), got %v", theme.BgSubtle)
	}

	fmt.Println("\nAll FLOSS theme verifications passed!")
}

// isSameColor compares two colors by their RGBA values
func isSameColor(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}