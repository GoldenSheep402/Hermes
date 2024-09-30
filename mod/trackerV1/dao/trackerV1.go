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

// CheckKey Check key and set into redis
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
			_, err = t.SetKeyWithTTL(ctx, key, _user.ID+"/"+trackerV1Values.OK, trackerV1Values.TTL)
			if err != nil {
				return "", _user.ID, nil
			}
			return trackerV1Values.OK, _user.ID, nil
		default:
			return "", _user.ID, nil
		}
	}

	parts := strings.Split(val, "/")
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
// Upload:Download
func (t *trackerV1) HandelDownloadAndUpload(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	switch status {
	case trackerV1Values.Downloading:
		t.HandleDownloading(ctx, torrentID, userID, status, uploadMB, downloadMB)
	case trackerV1Values.ReadySeeding:
		t.HandleReadySeeding(ctx, torrentID, userID, status, uploadMB, downloadMB)
	case trackerV1Values.Seeding:
		t.HandleSeeding(ctx, torrentID, userID, status, uploadMB, downloadMB)
	case trackerV1Values.Stopped:
		t.HandleStop(ctx, torrentID, userID, status, uploadMB, downloadMB)
	}

	return nil
}

// HandleDownloading Handle downloading.
// If it is downloading
func (t *trackerV1) HandleDownloading(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	key := TorrentSumKey(torrentID, userID)

	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check hash existence in redis: %v", err)
	}

	// if exist and the value is greater than old, update
	if exists != 0 {
		oldData, err := t.rds.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("failed to get old user data in redis: %v", err)
		}

		parts := strings.Split(oldData, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid old data format for key %s/%s", torrentID, userID)
		}

		oldUploadMB, _ := strconv.ParseInt(parts[0], 10, 64)
		oldDownloadMB, _ := strconv.ParseInt(parts[1], 10, 64)

		// is smaller it is reset
		uploadSmaller := uploadMB < oldUploadMB
		downloadSmaller := downloadMB < oldDownloadMB
		if uploadSmaller || downloadSmaller {
			if uploadSmaller {
				err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, 0)
				if err != nil {
					return fmt.Errorf("failed to increase single sum: %v", err)
				}
			}

			if downloadSmaller {
				err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, 0, downloadMB)
				if err != nil {
					return fmt.Errorf("failed to increase single sum: %v", err)
				}
			}

			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		} else {
			// the new data is greater than old data, refresh data in redis
			// THERE MIGHT BE SOME BUG HERE
			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
			uploadToIncrease := uploadMB - oldUploadMB
			downloadToIncrease := downloadMB - oldDownloadMB
			if uploadToIncrease < 0 {
				uploadToIncrease = 0
			}
			if downloadToIncrease < 0 {
				downloadToIncrease = 0
			}
			err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadToIncrease, downloadToIncrease)
			if err != nil {
				return fmt.Errorf("failed to increase single sum: %v", err)
			}
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		}
	} else {
		// no exist
		data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)

		err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
		if err != nil {
			return fmt.Errorf("failed to set user data in redis: %v", err)
		}

		if err := SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, downloadMB); err != nil {
			return fmt.Errorf("failed to update single sum: %v", err)
		}
	}

	if err := t.HandleClientStatus(ctx, torrentID, userID, status); err != nil {
		return fmt.Errorf("failed to update torrent status: %v", err)
	}

	return nil
}

func (t *trackerV1) HandleClientStatus(ctx context.Context, torrentID, userID string, status int) error {
	seedingKey := TorrentSeedingKey(torrentID, userID)
	// Handle seeding count
	switch status {
	case trackerV1Values.Downloading:
		t.rds.Set(ctx, seedingKey, trackerV1Values.Downloading, 24*time.Hour)
	case trackerV1Values.Seeding:
		t.rds.Set(ctx, seedingKey, trackerV1Values.Seeding, 24*time.Hour)
	case trackerV1Values.Finished:
		t.rds.Set(ctx, seedingKey, trackerV1Values.Finished, 24*time.Hour)
	case trackerV1Values.Stopped:
		t.rds.Set(ctx, seedingKey, trackerV1Values.Stopped, 1*time.Hour)

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

// HandleReadySeeding Handle ready seeding.
// If the key is exist, set the value to 0:0 and set the old value to db
func (t *trackerV1) HandleReadySeeding(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	key := TorrentSumKey(torrentID, userID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check hash existence in redis: %v", err)
	}

	// exist
	if exists != 0 {
		oldData, err := t.rds.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("failed to get old user data in redis: %v", err)
		}

		parts := strings.Split(oldData, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid old data format for key %s/%s", torrentID, userID)
		}

		oldUploadMB, _ := strconv.ParseInt(parts[0], 10, 64)
		oldDownloadMB, _ := strconv.ParseInt(parts[1], 10, 64)

		// is smaller it is reset
		uploadSmaller := uploadMB < oldUploadMB
		downloadSmaller := downloadMB < oldDownloadMB
		if uploadSmaller || downloadSmaller {
			//if uploadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, 0)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}
			//
			//if downloadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, 0, downloadMB)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}

			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		} else {
			// the new data is greater than old data, refresh data in redis
			// THERE MIGHT BE SOME BUG HERE
			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
			uploadToIncrease := uploadMB - oldUploadMB
			downloadToIncrease := downloadMB - oldDownloadMB
			if uploadToIncrease < 0 {
				uploadToIncrease = 0
			}
			if downloadToIncrease < 0 {
				downloadToIncrease = 0
			}
			err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadToIncrease, downloadToIncrease)
			if err != nil {
				return fmt.Errorf("failed to increase single sum: %v", err)
			}
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		}
	} else {
		// not exist
		err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, downloadMB)
		if err != nil {
			return fmt.Errorf("failed to increase single sum: %v", err)
		}

		data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
		err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
		if err != nil {
			return fmt.Errorf("failed to set user data in redis: %v", err)
		}
	}

	if err := t.HandleClientStatus(ctx, torrentID, userID, status); err != nil {
		return fmt.Errorf("failed to update torrent status: %v", err)
	}

	return nil
}

