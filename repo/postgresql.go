package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ervitis/spamtoputocorreos/models"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

type (
	postgresql struct {
		db  *sql.DB
		cfg *DBParameters
	}
)

func New(params *DBParameters) IRepository {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s", params.Host, params.Port, params.User, params.Pass, params.Name, params.Options)
	psqlConn = strings.TrimSpace(psqlConn)

	conn, err := sql.Open("postgres", psqlConn)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		panic(err)
	}

	return &postgresql{
		db:  conn,
		cfg: params,
	}
}

func (p *postgresql) Save(ctx context.Context, statusTrace *models.StatusTrace) error {
	query := fmt.Sprintf(`INSERT INTO %s (refCode, status, detail, date) VALUES ($1, $2, $3, $4)`, p.cfg.TableName)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	for _, v := range statusTrace.Statuses {
		_, err = stmt.ExecContext(ctx, statusTrace.RefCode, v.Status, v.Detail, v.Date)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresql) Get(ctx context.Context, refID string) (*models.StatusTrace, error) {
	query := fmt.Sprintf(`SELECT status, detail, date FROM %s WHERE refCode = $1`, p.cfg.TableName)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, refID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	statusTraces := &models.StatusTrace{RefCode: refID, Statuses: make([]*models.StatusData, 0)}
	for rows.Next() {
		var statusTrace models.StatusData
		if err := rows.Scan(&statusTrace.Status, &statusTrace.Detail, &statusTrace.Date); err != nil {
			return nil, err
		}
		statusTraces.Statuses = append(statusTraces.Statuses, &statusTrace)
	}

	return statusTraces, nil
}

func (p *postgresql) Delete(ctx context.Context) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, p.cfg.TableName)

	_, err := p.db.ExecContext(ctx, query)
	return err
}
