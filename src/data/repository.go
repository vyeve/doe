package data

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"doe/src/logger"
	"doe/src/models"

	"github.com/lib/pq"
	"go.uber.org/fx"
)

type Repository interface {
	Insert(ctx context.Context, ports []*models.Port) error
	GetAll(ctx context.Context, limit int32) ([]*models.Port, error)
	GetByID(ctx context.Context, portID string) (*models.Port, error)
}

type repoImpl struct {
	tx  SQLDb
	log logger.Logger
}

func NewRepository(params Params) Repository {
	repo := &repoImpl{
		tx:  params.DBConn,
		log: params.Logger,
	}
	params.LifeCycle.Append(fx.Hook{
		OnStart: func(context.Context) error { return nil },
		OnStop: func(ctx context.Context) error {
			return repo.tx.Close()
		},
	})
	return repo
}

func (r *repoImpl) guardPanic(err *error) {
	if p := recover(); p != nil {
		r.log.Errorf("caught panic: %v", p)
		*err = fmt.Errorf("caught panic: %v. err: %w", p, *err)
	}
}

func (r *repoImpl) Insert(ctx context.Context, ports []*models.Port) (err error) {
	defer r.guardPanic(&err)
	r.log.Debugf("Try to insert %d ports", len(ports))
	if len(ports) == 0 {
		return nil
	}
	return ExecuteTx(ctx, r.tx.DBConn(), func(tx SQLTx) error {
		ph := strings.Builder{}
		phCleanup := strings.Builder{}
		phAlias := strings.Builder{}
		values := make([]interface{}, 0)
		numFields := 10
		aliases := make([]interface{}, 0)
		k := 0
		portIDs := make([]interface{}, 0, len(ports))
		for i, p := range ports {
			values = append(values,
				p.PortID,
				p.Name,
				p.City,
				p.Country,
				pq.Array(p.Regions),
				pq.Array(p.Coordinates),
				p.Province,
				p.Timezone,
				pq.Array(p.Unlocs),
				p.Code,
			)
			portIDs = append(portIDs, p.PortID)
			if len(p.Alias) > 0 {
				for _, a := range p.Alias {
					if k != 0 {
						phAlias.WriteRune(',')
					}
					phAlias.WriteString(fmt.Sprintf("($%d,$%d)", k*2+1, k*2+2))
					k++
					aliases = append(aliases, p.PortID, a)
				}

			}
			if i != 0 {
				ph.WriteRune(',')
				phCleanup.WriteRune(',')
			}
			phCleanup.WriteString("$" + strconv.Itoa(i+1))
			ph.WriteRune('(')
			for j := 1; j <= numFields; j++ {
				if j != 1 {
					ph.WriteRune(',')
				}
				ph.WriteString("$" + strconv.Itoa(j+i*numFields))
			}
			ph.WriteRune(')')
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf(insertStmt, ph.String()), values...)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf(cleanupAliasesStm, phCleanup.String()), portIDs...)
		if err != nil {
			return err
		}
		if len(aliases) == 0 {
			return nil
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf(insertAliasesStmt, phAlias.String()), aliases...)
		if err != nil {
			return err
		}
		return nil

	})
}

func (r *repoImpl) GetAll(ctx context.Context, limit int32) (_ []*models.Port, err error) {
	defer r.guardPanic(&err)
	result := make([]*models.Port, 0, limit)
	var count int32
	for count != limit {
		rows, err := r.tx.QueryContext(ctx, listAllStmt, limit)
		if err != nil {
			r.log.Warnf("Failed to extract ports. err: %v", err)
			return nil, err
		}
		for ; rows.Next(); count++ {
			var (
				p       models.Port
				aliases []sql.NullString
			)

			if err = rows.Scan(
				&p.PortID,
				&p.Name,
				&p.City,
				&p.Province,
				&p.Country,
				pq.Array(&p.Regions),
				pq.Array(&p.Coordinates),
				&p.Timezone,
				pq.Array(&p.Unlocs),
				&p.Code,
				pq.Array(&aliases),
			); err != nil {
				r.log.Warnf("Failed to scan data. err: %v", err)
				return nil, err
			}
			for _, a := range aliases {
				p.Alias = append(p.Alias, a.String)
			}
			result = append(result, &p)
		}
	}
	return result, nil
}

func (r *repoImpl) GetByID(ctx context.Context, portID string) (_ *models.Port, err error) {
	defer r.guardPanic(&err)
	var (
		p       models.Port
		aliases []sql.NullString
	)
	err = r.tx.QueryRowContext(ctx, getByIDStmt, portID).Scan(
		&p.PortID,
		&p.Name,
		&p.City,
		&p.Province,
		&p.Country,
		pq.Array(&p.Regions),
		pq.Array(&p.Coordinates),
		&p.Timezone,
		pq.Array(&p.Unlocs),
		&p.Code,
		pq.Array(&aliases),
	)
	if err != nil {
		r.log.Errorf("Failed to get port. err: %v", err)
		return nil, err
	}
	for _, a := range aliases {
		p.Alias = append(p.Alias, a.String)
	}
	return &p, nil
}
