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
	SELECT sum(point)
	FROM rangkings
	WHERE user_id = $1 and next = true AND deleted_at IS NULL;
	`

	queryGetUserAndPoint = `
	WITH UserPoints AS (
		SELECT user_id, SUM(point) AS total_points 
		FROM rangkings and next = 'true' 
		WHERE deleted_at IS NULL
		GROUP BY user_id
	)
	SELECT u.user_id, u.username, u.profile, up.total_points
	FROM UserPoints up
	JOIN users u ON u.user_id  = up.user_id;
	`

	queryUpdateNext = `
		UPDATE rangkings SET next = true, updated_at = $1 where user_id = $2 and soal_id= $3 and deleted_at is null
	`
)

type RangkingRepository interface {
	CreateRangking(ctx context.Context, rank models.Rangking) (id int64, err error)
	GetPointByUserId(ctx context.Context, id int64) (rank models.Rangking, err error)
	GetUserAndPoint(ctx context.Context) (rank []models.RangkingUser, err error)
	UpdateNextUser(ctx context.Context, rank models.Rangking) (id int64, err error)
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

func (r *RangkingRepositoryImpl) GetUserAndPoint(ctx context.Context) (rank []models.RangkingUser, err error) {
	rows, err := r.db.QueryContext(ctx, queryGetUserAndPoint)
	if err != nil {
		if err != sql.ErrNoRows {
			err = fmt.Errorf("queryGetSoal err: %+v", err)
			return
		}
		err = nil
		return
	}
	defer rows.Close()

	var s models.RangkingUser

	for rows.Next() {
		err = rows.Scan(&s.UserId, &s.Username, &s.Profile, &s.Point)
		if err != nil {
			err = fmt.Errorf("row scan err: %+v", err)
			return nil, err
		}
		rank = append(rank, s)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return rank, nil
}

func (r *RangkingRepositoryImpl) UpdateNextUser(ctx context.Context, rank models.Rangking) (id int64, err error) {
	updatedAt := time.Now()

	_, err = r.db.ExecContext(ctx, queryUpdateNext, updatedAt, rank.UserId, rank.SoalId)
	if err != nil {
		return 0, err
	}

	return id, nil
}
