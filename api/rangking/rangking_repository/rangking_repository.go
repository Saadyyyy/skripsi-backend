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
	queryCreateRank = `
		insert into rangkings(user_id,category_id,point,created_at)values($1,$2,$3,$4)returning rank_id
	`

	queryUpdateRank = `
		update rangkings SET point=$1,updated_at =$2  where user_id= $3 and category_id=$4 and deleted_at is null
	`

	queryCheckRank = `select rank_id,user_id,category_id,point from rangkings where user_id =$1 and deleted_at is null`

	queryGetRank = `
	SELECT rangkings.user_id, SUM(rangkings.point) AS total_points ,u.username  as username,u.profile as profile
	FROM rangkings 
	JOIN  users u  ON rangkings.user_id = u.user_id
	WHERE rangkings.deleted_at IS NULL
	GROUP BY rangkings.user_id ,u.username ,u.profile
	ORDER BY total_points DESC
	LIMIT 10;
	`
)

type RangkingRepository interface {
	CreateRank(ctx context.Context, rank models.Rangking) (ID int64, err error)
	UpdatedRank(ctx context.Context, rank models.Rangking) (ID int64, err error)
	CheckRank(ctx context.Context, ID int64) (rank []models.Rangking, err error)
	GetRank(ctx context.Context, rank models.RangkingRes) (result []models.RangkingRes, err error)
}

type RangkingRepositoryImpl struct {
	db *sqlx.DB
}

func NewRangkingRepository(db *sqlx.DB) RangkingRepository {
	return &RangkingRepositoryImpl{db: db}
}

func (r *RangkingRepositoryImpl) CheckRank(ctx context.Context, ID int64) (result []models.Rangking, err error) {
	rows, err := r.db.QueryContext(ctx, queryCheckRank, ID)
	if err != nil {
		if err != sql.ErrNoRows {
			err = fmt.Errorf("queryCheckRank err: %+v", err)
			return
		}
		err = nil
	}
	defer rows.Close()
	for rows.Next() {
		rank := models.Rangking{}
		err = rows.Scan(&rank.RankId, &rank.UserId, &rank.CategoryId, &rank.Point)
		if err != nil {
			err = fmt.Errorf("row scan err: %+v", err)
			return nil, err
		}
		result = append(result, rank)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return result, nil
}

func (r *RangkingRepositoryImpl) GetRank(ctx context.Context, rank models.RangkingRes) (result []models.RangkingRes, err error) {
	rows, err := r.db.QueryContext(ctx, queryGetRank)
	if err != nil {
		if err != sql.ErrNoRows {
			err = fmt.Errorf("queryGetSoal err: %+v", err)
			return
		}
		err = nil
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&rank.UserId, &rank.Point, &rank.Username, &rank.Profile)
		if err != nil {
			err = fmt.Errorf("row scan err: %+v", err)
			return nil, err
		}
		result = append(result, rank)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return result, nil
}

func (r *RangkingRepositoryImpl) UpdatedRank(ctx context.Context, rank models.Rangking) (ID int64, err error) {
	updatedAt := time.Now()

	_, err = r.db.ExecContext(ctx, queryUpdateRank, rank.Point, updatedAt, rank.UserId, rank.CategoryId)
	if err != nil {
		fmt.Errorf("Query Updated Rank Error :", err)

	}
	return ID, nil
}

func (r *RangkingRepositoryImpl) CreateRank(ctx context.Context, rank models.Rangking) (ID int64, err error) {
	createdAt := time.Now()
	err = r.db.QueryRowContext(ctx, queryCreateRank, rank.UserId, rank.CategoryId, rank.Point, createdAt).Scan(&ID)
	if err != nil {
		err = fmt.Errorf("Query Rank Error :", err)
		return
	}
	fmt.Println("err", err)
	return ID, nil
}
