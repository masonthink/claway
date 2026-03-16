# Claway TODOs

## From CEO Review (2026-03-14)

### Pending Implementation (decided but not yet built)

- [ ] **5A: JWT → httpOnly cookie** — Migrate token from localStorage to httpOnly cookie. Backend Set-Cookie on auth, frontend stops touching token directly. Effort: M, Priority: P0.
- [ ] **3A: Sentry integration** — Frontend + backend error tracking. Next.js SDK + Go Sentry middleware. Effort: S, Priority: P1.
- [ ] **10A: zerolog + /metrics** — Replace Echo default logger with zerolog structured logging. Add Prometheus-compatible /metrics endpoint. Effort: M, Priority: P1.
- [ ] **9A: Handler + middleware tests** — Add tests for auth middleware (JWT validation, expiry, malformed tokens) and all HTTP handlers (request parsing, auth checks, error responses). Effort: M, Priority: P1.
- [ ] **2A: Auth sessions → PostgreSQL** — Migration 007 created the auth_sessions table. Still need to update store/auth_session.go to use DB instead of sync.Map. Effort: S, Priority: P2.

### Vision / Delight (future)

- [ ] **D1: Reveal countdown animation** — Last 24h: pulsing countdown timer + urgency copy ("Revealing in 3h 24m"). Effort: XS, Priority: P3.
- [ ] **D4: PDF export** — One-click export featured proposals as branded PDF. Needs PDF generation library. Effort: M, Priority: P3.
