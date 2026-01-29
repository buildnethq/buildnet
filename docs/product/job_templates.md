# BuildNET JobTemplates (v1)

JobTemplates compile into Action DAGs. Each template declares:
- Inputs (CAS trees/blobs, repo@rev)
- Environment (toolchain capsule)
- Effects policy (default deterministic)
- Constraints (required resources/tools)
- Outputs (artifacts + manifests + attestations)

## Template set (initial)

### 1) ReleaseBundle
Goal: signed multi-target bundle (+ SBOM + provenance).
Actions: FetchRepo → ResolveDeps → Build(target_i)* → Test(shards)* → SBOM → Sign → Bundle

### 2) CompileRepo
Goal: sharded compile/test with caching.
Actions: FetchRepo → Compile(shards)* → Link → Test(shards)*

### 3) ArchivePack
Goal: huge archive beyond single-host capacity (streaming).
Actions: PlanTar → TarShard(shards)* → CompressShard(shards)* → Assemble → Verify

### 4) TrustBuild
Goal: redundant execution integrity (k-of-n).
Actions: Execute(k diverse workers) → QuorumCheck → Attest

### 5) IndexBuild
Goal: symbols/search/docs indexes.
Actions: Parse(shards)* → BuildIndex → Merge

### 6) SBOMProvenance
Goal: SBOM + provenance without full rebuild.
Actions: ScanDeps → SBOM → Provenance → Sign

### 7) DataTransform (ETL)
Goal: deterministic map/reduce over immutable inputs.
Actions: Map(shards)* → Reduce → Verify

### 8) ContainerBuild
Goal: OCI image builds with layer caching.
Actions: BuildLayer_i(cacheable) → AssembleImage → Sign

### 9) CodecTranscode
Goal: chunked encode/decode where semantics permit.
Actions: Chunk → Encode(shards)* → Assemble → Verify

### 10) CacheWarm
Goal: pre-populate CAS across a fleet.
Actions: Pull → Replicate → Verify
