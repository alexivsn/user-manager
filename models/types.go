package models

import (
	"time"

	"io"

	"context"

	"git.containerum.net/ch/json-types/errors"
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"github.com/lib/pq"
)

// User describes user model. It should be used only inside this project.
type User struct {
	ID            string           `db:"id"`
	Login         string           `db:"login"`
	PasswordHash  string           `db:"password_hash"` // base64
	Salt          string           `db:"salt"`          // base64
	Role          umtypes.UserRole `db:"role"`
	IsActive      bool             `db:"is_active"`
	IsDeleted     bool             `db:"is_deleted"`
	IsInBlacklist bool             `db:"is_in_blacklist"`
}

// Profile describes user`s profile model. It should be used only inside this project.
type Profile struct {
	ID          string
	Referral    string
	Access      string
	CreatedAt   time.Time
	BlacklistAt pq.NullTime
	DeletedAt   pq.NullTime

	User *User

	Data map[string]interface{}
}

// Accounts describes user`s bound accounts. It should be used only inside this project.
type Accounts struct {
	ID       string
	Github   string
	Facebook string
	Google   string

	User *User
}

// Link describes link (for activation, password change, etc.) model. It should be used only inside this project.
type Link struct {
	Link      string
	Type      umtypes.LinkType
	CreatedAt time.Time
	ExpiredAt time.Time
	IsActive  bool
	SentAt    pq.NullTime

	User *User
}

// Token describes one-time token model (used for login), It should be used only inside this project.
type Token struct {
	Token     string
	CreatedAt time.Time
	IsActive  bool
	SessionID string

	User *User
}

// DomainBlacklistEntry describes one blacklisted email domain.
// Registration with email for this domain must be rejected. It should be used only inside this project.
type DomainBlacklistEntry struct {
	Domain    string
	CreatedAt time.Time
}

// Errors which may occur in transactional operations
var (
	ErrTransactionBegin    = errors.New("transaction begin error")
	ErrTransactionRollback = errors.New("transaction rollback error")
	ErrTransactionCommit   = errors.New("transaction commit error")
)

// DB is an interface for persistent data storage (also sometimes called DAO).
type DB interface {
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetBlacklistedUsers(ctx context.Context, limit, offset int) ([]User, error)
	BlacklistUser(ctx context.Context, user *User) error

	CreateProfile(ctx context.Context, profile *Profile) error
	GetProfileByID(ctx context.Context, id string) (*Profile, error)
	GetProfileByUser(ctx context.Context, user *User) (*Profile, error)
	UpdateProfile(ctx context.Context, profile *Profile) error
	GetAllProfiles(ctx context.Context, perPage, offset int) ([]Profile, error)

	GetUserByBoundAccount(ctx context.Context, service umtypes.OAuthResource, accountID string) (*User, error)
	GetUserBoundAccounts(ctx context.Context, user *User) (*Accounts, error)
	BindAccount(ctx context.Context, user *User, service umtypes.OAuthResource, accountID string) error

	BlacklistDomain(ctx context.Context, domain string) error
	UnBlacklistDomain(ctx context.Context, domain string) error
	IsDomainBlacklisted(ctx context.Context, domain string) (bool, error)
	/*GetBlacklistedDomains() ([]DomainBlacklistEntry, error)*/

	CreateLink(ctx context.Context, linkType umtypes.LinkType, lifeTime time.Duration, user *User) (*Link, error)
	GetLinkForUser(ctx context.Context, linkType umtypes.LinkType, user *User) (*Link, error)
	GetLinkFromString(ctx context.Context, strLink string) (*Link, error)
	UpdateLink(ctx context.Context, link *Link) error
	GetUserLinks(ctx context.Context, user *User) ([]Link, error)

	GetTokenObject(ctx context.Context, token string) (*Token, error)
	CreateToken(ctx context.Context, user *User, sessionID string) (*Token, error)
	GetTokenBySessionID(ctx context.Context, sessionID string) (*Token, error)
	DeleteToken(ctx context.Context, token string) error
	UpdateToken(ctx context.Context, token *Token) error

	// Perform operations inside transaction
	// Transaction commits if `f` returns nil error, rollbacks and forwards error otherwise
	// May return ErrTransactionBegin if transaction start failed,
	// ErrTransactionCommit if commit failed, ErrTransactionRollback if rollback failed
	Transactional(ctx context.Context, f func(ctx context.Context, tx DB) error) error

	io.Closer
}
