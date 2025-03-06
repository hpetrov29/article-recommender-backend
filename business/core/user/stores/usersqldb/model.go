package usersqldb

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/hpetrov29/resttemplate/business/core/user"
	"github.com/hpetrov29/resttemplate/business/data/dbsql/mysql/dbarray"
)

// dbUser represents the structure used to transfer user data
// between the application and the database.
type dbUser struct {
	Id           	int64      		`db:"id"`
	Username     	string      		`db:"username"`
	Email        	string         		`db:"email"`
	Roles        	dbarray.String 		`db:"roles"`
	PasswordHash 	[]byte         		`db:"password_hash"`
	CreatedAt  		time.Time      		`db:"created_at"`
}

// toDBUser converts a user.User instance (found in the service layer) to a dbUser struct suited for database operations.
//
// Parameters:
//   - usr: the user instance to be converted.
func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return dbUser{
		Id:           usr.Id,
		Username:     usr.Username,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		CreatedAt: 	  usr.CreatedAt.UTC(),
	}
}

// toCoreUser converts a dbUser instance (found in the repository layer) to a user.User struct.
//
// Parameters:
//   - dbUsr: the dbUsr instance to be converted.
func toCoreUser(dbUsr dbUser) (user.User, error) {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.User{}, fmt.Errorf("parse role: %w", err)
		}
	}

	usr := user.User{
		Id:           dbUsr.Id,
		Username:     dbUsr.Username,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		CreatedAt:    dbUsr.CreatedAt.In(time.Local),
	}

	return usr, nil
}