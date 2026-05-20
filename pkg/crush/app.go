package crush

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/crush/internal/app"
	"github.com/charmbracelet/crush/internal/db"
)

// App is the top-level application orchestrator.
type App = app.App

// NewApp initializes a new application instance from a database connection
// and configuration store.
func NewApp(ctx context.Context, conn *sql.DB, store *ConfigStore, skillsMgr *SkillsManager) (*App, error) {
	return app.New(ctx, conn, store, skillsMgr)
}

// NewAppWithConfig loads the configuration from default paths, opens the
// SQLite database (using the shared connection pool), and creates a new
// App instance. The returned App owns the database connection; call
// Shutdown() when finished.
func NewAppWithConfig(ctx context.Context, workingDir, dataDir string, debug bool, skillsMgr *SkillsManager) (*App, error) {
	store, err := Load(workingDir, dataDir, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Ensure the data directory exists before opening the database.
	if err := os.MkdirAll(store.Config().Options.DataDirectory, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	conn, err := db.Connect(ctx, store.Config().Options.DataDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	a, err := NewApp(ctx, conn, store, skillsMgr)
	if err != nil {
		db.Release(store.Config().Options.DataDirectory)
		return nil, err
	}

	return a, nil
}

// RunPrompt runs a single prompt in non-interactive mode with sensible defaults.
// Output is written to os.Stdout, models come from config, the spinner is shown,
// and a new session is created. For more control (custom output, model overrides,
// session continuation), use [App.RunNonInteractive] directly.
func RunPrompt(app *App, ctx context.Context, prompt string) error {
	return app.RunNonInteractive(ctx, os.Stdout, prompt, "", "", false, "", false)
}

// RunPromptAndCreateSession creates a new named session and runs a prompt in it,
// returning the session ID so callers can retrieve messages later via
// [App.Sessions] and [App.Messages].
func RunPromptAndCreateSession(app *App, ctx context.Context, output io.Writer, title, prompt string) (string, error) {
	sess, err := app.Sessions.Create(ctx, title)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sess.ID, app.RunNonInteractive(ctx, output, prompt, "", "", false, sess.ID, false)
}

// RunPromptInSession runs a prompt in an existing session (or a new one if
// sessionID is empty), streaming output to the given writer.
func RunPromptInSession(app *App, ctx context.Context, output io.Writer, sessionID, prompt string) error {
	return app.RunNonInteractive(ctx, output, prompt, "", "", false, sessionID, false)
}
