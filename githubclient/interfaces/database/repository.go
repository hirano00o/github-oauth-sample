package database

import "time"

// OauthRepository shows database handler
type OauthRepository struct {
	DBHandler
}

// StoreState stores user token and login state
func (r *OauthRepository) StoreState(token, state string, expiry time.Time) (err error) {
	statement := `INSERT INTO users (session_id, state, expiry) VALUES (?, ?, ?)`
	_, err = r.Execute(statement, token, state, expiry.Format("2006-01-02 15:04:05"))
	if err != nil {
		return
	}
	return
}

// StoreGithubToken stores github token
func (r *OauthRepository) StoreGithubToken(tk, tp, reftk string, exp time.Time, usersId int) (err error) {
	statement := `INSERT INTO github_tokens (token, type, refresh_token, expiry, users_id) VALUES (?, ?, ?, ?, ?)`
	_, err = r.Execute(statement, tk, tp, reftk, exp.Format("2006-01-02 15:04:05"), usersId)
	if err != nil {
		return
	}
	return
}

// UpdateUserTokenExpiry updates user token expiry
func (r *OauthRepository) UpdateUserTokenExpiry(token string, expiry time.Time) (id int, err error) {
	statement := `UPDATE users SET expiry = ? WHERE token = ?`
	res, err := r.Execute(statement, expiry.Format("1990-01-01 01:00:00"), token)
	if err != nil {
		return
	}
	id64, err := res.RowAffected()
	if err != nil {
		return
	}
	id = int(id64)
	return
}

// FindStateByUserToken is ...
func (r *OauthRepository) FindStateByUserToken(token string) (id int, state string, err error) {
	statement := `SELECT id, state from users where token = ?`
	row, err := r.Query(statement, token)
	defer row.Close()
	if err != nil {
		return
	}
	row.Next()
	if err = row.Scan(&id, &state); err != nil {
		return
	}
	return
}

// FindGithubTokenIDByUserToken is ...
func (r *OauthRepository) FindGithubTokenIDByUserToken(token string) (expiry time.Time, id int, err error) {
	statement := `SELECT expiry, github_tokens_id from users where token = ?`
	row, err := r.Query(statement, token)
	defer row.Close()
	if err != nil {
		return
	}
	row.Next()
	if err = row.Scan(&expiry, &id); err != nil {
		return
	}
	return
}

// FindByUserID is ...
func (r *OauthRepository) FindByUserID(tokenID int) (acsTk, tp, refTk string, exp time.Time, err error) {
	statement := `SELECT token, type, refresh_token, expiry from github_tokens where id = ?`
	row, err := r.Query(statement, tokenID)
	defer row.Close()
	if err != nil {
		return
	}
	row.Next()
	if err = row.Scan(&acsTk, &tp, &refTk, &exp); err != nil {
		return
	}
	return
}
