package core

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/nom-nom-hub/floss/internal/app"
	"github.com/nom-nom-hub/floss/internal/history"
	"github.com/nom-nom-hub/floss/internal/lsp"
	"github.com/nom-nom-hub/floss/internal/tui/components/chat"
	"github.com/nom-nom-hub/floss/internal/tui/components/chat/header"
	"github.com/nom-nom-hub/floss/internal/tui/components/chat/sidebar"
	"github.com/nom-nom-hub/floss/internal/tui/components/core/status"
	"github.com/nom-nom-hub/floss/internal/tui/components/floss"
	"github.com/nom-nom-hub/floss/internal/csync"
)

// FLOSSFactory creates FLOSS-specific UI components
type FLOSSFactory struct {
	app        *app.App
	history    history.Service
	lspClients *csync.Map[string, *lsp.Client]
}

// NewFLOSSFactory creates a new factory for FLOSS-specific components
func NewFLOSSFactory(app *app.App, history history.Service, lspClients *csync.Map[string, *lsp.Client]) *FLOSSFactory {
	return &FLOSSFactory{
		app:        app,
		history:    history,
		lspClients: lspClients,
	}
}

// CreateChat creates a FLOSS-styled chat component
func (f *FLOSSFactory) CreateChat() chat.MessageListCmp {
	return chat.NewFloss(f.app)
}

// CreateSidebar creates a FLOSS-styled sidebar component
func (f *FLOSSFactory) CreateSidebar(compact bool) sidebar.Sidebar {
	return sidebar.New(f.history, f.lspClients, compact)
}

// CreateHeader creates a FLOSS-styled header component
func (f *FLOSSFactory) CreateHeader() header.Header {
	return header.NewFloss(f.lspClients)
}

// CreateStatus creates a FLOSS-styled status component
func (f *FLOSSFactory) CreateStatus() status.StatusCmp {
	return status.NewFlossStatusCmp()
}

// CreateButton creates a FLOSS-styled button component
func (f *FLOSSFactory) CreateButton(buttonType floss.ButtonType, label string, onPress func() tea.Msg) *floss.Button {
	return floss.NewButton(buttonType, label, onPress)
}

// CreatePrimaryButton creates a FLOSS-styled primary button
func (f *FLOSSFactory) CreatePrimaryButton(label string, onPress func() tea.Msg) *floss.Button {
	return floss.NewPrimaryButton(label, onPress)
}

// CreateSecondaryButton creates a FLOSS-styled secondary button
func (f *FLOSSFactory) CreateSecondaryButton(label string, onPress func() tea.Msg) *floss.Button {
	return floss.NewSecondaryButton(label, onPress)
}

// CreateDangerButton creates a FLOSS-styled danger button
func (f *FLOSSFactory) CreateDangerButton(label string, onPress func() tea.Msg) *floss.Button {
	return floss.NewDangerButton(label, onPress)
}