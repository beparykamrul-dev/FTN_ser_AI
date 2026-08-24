# FTN Core Router / MikroTik / Android

FTN Phase-1 exposes a vendor-neutral router registry with a first-class MikroTik profile and an Android client contract.

## Core router

The control plane tracks router identity, vendor, address, API version, authorization, enabled state and capabilities. Runtime adapters are responsible for actual device operations.

## MikroTik

MikroTik RouterOS is represented as an authorized adapter target. Credentials and certificates are runtime secrets and are never stored in the repository.

## Android

The Android client contract supports:

- FTNVPN
- authorized router-control workflows
- service-health visibility
- network diagnostics

Android does not receive unrestricted router credentials. Device identity and authorization are required before control-plane operations are exposed.

## Production path

`Android -> FTN API -> authorization -> router adapter -> MikroTik/FTN router -> health -> telemetry -> FTN Silk`

The registry does not itself execute router changes; deployment/runtime adapters apply approved operations and report the resulting state.
