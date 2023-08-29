package server

import (
	"bytes"
	"context"
	"encoding/csv"
	"slices"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metafates/avito-task/db"
	"github.com/metafates/avito-task/server/api"
)

type auditAction string

const (
	auditActionAssign  auditAction = "ASSIGN"
	auditActionDeprive auditAction = "DEPRIVE"
)

type auditEntry struct {
	Stamp       time.Time
	UserID      api.UserID
	SegmentSlug api.Slug
	Action      auditAction
}

type audit []auditEntry

func (a audit) CSV() (string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	columns := []string{"user_id", "segment_slug", "action", "stamp"}
	if err := writer.Write(columns); err != nil {
		return "", err
	}

	for _, entry := range a {
		record := []string{
			strconv.Itoa(int(entry.UserID)),
			entry.SegmentSlug,
			string(entry.Action),
			entry.Stamp.Format(time.RFC3339), // as defined by openapi format date-time which uses RFC3339
		}

		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type auditFilter struct {
	From, To *time.Time
	User     *api.UserID
}

func (s *Server) audit(ctx context.Context, conn *pgxpool.Conn, filter auditFilter) (audit, error) {
	query := s.
		psql().
		Select("user_id", "segment_slug", "action", "stamp", "expires_at").
		From(db.TableAssignmentsAudit)

	if filter.From != nil {
		query = query.Where(squirrel.GtOrEq{"stamp": filter.From})
	}

	if filter.To != nil {
		query = query.Where(squirrel.LtOrEq{"stamp": filter.To})
	}

	if filter.User != nil {
		query = query.Where(squirrel.Eq{"user_id": filter.User})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	type auxAuditEntry struct {
		auditEntry

		expiresAt *api.Timestamp
	}

	var auxAudit []auxAuditEntry
	for rows.Next() {
		var entry auxAuditEntry

		if err = rows.Scan(
			&entry.UserID,
			&entry.SegmentSlug,
			&entry.Action,
			&entry.Stamp,
			&entry.expiresAt,
		); err != nil {
			return nil, err
		}

		auxAudit = append(auxAudit, entry)
	}

	var audit audit

	for _, auxEntry := range auxAudit {
		audit = append(audit, auxEntry.auditEntry)

		expiresAt := auxEntry.expiresAt
		if expiresAt == nil || expiresAt.After(time.Now()) {
			continue
		}

		// mark expired segment as deprivation

		// TODO: should be handled from the sql
		if filter.To != nil && expiresAt.After(*filter.To) {
			continue
		}

		audit = append(audit, auditEntry{
			Stamp:       *expiresAt,
			UserID:      auxEntry.UserID,
			SegmentSlug: auxEntry.SegmentSlug,
			Action:      auditActionDeprive,
		})
	}

	// sort by time.
	slices.SortFunc(audit, func(a, b auditEntry) int {
		if a.Stamp.After(b.Stamp) {
			return 1
		}

		if a.Stamp.Equal(b.Stamp) {
			return 0
		}

		return -1
	})

	return audit, nil
}