// HandleSeeding Handle seeding.
func (t *trackerV1) HandleSeeding(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	key := TorrentSumKey(torrentID, userID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check hash existence in redis: %v", err)
	}

	// if exist and the value is greater than old, update
	// if not exist, set the value to db and reset the value
	if exists != 0 {
		oldData, err := t.rds.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("failed to get old user data in redis: %v", err)
		}

		parts := strings.Split(oldData, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid old data format for key %s/%s", torrentID, userID)
		}

		oldUploadMB, _ := strconv.ParseInt(parts[0], 10, 64)
		oldDownloadMB, _ := strconv.ParseInt(parts[1], 10, 64)

		// is smaller it is reset
		uploadSmaller := uploadMB < oldUploadMB
		downloadSmaller := downloadMB < oldDownloadMB
		if uploadSmaller || downloadSmaller {
			//if uploadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, 0)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}
			//
			//if downloadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, 0, downloadMB)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}

			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		} else {
			// the new data is greater than old data, refresh data in redis
			// THERE MIGHT BE SOME BUG HERE
			data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)

			parts := strings.Split(oldData, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid old data format for key %s/%s", torrentID, userID)
			}

			_oldUploadMB, _ := strconv.ParseInt(parts[0], 10, 64)
			_oldDownloadMB, _ := strconv.ParseInt(parts[1], 10, 64)

			uploadToIncrease := uploadMB - _oldUploadMB
			downloadToIncrease := downloadMB - _oldDownloadMB
			if uploadToIncrease < 0 {
				uploadToIncrease = 0
			}
			if downloadToIncrease < 0 {
				downloadToIncrease = 0
			}
			err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadToIncrease, downloadToIncrease)
			if err != nil {
				return fmt.Errorf("failed to increase single sum: %v", err)
			}
			err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
			if err != nil {
				return fmt.Errorf("failed to set user data in redis: %v", err)
			}
		}
	} else {
		// not exist
		err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, downloadMB)
		if err != nil {
			return fmt.Errorf("failed to increase single sum: %v", err)
		}

		data := fmt.Sprintf("%v:%v", uploadMB, downloadMB)
		err = t.rds.Set(ctx, key, data, 24*time.Hour).Err()
		if err != nil {
			return fmt.Errorf("failed to set user data in redis: %v", err)
		}
	}

	if err := t.HandleClientStatus(ctx, torrentID, userID, status); err != nil {
		return fmt.Errorf("failed to update torrent status: %v", err)
	}
	return nil
}

// HandleStop Handle stop.
func (t *trackerV1) HandleStop(ctx context.Context, torrentID, userID string, status int, uploadMB, downloadMB int64) error {
	key := TorrentSumKey(torrentID, userID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check hash existence in redis: %v", err)
	}

	// if exist and the value is greater than old, update
	// if not exist, set the value to db and reset the value
	if exists != 0 {
		oldData, err := t.rds.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("failed to get old user data in redis: %v", err)
		}

		parts := strings.Split(oldData, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid old data format for key %s/%s", torrentID, userID)
		}

		oldUploadMB, _ := strconv.ParseInt(parts[0], 10, 64)
		oldDownloadMB, _ := strconv.ParseInt(parts[1], 10, 64)

		// is smaller it is reset
		uploadSmaller := uploadMB < oldUploadMB
		downloadSmaller := downloadMB < oldDownloadMB
		if uploadSmaller || downloadSmaller {
			//if uploadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, 0)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}
			//
			//if downloadSmaller {
			//	err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, 0, downloadMB)
			//	if err != nil {
			//		return fmt.Errorf("failed to increase single sum: %v", err)
			//	}
			//}

			err = t.rds.Del(ctx, key).Err()
			if err != nil {
				return fmt.Errorf("failed to delete key: %v", err)
			}
		} else {
			// THERE MIGHT BE SOME BUG HERE
			uploadToIncrease := uploadMB - oldUploadMB
			downloadToIncrease := downloadMB - oldDownloadMB
			if uploadToIncrease < 0 {
				uploadToIncrease = 0
			}
			if downloadToIncrease < 0 {
				downloadToIncrease = 0
			}
			err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadToIncrease, downloadToIncrease)
			if err != nil {
				return fmt.Errorf("failed to increase single sum: %v", err)
			}
			err = t.rds.Del(ctx, key).Err()
			if err != nil {
				return fmt.Errorf("failed to delete key: %v", err)
			}
		}
	} else {
		// not exist
		err = SingleSum.IncreaseSingleSum(ctx, torrentID, userID, uploadMB, downloadMB)
		if err != nil {
			return fmt.Errorf("failed to increase single sum: %v", err)
		}

		err = t.rds.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to delete key: %v", err)
		}
	}

	return nil
}

func TorrentDownloadingKey(torrentID string, uid string) string {
	return "TorrentDownloading" + "/" + torrentID + "/" + uid
}

func TorrentSumKey(torrentID, uid string) string {
	return "TorrentSum" + "/" + torrentID + "/" + uid
}

func TorrentSeedingKey(torrentID, uid string) string {
	return "TorrentSeeding" + "/" + torrentID + "/" + uid
}

func TorrentStopKey(torrentID, uid string) string {
	return "TorrentStop" + "/" + torrentID + "/" + uid
}
