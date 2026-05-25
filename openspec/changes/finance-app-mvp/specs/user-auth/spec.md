## ADDED Requirements

### Requirement: User registration
The system SHALL allow users to register with email and password.

#### Scenario: Successful registration
- **WHEN** a user submits a registration form with email and password
- **THEN** the system creates a new user account and returns a success response

#### Scenario: Duplicate email registration
- **WHEN** a user tries to register with an email that already exists
- **THEN** the system returns an error indicating the email is already in use

### Requirement: User authentication
The system SHALL authenticate users via email and password, maintaining a session.

#### Scenario: Successful login
- **WHEN** a user submits valid email and password
- **THEN** the system creates a session and redirects to the dashboard

#### Scenario: Invalid credentials
- **WHEN** a user submits invalid email or password
- **THEN** the system returns an authentication error

### Requirement: Password hashing
The system SHALL store passwords hashed using Django's default PBKDF2 algorithm.

#### Scenario: Password storage
- **WHEN** a user registers or changes password
- **THEN** the system stores only the hashed version of the password

### Requirement: Session management
The system SHALL manage user sessions securely with configurable timeout.

#### Scenario: Session timeout
- **WHEN** a user is inactive for the configured session duration
- **THEN** the system invalidates the session and requires re-login

#### Scenario: Logout
- **WHEN** a user clicks logout
- **THEN** the system destroys the session
