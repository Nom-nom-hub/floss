// Package anim provides enhanced animated spinner functionality.
package anim

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/zeebo/xxh3"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/nom-nom-hub/floss/internal/csync"
)

const (
	enhancedFps           = 30  // Increased FPS for smoother animations according to UI/UX specification
	enhancedInitialChar   = '.'
	enhancedLabelGap      = " "
	enhancedLabelGapWidth = 1

	// Periods of ellipsis animation speed in steps.
	//
	// If the FPS is 30 (33 milliseconds) this means that the ellipsis will
	// change every 5 frames (165 milliseconds) for smoother transitions.
	enhancedEllipsisAnimSpeed = 5

	// The maximum amount of time that can pass before a character appears.
	// This is used to create a staggered entrance effect.
	// Reduced for more responsive animations according to UI/UX specification.
	enhancedMaxBirthOffset = 500 * time.Millisecond

	// Number of frames to prerender for the animation. After this number
	// of frames, the animation will loop. This only applies when color
	// cycling is disabled.
	// Reduced for better performance according to UI/UX specification.
	enhancedPrerenderedFrames = 5

	// Default number of cycling chars.
	enhancedDefaultNumCyclingChars = 10
)

// Default colors for gradient.
var (
	enhancedDefaultGradColorA = color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}
	enhancedDefaultGradColorB = color.RGBA{R: 0, G: 0, B: 0xff, A: 0xff}
	enhancedDefaultLabelColor = color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xff}
	
	enhancedAvailableRunes = []rune("0123456789abcdefABCDEF~!@#$£€%^&*()+=_")
	
	// Enhanced ellipsis frames for more sophisticated animations
	enhancedEllipsisFrames = []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
	}
)

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var enhancedLastID int64

func enhancedNextID() int {
	return int(atomic.AddInt64(&enhancedLastID, 1))
}

// Cache for expensive animation calculations
type enhancedAnimCache struct {
	initialFrames  [][]string
	cyclingFrames  [][]string
	width          int
	labelWidth     int
	label          []string
	ellipsisFrames []string
}

var enhancedAnimCacheMap = csync.NewMap[string, *enhancedAnimCache]()

// enhancedSettingsHash creates a hash key for the settings to use for caching
func enhancedSettingsHash(opts Settings) string {
	h := xxh3.New()
	fmt.Fprintf(h, "%d-%s-%v-%v-%v-%t",
		opts.Size, opts.Label, opts.LabelColor, opts.GradColorA, opts.GradColorB, opts.CycleColors)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// EnhancedStepMsg is a message type used to trigger the next step in the enhanced animation.
type EnhancedStepMsg struct{ id int }

// FlossAnimationType represents different types of enhanced animations
type FlossAnimationType int

const (
	// Standard cycling animation
	FlossStandard FlossAnimationType = iota
	// Pulse animation with intensity variations
	FlossPulse
	// Wave animation with sequential character lighting
	FlossWave
	// Bounce animation with characters moving up and down
	FlossBounce
)

// FlossAnimationSettings defines enhanced settings for the animation with more options.
type FlossAnimationSettings struct {
	Size        int
	Label       string
	LabelColor  color.Color
	GradColorA  color.Color
	GradColorB  color.Color
	CycleColors bool
	// New enhanced options
	AnimType  FlossAnimationType
	Speed          int // 1-10, with 10 being fastest
	PulseIntensity float64 // 0.0-1.0 for pulse effect intensity
}

// FlossEnhancedAnim extends the basic Anim with more sophisticated animations.
type FlossEnhancedAnim struct {
	*Anim // Embed the base Anim
	
	// Enhanced properties
	animationType  FlossAnimationType
	speed          int
	pulseIntensity float64
	pulseStep      atomic.Int64
	waveOffset     atomic.Int64
	bounceStep     atomic.Int64
}

// NewFlossEnhanced creates a new FlossEnhancedAnim instance with enhanced animation options.
func NewFlossEnhanced(opts FlossAnimationSettings) *FlossEnhancedAnim {
	// Convert FlossAnimationSettings to Settings for compatibility
	baseOpts := Settings{
		Size:        opts.Size,
		Label:       opts.Label,
		LabelColor:  opts.LabelColor,
		GradColorA:  opts.GradColorA,
		GradColorB:  opts.GradColorB,
		CycleColors: opts.CycleColors,
	}
	
	// Create the base animation
	baseAnim := New(baseOpts)
	
	// Create enhanced animation
	a := &FlossEnhancedAnim{
		Anim:           baseAnim,
		animationType:  opts.AnimType,
		speed:          opts.Speed,
		pulseIntensity: opts.PulseIntensity,
	}
	
	// Validate settings
	if a.speed < 1 {
		a.speed = 1
	}
	if a.speed > 10 {
		a.speed = 10
	}
	if a.pulseIntensity < 0.0 {
		a.pulseIntensity = 0.0
	}
	if a.pulseIntensity > 1.0 {
		a.pulseIntensity = 1.0
	}
	
	// Override some properties for enhanced functionality
	a.id = enhancedNextID()
	
	return a
}

// Update processes animation steps (or not).
func (a *FlossEnhancedAnim) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EnhancedStepMsg:
		if msg.id != a.id {
			// Reject messages that are not for this instance.
			return a, nil
		}

		step := a.step.Add(1)
		if int(step) >= len(a.cyclingFrames) {
			a.step.Store(0)
		}

		// Update animation-specific steps based on type
		switch a.animationType {
		case FlossPulse:
			pulseStep := a.pulseStep.Add(1)
			if int(pulseStep) >= 20 { // 20 steps for a full pulse cycle
				a.pulseStep.Store(0)
			}
		case FlossWave:
			waveOffset := a.waveOffset.Add(1)
			if int(waveOffset) >= a.width {
				a.waveOffset.Store(0)
			}
		case FlossBounce:
			bounceStep := a.bounceStep.Add(1)
			if int(bounceStep) >= 20 { // 20 steps for a full bounce cycle
				a.bounceStep.Store(0)
			}
		}

		if a.initialized.Load() && a.labelWidth > 0 {
			// Manage the ellipsis animation.
			ellipsisStep := a.ellipsisStep.Add(1)
			if int(ellipsisStep) >= enhancedEllipsisAnimSpeed*len(enhancedEllipsisFrames) {
				a.ellipsisStep.Store(0)
			}
		} else if !a.initialized.Load() && time.Since(a.startTime) >= enhancedMaxBirthOffset {
			a.initialized.Store(true)
		}
		return a, a.Step()
	default:
		return a.Anim.Update(msg) // Delegate to base animation for other messages
	}
}

