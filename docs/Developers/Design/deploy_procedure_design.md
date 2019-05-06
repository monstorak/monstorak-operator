This document describes the actions that have to be executed successfully for
monstorak's CRDs to be considered reconciled. They are implemented as
`reconciler.Action` objects, and are enumerated in a `reconciler.Procedure`
object. Procedure level actions are ones that modify state and should have a
corresponding entry populated in
`.Status.ReconcileActions map[string]reconciler.Result`.
Top level actions are executed in an arbitrary order so they must define any
prerequisite actions explicitly.
An action may be a top-level action and still defined as a prerequisite and the
caching implementation will ensure that it is executed a maximum of once per
Procedure execution.

# StorageMixinClass actions

## clusterMonitoringDeployed

## storageClusterDeployed

## monitoringEnabledForStorage

## deployStorageMixin

