package dao

import (
	"context"
	"errors"
	"fmt"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

type trackerV1 struct {
	rds *redis.Client
}

func (t *trackerV1) Init(rds *redis.Client) error {
	t.rds = rds
	return nil
}

func (t *trackerV1) CheckKey(ctx context.Context, key string) (string, string, error) {
	val, err := t.rds.Get(ctx, key).Result()
	var _user userModel.User

	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			_user, err := userDao.User.CheckKey(ctx, key)
			if err != nil {
				return "", "", err
			}

			// TODO: status check
			_, err = t.SetKeyWithTTL(ctx, key, _user.ID+":"+trackerV1Values.OK, trackerV1Values.TTL)
			if err != nil {
				return "", _user.ID, nil
			}
			return trackerV1Values.OK, _user.ID, nil
		default:
			return "", _user.ID, nil
		}
	}

	//if errors.Is(err, redis.Nil) {
	//	_user, err := userDao.User.CheckKey(ctx, key)
	//	if err != nil {
	//		return "", "", err
	//	}
	//
	//	// TODO: status check
	//	_, err = t.SetKeyWithTTL(ctx, key, _user.ID+":"+trackerV1Values.OK, trackerV1Values.TTL)
	//	if err != nil {
	//		return "", _user.ID, nil
	//	}
	//	return trackerV1Values.OK, _user.ID, nil
	//}

	//if err != nil {
	//return "", _user.ID, nil
	//}

	parts := strings.Split(val, ":")
	if len(parts) != 2 {
		return "", _user.ID, fmt.Errorf("unknown key status: %s", val)
	}

	status := parts[1]
	_user.ID = parts[0]

	if status == trackerV1Values.OK {
		if err := t.rds.Expire(ctx, key, trackerV1Values.TTL).Err(); err != nil {
			return "", _user.ID, err
		}
		return trackerV1Values.OK, _user.ID, nil
	} else if val == trackerV1Values.Banned {
		if err := t.rds.Expire(ctx, _user.ID+":"+trackerV1Values.Banned, trackerV1Values.TTL).Err(); err != nil {
			return "", _user.ID, err
		}
		return trackerV1Values.Banned, _user.ID, nil
	}

	return "", _user.ID, fmt.Errorf("unknown key status: %s", val)
}

func (t *trackerV1) SetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (string, error) {
	if err := t.rds.Set(ctx, key, value, ttl).Err(); err != nil {
		return "", err
	}
	return value, nil
}

// HandelDownloadAndUpload Handel download and upload.
func (t *trackerV1) HandelDownloadAndUpload(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	// key: torrent:ID
	key := "TorrentSum:" + torrentID

	exists, err := t.rds.HExists(ctx, key, userID).Result()
	if err != nil {
		return fmt.Errorf("failed to check hash existence in redis: %v", err)
	}

	//uploadBytes, err := strconv.ParseInt(uploadStr, 10, 64)
	//if err != nil {
	//	return fmt.Errorf("failed to parse uploaded data: %v", err)
	//}
	//downloadBytes, err := strconv.ParseInt(downloadStr, 10, 64)
	//if err != nil {
	//	return fmt.Errorf("failed to parse downloaded data: %v", err)
	//}

	//uploadMB := float64(uploadBytes) / (1024 * 1024)
	//downloadMB := float64(downloadBytes) / (1024 * 1024)

	if exists {
		oldData, err := t.rds.HGet(ctx, key, userID).Result()
		if err != nil {
			return fmt.Errorf("failed to get old user data in redis: %v", err)
		}

		parts := strings.Split(oldData, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid old data format for userID %s", userID)
		}

		oldUploadMB, _ := strconv.ParseInt(parts[0], 64, 10)
		oldDownloadMB, _ := strconv.ParseInt(parts[1], 64, 10)

		if uploadMB > oldUploadMB || downloadMB > oldDownloadMB {
			if err := SingleSum.UpdateSingleSum(ctx, torrentID, userID, uploadMB, downloadMB); err != nil {
				return fmt.Errorf("failed to update single sum: %v", err)
			}

			if err := TorrentStatus.IncrementUploadAndDownload(ctx, torrentID, userID, status, uploadMB, downloadMB); err != nil {
				return fmt.Errorf("failed to update torrent status: %v", err)
			}
		}

		// no update
		return nil
	}

	data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)

	err = t.rds.HSet(ctx, key, userID, data).Err()
	if err != nil {
		return fmt.Errorf("failed to set user data in redis: %v", err)
	}

	err = t.rds.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration in redis: %v", err)
	}

	if err := SingleSum.UpdateSingleSum(ctx, torrentID, userID, uploadMB, downloadMB); err != nil {
		return fmt.Errorf("failed to update single sum: %v", err)
	}

	if err := TorrentStatus.IncrementUploadAndDownload(ctx, torrentID, userID, status, uploadMB, downloadMB); err != nil {
		return fmt.Errorf("failed to update torrent status: %v", err)
	}

	return nil
}

func (t *trackerV1) GetTorrentID(ctx context.Context, hash string) (string, error) {
	key := "Hash:" + hash

	id, err := t.rds.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			//	get torrent
			torrent, _err := torrentDao.Torrent.GetTorrentByHash(ctx, hash)
			if _err != nil {
				return "", _err
			}
			_, _err = t.SetKeyWithTTL(ctx, key, torrent.ID, 24*time.Hour)
			if _err != nil {
				return "", _err
			}
			return torrent.ID, nil
		default:
			return "", err
		}
	}

	err = t.rds.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return id, nil
}
