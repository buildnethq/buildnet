# Contributing to BuildNET

## Ground rules
- Keep changes focused and reviewable.
- Prefer additive protobuf changes; never reuse field numbers.
- Keep CLI/GUI parity by routing everything through the same API surface.

## Development workflow
1. Fork + branch from `main`
2. Open a PR
3. Ensure tests/build pass

## Commit style
Use clear, imperative messages:
- "Add FleetService Join/Heartbeat"
- "Implement ActionDigest calculation"
