# BuildNET

BuildNET is a federated distributed execution + caching system for build workloads (compile/test shards), compression pipelines, and other coarse-grained, mostly-immutable data transforms.

**Core concepts**
- **BuildNet**: the network (hosts become buildbots).
- **Buildbot**: an agent that contributes resources and executes work.
- **BuildVM**: the virtual execution machine (actions + CAS + scheduling + audit), not pooled RAM.
- **BuildNET Protocol**: Protobuf + ConnectRPC API surface.
- **BuildNET Language (BNL)**: a structured command language compiled into protocol calls (CLI + shell + UI parity).

**Status**
- Version line: **5.x** (major API namespace: `buildnet.v5`)
- Current: early alpha scaffolding / architecture lock-in.

## Goals
- Remote-first: run builds from anywhere.
- End-to-end encryption, capability-based enrollment, full audit trail.
- Movable control plane (leader is a role, not a host).
- Worker classes (FULL/MOBILE/EDGE/AI) with governance + performance views.
- Deterministic actions, content-addressed storage, cache-aware scheduling.

## Quickstart (placeholder)
Coming soon:
- `buildnet version`
- `buildnet hub start`
- `buildnet worker join ...`
- `buildnet workers ls`

## License
Dual-licensed under **MIT OR Apache-2.0** (see `LICENSE`).
