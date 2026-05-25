## ADDED Requirements

### Requirement: Manual sync trigger
The system SHALL allow users to manually trigger synchronization via a pull-to-refresh action.

#### Scenario: Pull-to-refresh sync
- **WHEN** a user clicks the sync button or pulls to refresh
- **THEN** the system sends local changes to the server and pulls remote changes

### Requirement: Sync protocol
The system SHALL track changes using UUIDs and `updated_at` timestamps per entity.

#### Scenario: Sync local changes
- **WHEN** a user creates or modifies a transaction offline
- **THEN** the system records the change with an `updated_at` timestamp and sync status as "pending"

#### Scenario: Push pending changes
- **WHEN** sync is triggered
- **THEN** the system pushes all entities with "pending" status to the server

### Requirement: Conflict resolution
The system SHALL resolve sync conflicts using last-write-wins strategy.

#### Scenario: Last-write-wins
- **WHEN** two users modify the same entity before syncing
- **THEN** the entity with the most recent `updated_at` timestamp is kept

### Requirement: Offline support
The system SHALL allow creating and editing transactions while offline, queuing them for the next sync.

#### Scenario: Offline transaction creation
- **WHEN** a user creates a transaction while offline
- **THEN** the system saves it locally with a "pending" sync status and available in the local UI

#### Scenario: Sync after reconnect
- **WHEN** the user triggers sync after being offline
- **THEN** the system pushes queued changes and pulls remote updates
