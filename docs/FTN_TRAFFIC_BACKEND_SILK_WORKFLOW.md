# FTN Traffic, Backend and SiLK Workflow

```text
DISCOVER
  -> REGISTER
  -> HEALTH/CAPACITY
  -> POLICY
  -> ROUTE/GATEWAY DECISION
  -> FLOW TELEMETRY
  -> SiLK/ANALYTICS
  -> ClickHouse
  -> AI ANALYSIS
  -> RECOMMENDATION
  -> AUTHORIZED APPLY
  -> VERIFY
  -> REBALANCE/ROLLBACK
```

## Data boundaries

FTN AI consumes authorized FTN telemetry and service state. It does not infer or access private provider databases or private traffic without an explicit authorized interface.

## Backend integration

Service AI can coordinate PostgreSQL, CockroachDB, TimescaleDB, ClickHouse, Kafka and Flink through versioned adapters. Long-running changes are idempotent and emit audit events.

## Routing

The workflow may prepare GoBGP/BFD/ECMP changes only from an authorized peer/prefix policy. PKI/mTLS protects control-plane communication; BGP authentication is selected according to router/provider capability.

## Failover

A node or route becomes ineligible after failed health checks. Recovery requires successful reconciliation and health verification before returning it to the eligible pool.
