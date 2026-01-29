# BuildVM Core (Design Notes)

BuildVM is not a distributed process VM. It is an Action-DAG evaluation engine:
- compiles user intents (JobTemplates) into Action graphs
- executes actions on a fleet with caching, retries, auditing, and policy enforcement
- produces canonical artifacts (digests) with provenance

## Core model
- Digest: BLAKE3
- CAS objects: Blob, Tree, Manifest
- ActionSpec: inputs + env_digest + command + effects + constraints + limits
- ActionDigest: H(ActionSpec + input digests + env_digest + policy)
- ActionResult: output digests + logs + metrics + taint flags

## Evaluation loop
1) compute ActionDigest
2) cache lookup; on hit emit ledger event and short-circuit
3) select worker by constraints + normalized compute + mem window + RTT + locality
4) dispatch, stream logs
5) store outputs in CAS; store ActionDigest→ResultDigest mapping
6) ledger scheduled/started/finished (+ optional double-entry links)

## Executors (pluggable)
v1 local runner → v2 container → v3 WASI/WASM → v4 microVM
