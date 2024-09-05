package server

import (
    "github.com/Gewinum/go-df-discord/utils"
    "log/slog"
    "os"
)

type Opts struct {
    Logger  *slog.Logger
    Repo    Repository
    CodeStr CodeStore
}

func FillEmptyOpts(opts *Opts) {
    if opts.Logger == nil {
        opts.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
    }

    if opts.Repo == nil {
        repo, err := NewDefaultRepository()
        utils.ErrorPanic(err)
        opts.Repo = repo
    }

    if opts.CodeStr == nil {
        opts.CodeStr = newDefaultCodeStore()
    }
}

func DefaultOpts() *Opts {
    opts := &Opts{}
    FillEmptyOpts(opts)
    return opts
}
