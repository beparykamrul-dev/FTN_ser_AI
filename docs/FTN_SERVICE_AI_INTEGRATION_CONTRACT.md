# FTN Service AI Integration Contract

## Purpose

FTN Service AI is the service-intelligence layer for the FTN control plane. It observes authorized service/node telemetry, evaluates policy, produces recommendations, and can prepare configuration workflows. Autonomous changes must pass policy, authorization and health validation.

## Mesh boundaries

- FTNDNS Global Mesh remains responsible for DNS.
- FTN Server/Service Mesh remains responsible for server and application workloads.
- Service AI coordinates between them through versioned contracts; it does not bypass routing, PKI, or gateway policy.

## Inputs

- Node health and capacity
- Service health
- Authorized NetFlow/IPFIX/sFlow metadata
- DNS health and latency
- Database/queue health
- Certificate lifecycle state
- BGP session and route telemetry
- Storage capacity and health
- Customer/service SLO metrics

## Decision loop

`OBSERVE -> VALIDATE -> ANALYZE -> PROPOSE -> AUTHORIZE -> APPLY -> VERIFY -> ROLLBACK/ACCEPT`

## Safety boundaries

AI must never invent provider credentials, peer authorization, ASN ownership, prefixes, CA private keys, or access rights. BGP peers and route advertisements require explicit policy. Provider integrations use documented or authorized interfaces only.

## Data destinations

- TimescaleDB for time-series service metrics
- ClickHouse for high-volume logs/flow analytics
- PostgreSQL for service/business state
- CockroachDB where distributed control state is required
- Kafka/Flink for streaming workloads

## Automation targets

The service workflow may prepare and validate:

- service placement and rebalancing
- gateway selection
- health-based failover
- certificate renewal workflow
- DNS endpoint health state
- telemetry registration
- capacity alerts
- deployment/rollback plans

The workflow must remain idempotent and auditable.
