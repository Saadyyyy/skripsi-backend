package rangkingrepository

import (
	"bank_soal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	queryCreateRangking = `
		insert into rangkings(user_id,category_id,soal_id,point,next,created_at)values($1,$2,$3,$4,$5,$6)returning rangking_id
	`

	queryGetPoint = `
	SELECT COALESCE(SUM((point + point) - point), 0) AS total_points
FROM rangkings
WHERE user_id = $1 AND deleted_at IS NULL;

	`
)

type RangkingRepository interface {
	CreateRangking(ctx context.Context, rank models.Rangking) (id int64, err error)
	GetPointByUserId(ctx context.Context, id int64) (rank models.Rangking, err error)
}

type RangkingRepositoryImpl struct {
	db *sqlx.DB
}

func NewRangkingRepository(db *sqlx.DB) RangkingRepository {
	return &RangkingRepositoryImpl{db: db}
}

func (r *RangkingRepositoryImpl) CreateRangking(ctx context.Context, rank models.Rangking) (id int64, err error) {
	created_at := time.Now()
	err = r.db.QueryRowContext(ctx, queryCreateRangking, rank.UserId, rank.CategoryId, rank.SoalId, rank.Point, rank.Next, created_at).Scan(&id)
	if err != nil {
		err = fmt.Errorf("queryCreateRangking err%+v", err)
		return
	}
	return id, nil

}

// // GetPointByUserId implements RangkingRepository.
// func (r *RangkingRepositoryImpl) GetPointByUserId(ctx context.Context, id int64) (rank models.Rangking, err error) {

// 	err = r.db.QueryRowContext(ctx, queryGetPoint, id).Scan(&rank.Point)
// 	if err != nil {
// 		err = fmt.Errorf("queryGetPoint err%+v", err)
// 		return
// 	}
// 	return rank, nil
// }

func (r *RangkingRepositoryImpl) GetPointByUserId(ctx context.Context, id int64) (models.Rangking, error) {
	var totalPoints sql.NullInt64
	err := r.db.QueryRowContext(ctx, queryGetPoint, id).Scan(&totalPoints)
	if err != nil {
		err = fmt.Errorf("queryGetPoint err%+v", err)
		return models.Rangking{}, err
	}

	// Convert to int64 if valid, otherwise use 0
	var points int64
	if totalPoints.Valid {
		points = totalPoints.Int64
	} else {
		points = 0
	}

	rank := models.Rangking{
		Point: points,
	}

	return rank, nil
}
