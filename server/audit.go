package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
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
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	records := make([][]string, len(a)+1)
	records = append(records, []string{"user_id", "segment_slug", "action", "stamp"})

	for _, entry := range a {
		records = append(records, []string{
			entry.UserID,
			entry.SegmentSlug,
			string(entry.Action),
			entry.Stamp.Format(time.RFC3339), // as defined by openapi format date-time which uses RFC3339
		})
	}

	for _, record := range records {
		fmt.Printf("record: %v\n", record)
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

// TODO: implement, use time range
func (s *Server) audit(ctx context.Context, from *time.Time, to *time.Time) (audit, error) {
	query := s.
		psql().
		Select("user_id", "segment_slug", "action", "stamp", "expires_at").
		From(db.TableAssignmentsAudit).
		JoinClause(fmt.Sprintf("natural join %s", db.TableAssignedSegments))

	if from != nil && to != nil {
		query = query.Where(squirrel.And{
			squirrel.GtOrEq{"stamp": from},
			squirrel.LtOrEq{"stamp": to},
		})
	} else if from != nil {
		query = query.Where(squirrel.GtOrEq{"stamp": from})
	} else if to != nil {
		query = query.Where(squirrel.LtOrEq{"stamp": from})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.pg().Query(ctx, sql, args...)
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

		// TODO: should be handled from the sql
		if to != nil && expiresAt.After(*to) {
			continue
		}

		audit = append(audit, auditEntry{
			Stamp:       *expiresAt,
			UserID:      auxEntry.UserID,
			SegmentSlug: auxEntry.SegmentSlug,
			Action:      auditActionDeprive,
		})
	}

	// sort by time
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
