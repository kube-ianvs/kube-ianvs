package oidc

// CodeFetcher comment lint rebel
type CodeFetcher interface {
	Fetch(a *App) (code string)
}