// View renders the current state of the enhanced animation.
func (a *FlossEnhancedAnim) View() string {
	var b strings.Builder
	step := int(a.step.Load())
	
	// Adjust speed - higher speed means faster animation
	speedAdjustedStep := step * (11 - a.speed) / 10
	
	for i := range a.width {
		switch {
		case !a.initialized.Load() && i < len(a.birthOffsets) && time.Since(a.startTime) < a.birthOffsets[i]:
			// Birth offset not reached: render initial character.
			b.WriteString(a.initialFrames[speedAdjustedStep][i])
		case i < a.cyclingCharWidth:
			// Render a cycling character with enhanced animation effects.
			b.WriteString(a.enhancedCyclingChar(speedAdjustedStep, i))
		case i == a.cyclingCharWidth:
			// Render label gap.
			b.WriteString(enhancedLabelGap)
		case i > a.cyclingCharWidth:
			// Label.
			if labelChar, ok := a.label.Get(i - a.cyclingCharWidth - enhancedLabelGapWidth); ok {
				b.WriteString(labelChar)
			}
		}
	}
	// Render animated ellipsis at the end of the label if all characters
	// have been initialized.
	if a.initialized.Load() && a.labelWidth > 0 {
		ellipsisStep := int(a.ellipsisStep.Load())
		if ellipsisFrame, ok := a.ellipsisFrames.Get(ellipsisStep / enhancedEllipsisAnimSpeed); ok {
			b.WriteString(ellipsisFrame)
		}
	}

	return b.String()
}

// enhancedCyclingChar renders a character with enhanced animation effects based on type.
func (a *FlossEnhancedAnim) enhancedCyclingChar(step, index int) string {
	if step >= len(a.cyclingFrames) {
		step = 0
	}
	if index >= len(a.cyclingFrames[step]) {
		return ""
	}
	
	char := a.cyclingFrames[step][index]
	
	// Apply animation effects based on type
	switch a.animationType {
	case FlossPulse:
		// Apply pulsing effect based on pulse intensity
		pulseStep := int(a.pulseStep.Load())
		pulseIntensity := (1.0 + a.pulseIntensity*math.Sin(2*math.Pi*float64(pulseStep)/20.0)) / 2.0
		if pulseIntensity > 0.7 {
			// Make character bold during peak pulse
			return lipgloss.NewStyle().Bold(true).Render(char)
		}
		return char
	case FlossWave:
		// Apply wave effect where characters light up sequentially
		waveOffset := int(a.waveOffset.Load())
		wavePosition := (index - waveOffset + a.width) % a.width
		if wavePosition < 3 { // Highlight 3 characters in sequence
			return lipgloss.NewStyle().Bold(true).Render(char)
		}
		return char
	case FlossBounce:
		// Apply bouncing effect
		bounceStep := int(a.bounceStep.Load())
		bouncePosition := int(float64(bounceStep) / 20.0 * float64(a.width))
		if index == bouncePosition || index == (bouncePosition+1)%a.width {
			// Make bouncing characters bold
			return lipgloss.NewStyle().Bold(true).Render(char)
		}
		return char
	default:
		// Standard cycling animation
		return char
	}
}

// Step is a command that triggers the next step in the enhanced animation.
func (a *FlossEnhancedAnim) Step() tea.Cmd {
	// Adjust timing based on speed setting
	interval := time.Second / time.Duration(enhancedFps*(11-a.speed)/10)
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return EnhancedStepMsg{id: a.id}
	})
}