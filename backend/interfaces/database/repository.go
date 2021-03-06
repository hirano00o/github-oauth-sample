package database

import "time"

// OauthRepository is ...
type OauthRepository struct {
	DBHandler
}

// StoreState is ...
func (r *OauthRepository) StoreState(state, session string, expiry time.Time) (err error) {
	statement := `INSERT INTO users (session_id, state, expiry) VALUES (?, ?, ?)`
	_, err = r.Execute(statement, session, state, expiry.Format("1990-01-01 01:00:00"))
	if err != nil {
		return
	}
	return
}

// StoreGithubToken is ...
func (r *OauthRepository) StoreGithubToken(tk, tp, reftk string, exp time.Time) (id int, err error) {
	statement := `INSERT INTO github_tokens (token, type, refresh_token, expiry) VALUES (?, ?, ?, ?)`
	res, err := r.Execute(statement, tk, tp, reftk, exp.Format("1990-01-01 01:00:00"))
	if err != nil {
		return
	}
	id64, err := res.LastInsertedId()
	if err != nil {
		return
	}
	id = int(id64)
	return
}

// StoreUserToken is ...
func (r *OauthRepository) StoreUserToken(session, token string, expiry time.Time, tkid int) (id int, err error) {
	statement := `UPDATE users SET token = ?, expiry = ?, github_tokens_id = ? WHERE session_id = ?`
	res, err := r.Execute(statement, token, expiry.Format("1990-01-01 01:00:00"), tkid, session)
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

// FindBySessionID is ...
func (r *OauthRepository) FindBySessionID(session string) (state string, err error) {
	statement := `SELECT state from users where session_id = ?`
	row, err := r.Query(statement, session)
	defer row.Close()
	if err != nil {
		return
	}
	row.Next()
	if err = row.Scan(&state); err != nil {
		return
	}
	return
}

// FindBySessionIDAndUserToken is ...
func (r *OauthRepository) FindBySessionIDAndUserToken(session, token string) (expiry time.Time, id int, err error) {
	statement := `SELECT expiry, github_tokens_id  from users where session_id = ? AND token = ?`
	row, err := r.Query(statement, session, token)
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

// FindByUserTokenID is ...
func (r *OauthRepository) FindByUserTokenID(tokenID int) (acsTk, tp, refTk string, exp time.Time, err error) {
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
