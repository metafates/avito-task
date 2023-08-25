package api

import "context"

func (a *API) createUser(ctx context.Context, id UserID) error {
	sql, args, err := a.
		psql().
		Insert("users").
		Columns("id").
		Values(id).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.pg().Exec(ctx, sql, args...)
	return err
}
