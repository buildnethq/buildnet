# BuildNET Protocol v5 (Draft)

## Goals
- End-to-end encrypted messaging across heterogeneous transports.
- Suite agility with a safe default and policy-gated stronger suites.
- Multiplexed bidirectional channels (control/data/event/telemetry).
- Tamper-evident audit ledger with double-entry linking.
- Federation + controller mobility.

## Non-goals / Limits
- Perfect anonymity on the public internet is not guaranteed.
- Timing metadata can be reduced but not eliminated; FlowPolicy describes tradeoffs.
- No shared-memory SMP semantics across hosts.

## Design
- Underlay: QUIC preferred, HTTP/2 fallback.
- Overlay: application-layer E2E session + SecureEnvelope framing.
- Audit: append-only signed events with hash chaining + optional double-entry.
- Resources: generic ResourceDescriptor catalog with event updates and visibility policy.

## Next
- Implement TransportService with real crypto core (Noise-style XX/IK).
- Add SQLite persistence for Ledger and Resources.
- Add CLI: resources ls/watch, ledger tail.
