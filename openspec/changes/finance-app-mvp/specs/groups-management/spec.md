## ADDED Requirements

### Requirement: Create group
The system SHALL allow any authenticated user to create a budget group/workspace.

#### Scenario: Create a new group
- **WHEN** an authenticated user creates a new group with a name
- **THEN** the system creates the group and sets the creator as the owner

### Requirement: Invite members
The system SHALL allow group members to invite other users to the group via email.

#### Scenario: Invite a user
- **WHEN** a group member sends an invitation to an email address
- **THEN** the system sends an invitation email and adds a pending membership

#### Scenario: Accept invitation
- **WHEN** an invited user accepts the invitation
- **THEN** the system adds the user as a full member of the group

### Requirement: Shared editing
All group members SHALL have edit permissions on the shared budget.

#### Scenario: Member edits budget
- **WHEN** any group member modifies a transaction or budget allocation
- **THEN** the change is visible to all other group members after sync

### Requirement: Group membership list
The system SHALL display all members of a group.

#### Scenario: View members
- **WHEN** a group member views the group settings
- **THEN** the system shows a list of all members with their roles
