# Decision 002: CI/CD Strategy

**Date:** 2026-03-15 | **Decider:** Alex (CEO)

## Decision
- GitHub Actions for CI (test + lint on push/PR)
- GoReleaser for binary releases on tag push
- Support Go 1.24 and 1.25 in CI matrix
- Platforms: Linux + macOS, amd64 + arm64

## Rationale
- GitHub Actions: Already on GitHub, free for public repos
- GoReleaser: Standard for Go projects, handles cross-compilation + checksums
- Multi-version testing catches compatibility issues early
