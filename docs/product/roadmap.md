# BuildNET Roadmap (BuildVM-first)

## Vision
BuildNET is a federated, auditable, end-to-end encrypted execution + caching network.
BuildVM is the execution core: it compiles user intents into deterministic, cacheable Action DAGs and schedules shards across heterogeneous buildbots.

## Guiding Principles
- Prefer Action DAG semantics over “distributed process” semantics.
- Content-addressed inputs/outputs (CAS) are the ground truth.
- Determinism is enforced by policy (effects declared, sandboxed).
- Every significant event is ledgered (tamper-evident, double-entry).
- Controller role is mobile; no fixed “hub”.

---

## Phase 0 — Repo & Protocol Wedge (alpha)
### Outcomes
- Protobuf contracts: Fleet + (Transport, Ledger, Resource catalog) stubs.
- Hub boots and serves endpoints; worker can join; basic UI scaffold exists.
- Resources (tools) are discoverable and watchable.

### Milestones
- [ ] TransportService stub → real session (Noise-style handshake later)
- [ ] LedgerService stub → SQLite persistence + append-only semantics
- [ ] ResourceService stub → PATH tool discovery + watch stream
- [ ] CLI: `buildnet resources ls/watch`, `buildnet ledger tail`

---

## Phase 1 — BuildVM Core MVP: Action Engine + CAS (alpha)
### Outcomes
- Canonical ActionSpec model: Inputs + Environment + Command + Effects → Outputs.
- CAS chunking + streaming transfers.
- Scheduler dispatches actions by constraints (os/arch/tools/trust).
- Cache hits across workers by ActionDigest.

### Milestones
- [ ] CAS: Put/Get blobs by digest; trees/manifests (Merkle DAG)
- [ ] Action service: SubmitAction, StreamLogs, GetResult
- [ ] Worker executor v1: local process runner, bounded resources
- [ ] Cache: ActionDigest → ResultDigest persisted in SQLite; optional Dragonfly/Redis

---

## Phase 2 — JobTemplates: “One-click deliverables” (alpha→beta)
### Outcomes
- JobTemplates compile to Action DAGs.
- Flagship deliverables:
  - ReleaseBundle (multi-target artifacts + SBOM + provenance)
  - ArchivePack (distributed “giant tar.gz” via member concatenation + deterministic assembly)
  - TrustBuild (k-of-n redundant execution + quorum)
  - IndexBuild (symbols/search/docs indexes)

### Milestones
- [ ] JobTemplate schema + CLI: `buildnet jobs templates`, `buildnet jobs run`
- [ ] Template compiler: TemplateSpec → Action DAG + constraints
- [ ] UI: DAG viewer + artifact browser + log streaming

---

## Phase 3 — Determinism & Sandboxing (beta)
### Outcomes
- Declared effects enforced (no ambient network/clock/random by default).
- Sandbox backends (pluggable): container, WASI/WASM, microVM later.

### Milestones
- [ ] Effects lattice: Read/Write/Network/Clock/Random budgets
- [ ] “Tainted outputs” for nondeterministic actions
- [ ] Sandbox policy by trust tier and worker class

---

## Phase 4 — Federation + Controller Mobility (beta→stable)
### Outcomes
- Sharded Raft for authoritative state (policy, catalog, schedules).
- Gossip/stream plane for telemetry.
- Controller role transfer with rekeying + ledger continuity.

### Milestones
- [ ] Leader election + role transfer
- [ ] Replicated ledger snapshots
- [ ] Multi-path connectivity (direct + relay)

---

## Phase 5 — Capability Fabric + Marketplace (stable+)
### Outcomes
- Rich Resource catalog with discovery plugins: GPU, toolchains, DB endpoints, AI endpoints (metadata only).
- Capability-aware scheduling across org/community pools.

### Milestones
- [ ] Resource attestation (signed announcements)
- [ ] Cost models, quotas, rate limits
- [ ] “Compute donation pool” support
