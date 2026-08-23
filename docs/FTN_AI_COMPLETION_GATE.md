# FTN AI Completion Gate

FTN AI is complete only when its decisions are grounded in authorized telemetry and service state and every automated change is policy-controlled, auditable and verifiable.

## Pipeline

`OBSERVE -> VALIDATE -> ANALYZE -> PROPOSE -> AUTHORIZE -> APPLY -> VERIFY -> ROLLBACK/ACCEPT`

## Supported data domains

- node/service health
- DNS health
- routing/BGP/BFD state
- gateway capacity
- IPFIX/NetFlow flow metadata
- TimescaleDB metrics
- ClickHouse analytics
- PostgreSQL service state
- CockroachDB distributed state
- Kafka/Flink streams
- certificate lifecycle

## Safety

AI must not invent credentials, peer authorization, ASN ownership, prefixes, private CA keys or access rights. Provider integrations are restricted to public, documented or explicitly authorized interfaces.

## Completion criteria

- deterministic workflow contracts
- idempotent jobs
- audit events
- retry/rollback
- health verification
- human approval where required
- CI/test coverage for critical workflows
