# FTN Control Plane additions

This branch adds the control-plane foundation without replacing the existing FTN SER AI state API.

Modules: monitoring, runtime telemetry, network/storage telemetry models, database health model, DNS mesh and Anycast state, provider abstraction, IPAM types, jobs, audit events, access decisions, latency probes, and a unified control snapshot.

Zones represented by the mesh model: `familytimenet.com` and `ftnddns.net`.

Anycast state is declarative; route activation must come from real node/BGP telemetry and must not be represented as active merely because configuration exists.
