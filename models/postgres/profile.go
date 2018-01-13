package postgres

import (
	"github.com/json-iterator/go"

	. "git.containerum.net/ch/user-manager/models"
)

const profileQueryColumnsWithUser = "profiles.id, profiles.referral, profiles.access, profiles.created_at, profiles.blacklist_at, profiles.deleted_at, " +
	"users.id, users.login, users.password_hash, users.salt, users.role, users.is_active, users.is_deleted, users.is_in_blacklist, profiles.data"
const profileQueryColumns = "id, referral, access, created_at, blacklist_at, deleted_at, data"

func (db *PgDB) CreateProfile(profile *Profile) error {
	db.log.Debugln("Create profile for", profile.User.Login)
	profileData, err := jsoniter.MarshalToString(profile.Data)
	if err != nil {
		return err
	}
	rows, err := db.qLog.Queryx("INSERT INTO profiles (referral, access, user_id, data) VALUES "+
		"($1, $2, $3, $4) RETURNING id, created_at", profile.Referral, profile.Access, profile.User.ID, profileData)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}

	err = rows.Scan(&profile.ID, &profile.CreatedAt)
	return err
}

func (db *PgDB) GetProfileByID(id string) (*Profile, error) {
	db.log.Debugln("Get profile by id", id)
	rows, err := db.qLog.Queryx("SELECT "+profileQueryColumnsWithUser+" FROM profiles "+
		"JOIN users ON profiles.user_id = user.id WHERE profiles.id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	profile := Profile{User: &User{}}
	var profileData string
	err = rows.Scan(
		&profile.ID, &profile.Referral, &profile.Access, &profile.CreatedAt, &profile.BlacklistAt, &profile.DeletedAt,
		&profile.User.ID, &profile.User.Login, &profile.User.PasswordHash, &profile.User.Salt, &profile.User.Role,
		&profile.User.IsActive, &profile.User.IsDeleted, &profile.User.IsInBlacklist,
		&profileData,
	)
	if err != nil {
		return nil, err
	}
	if err := jsoniter.UnmarshalFromString(profileData, &profile.Data); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *PgDB) GetProfileByUser(user *User) (*Profile, error) {
	db.log.Debugf("Get profile by user %#v", user)
	rows, err := db.qLog.Queryx("SELECT "+profileQueryColumns+" FROM profiles "+
		"WHERE profiles.user_id = $1", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	profile := Profile{User: user}
	var profileData string

	err = rows.Scan(&profile.ID, &profile.Referral, &profile.Access, &profile.CreatedAt, &profile.BlacklistAt, &profile.DeletedAt, &profileData)
	if err != nil {
		return nil, err
	}
	if err := jsoniter.UnmarshalFromString(profileData, &profile.Data); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *PgDB) UpdateProfile(profile *Profile) error {
	db.log.Debugf("Update profile %#v", profile)
	_, err := db.eLog.Exec("UPDATE profiles SET referal = $2, access = $3, data = '$4 WHERE id = $1",
		profile.ID, profile.Referral, profile.Access, profile.Data)
	return err
}

func (db *PgDB) GetAllProfiles(perPage, offset int) ([]Profile, error) {
	db.log.Debugln("Get all profiles")
	profiles := make([]Profile, 0) // return empty slice instead of nil if no records found

	rows, err := db.qLog.Queryx("SELECT "+profileQueryColumnsWithUser+" FROM profiles JOIN users ON profiles.user_id = users.id "+
		"LIMIT $1 OFFSET $2", perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		profile := Profile{User: &User{}}
		var profileData string
		err := rows.Scan(
			&profile.ID, &profile.Referral, &profile.Access, &profile.CreatedAt, &profile.BlacklistAt, &profile.DeletedAt,
			&profile.User.ID, &profile.User.Login, &profile.User.PasswordHash, &profile.User.Salt, &profile.User.Role,
			&profile.User.IsActive, &profile.User.IsDeleted, &profile.User.IsInBlacklist,
			&profileData,
		)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(profileData, &profile.Data); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, rows.Err()
}
