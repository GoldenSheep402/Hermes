package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/redis/go-redis/v9"
)

type AnnounceResult struct {
	UploadDelta    int64
	DownloadDelta  int64
	SeedTimeDelta  int64
	StartedAt      time.Time
	StoredUpload   int64
	StoredDownload int64
	StoredLeft     int64
	FinishedAt     time.Time
}

var applyAnnounceScript = redis.NewScript(`
local event = ARGV[1]
local nowUnix = tonumber(ARGV[2]) or 0
local peerTTL = tonumber(ARGV[3]) or 2400
local uploaded = tonumber(ARGV[4]) or 0
local downloaded = tonumber(ARGV[5]) or 0
local left = tonumber(ARGV[6]) or 0
local peerRecordID = ARGV[7]
local torrentID = ARGV[8]
local userID = ARGV[9]
local peerID = ARGV[10]
local ip = ARGV[11]
local lanIP = ARGV[12]
local port = tonumber(ARGV[13]) or 0
local agent = ARGV[14]
local realtimeField = ARGV[15]
local expiredField = ARGV[16]
local realtimeTTL = tonumber(ARGV[17]) or 600
local pairMember = ARGV[18]
local isActive = tonumber(ARGV[19]) or 0

local hasPrev = redis.call('EXISTS', KEYS[1]) == 1
local prevUploaded = tonumber(redis.call('HGET', KEYS[1], 'uploaded') or '0')
local prevDownloaded = tonumber(redis.call('HGET', KEYS[1], 'downloaded') or '0')
local prevStartedAt = tonumber(redis.call('HGET', KEYS[1], 'started_at') or '0')
local prevLeftRaw = redis.call('HGET', KEYS[1], 'left')
local prevLeft = nil
if prevLeftRaw ~= false then
	prevLeft = tonumber(prevLeftRaw) or 0
end

local prevPairActive = tonumber(redis.call('HGET', KEYS[5], 'is_active') or '0')
local prevPairFinished = tonumber(redis.call('HGET', KEYS[5], 'is_finished') or '0')
local prevPairLastAction = tonumber(redis.call('HGET', KEYS[5], 'last_action') or '0')
local prevPairSeedTime = tonumber(redis.call('HGET', KEYS[5], 'seed_time') or '0')
local prevFinishedAt = tonumber(redis.call('HGET', KEYS[5], 'finished_at') or '0')

local uploadDelta = 0
local downloadDelta = 0
if hasPrev then
	if uploaded >= prevUploaded then
		uploadDelta = uploaded - prevUploaded
	elseif event == 'started' then
		uploadDelta = uploaded
	end

	if downloaded >= prevDownloaded then
		downloadDelta = downloaded - prevDownloaded
	elseif event == 'started' then
		downloadDelta = downloaded
	end
elseif event == 'started' or event == '' then
	uploadDelta = uploaded
	downloadDelta = downloaded
end

local storedUploaded = uploaded
local storedDownloaded = downloaded
local storedLeft = left
local startedAt = nowUnix

if hasPrev and event ~= 'started' then
	startedAt = prevStartedAt
	if startedAt <= 0 then
		startedAt = nowUnix
	end
	if uploaded < prevUploaded then
		storedUploaded = prevUploaded
	end
	if downloaded < prevDownloaded then
		storedDownloaded = prevDownloaded
	end
	if prevLeft ~= nil and storedLeft > prevLeft then
		storedLeft = prevLeft
	end
else
	startedAt = nowUnix
end

if storedLeft < 0 then
	storedLeft = 0
end
local storedSeeder = 0
if storedLeft == 0 then
	storedSeeder = 1
end

local pairFresh = 0
if prevPairActive == 1 and prevPairLastAction > 0 then
	if hasPrev or (nowUnix - prevPairLastAction) <= peerTTL then
		pairFresh = 1
	end
end

local seedTimeDelta = 0
if pairFresh == 1 and prevPairFinished == 1 and nowUnix > prevPairLastAction then
	seedTimeDelta = nowUnix - prevPairLastAction
end
local pairSeedTime = prevPairSeedTime + seedTimeDelta

local nextPairActive = isActive
local currentSeeder = storedSeeder
local nextPairFinished = prevPairFinished
if currentSeeder == 1 then
	nextPairFinished = 1
end
if event == 'stopped' then
	nextPairActive = 0
end

local finishedAt = prevFinishedAt
if currentSeeder == 1 and finishedAt <= 0 then
	finishedAt = nowUnix
end

local snatchDelta = 0
if prevPairFinished == 0 and currentSeeder == 1 then
	if event == 'completed' or hasPrev or pairFresh == 1 then
		snatchDelta = 1
	end
end

if event == 'stopped' then
	redis.call('DEL', KEYS[1])
	redis.call('ZREM', KEYS[2], peerID)
	redis.call('ZREM', KEYS[13], peerID)
	redis.call('ZREM', KEYS[14], peerID)
else
	redis.call('HSET', KEYS[1],
		'id', peerRecordID,
		'torrent_id', torrentID,
		'user_id', userID,
		'peer_id', peerID,
		'ip', ip,
		'lan_ip', lanIP,
		'port', port,
		'uploaded', storedUploaded,
		'downloaded', storedDownloaded,
		'left', storedLeft,
		'agent', agent,
		'is_seeder', storedSeeder,
		'started_at', startedAt,
		'last_action', nowUnix
	)
	redis.call('EXPIRE', KEYS[1], peerTTL)
	redis.call('ZADD', KEYS[2], nowUnix, peerID)
	redis.call('EXPIRE', KEYS[2], peerTTL)

	if currentSeeder == 1 then
		redis.call('ZADD', KEYS[13], nowUnix, peerID)
		redis.call('ZREM', KEYS[14], peerID)
	else
		redis.call('ZADD', KEYS[14], nowUnix, peerID)
		redis.call('ZREM', KEYS[13], peerID)
	end
	redis.call('EXPIRE', KEYS[13], peerTTL)
	redis.call('EXPIRE', KEYS[14], peerTTL)
end

if nextPairActive == 1 then
	if currentSeeder == 1 then
		redis.call('ZADD', KEYS[15], nowUnix, torrentID)
		redis.call('ZREM', KEYS[16], torrentID)
	else
		redis.call('ZADD', KEYS[16], nowUnix, torrentID)
		redis.call('ZREM', KEYS[15], torrentID)
	end
	redis.call('EXPIRE', KEYS[15], peerTTL)
	redis.call('EXPIRE', KEYS[16], peerTTL)
else
	redis.call('ZREM', KEYS[15], torrentID)
	redis.call('ZREM', KEYS[16], torrentID)
end

local userDirty = 0
local torrentDirty = 0

-- Resolve promotion multipliers (Lua-only; no DB round-trip).
local upFactor = 1
local downFactor = 1
local promoUp = tonumber(redis.call('HGET', KEYS[18], 'up') or '1') or 1
local promoDown = tonumber(redis.call('HGET', KEYS[18], 'down') or '1') or 1
local promoUpUntil = tonumber(redis.call('HGET', KEYS[18], 'up_until') or '0') or 0
local promoDownUntil = tonumber(redis.call('HGET', KEYS[18], 'down_until') or '0') or 0
if promoUpUntil > 0 and nowUnix > promoUpUntil then
	promoUp = 1
end
if promoDownUntil > 0 and nowUnix > promoDownUntil then
	promoDown = 1
end
if promoUp > 0 then
	upFactor = promoUp
end
if promoDown >= 0 then
	downFactor = promoDown
end
if redis.call('EXISTS', KEYS[19]) == 1 then
	downFactor = 0
end

-- First announce after upgrade: seed credited_* from current real_* so history is not reset.
if redis.call('HEXISTS', KEYS[3], 'credited_upload') == 0 then
	local ru = tonumber(redis.call('HGET', KEYS[3], 'real_upload') or '0') or 0
	local rd = tonumber(redis.call('HGET', KEYS[3], 'real_download') or '0') or 0
	redis.call('HSET', KEYS[3], 'credited_upload', ru, 'credited_download', rd)
end

local creditedUploadDelta = uploadDelta * upFactor
local creditedDownloadDelta = downloadDelta * downFactor

if uploadDelta > 0 or downloadDelta > 0 then
	if uploadDelta > 0 then
		redis.call('HINCRBY', KEYS[3], 'real_upload', uploadDelta)
		redis.call('HINCRBY', KEYS[17], 'real_upload', uploadDelta)
		redis.call('HINCRBY', KEYS[4], 'total_upload', uploadDelta)
		redis.call('HINCRBY', KEYS[5], 'uploaded', uploadDelta)
		redis.call('HINCRBY', KEYS[9], realtimeField, uploadDelta)
		redis.call('HINCRBY', KEYS[11], realtimeField, uploadDelta)
		redis.call('HDEL', KEYS[9], expiredField)
		redis.call('HDEL', KEYS[11], expiredField)
		redis.call('EXPIRE', KEYS[9], realtimeTTL)
		redis.call('EXPIRE', KEYS[11], realtimeTTL)
	end
	if downloadDelta > 0 then
		redis.call('HINCRBY', KEYS[3], 'real_download', downloadDelta)
		redis.call('HINCRBY', KEYS[17], 'real_download', downloadDelta)
		redis.call('HINCRBY', KEYS[4], 'total_download', downloadDelta)
		redis.call('HINCRBY', KEYS[5], 'downloaded', downloadDelta)
		redis.call('HINCRBY', KEYS[10], realtimeField, downloadDelta)
		redis.call('HINCRBY', KEYS[12], realtimeField, downloadDelta)
		redis.call('HDEL', KEYS[10], expiredField)
		redis.call('HDEL', KEYS[12], expiredField)
		redis.call('EXPIRE', KEYS[10], realtimeTTL)
		redis.call('EXPIRE', KEYS[12], realtimeTTL)
	end
	if creditedUploadDelta > 0 then
		redis.call('HINCRBY', KEYS[3], 'credited_upload', creditedUploadDelta)
	end
	if creditedDownloadDelta > 0 then
		redis.call('HINCRBY', KEYS[3], 'credited_download', creditedDownloadDelta)
	end
	userDirty = 1
	torrentDirty = 1
end

if seedTimeDelta > 0 then
	redis.call('HINCRBY', KEYS[3], 'seed_time', seedTimeDelta)
	userDirty = 1
end

if snatchDelta > 0 then
	redis.call('HINCRBY', KEYS[4], 'snatch_count', snatchDelta)
	torrentDirty = 1
end

if userDirty == 1 then
	redis.call('HINCRBY', KEYS[3], '_version', 1)
	redis.call('SADD', KEYS[6], userID)
end
if torrentDirty == 1 or event == 'started' or event == 'stopped' or prevLeft ~= storedLeft then
	redis.call('HINCRBY', KEYS[4], '_version', 1)
	redis.call('SADD', KEYS[7], torrentID)
end

redis.call('HSET', KEYS[5],
	'is_active', nextPairActive,
	'is_finished', nextPairFinished,
	'last_action', nowUnix,
	'seed_time', pairSeedTime,
	'finished_at', finishedAt
)
redis.call('HINCRBY', KEYS[5], '_version', 1)
redis.call('SADD', KEYS[8], pairMember)

-- Track finished+active pairs for hourly seeding bonus settlement.
if nextPairActive == 1 and nextPairFinished == 1 then
	redis.call('SADD', KEYS[20], pairMember)
else
	redis.call('SREM', KEYS[20], pairMember)
end

return {uploadDelta, downloadDelta, seedTimeDelta, startedAt, storedUploaded, storedDownloaded, storedLeft, finishedAt}
`)

func (t *traffic) ApplyAnnounce(ctx context.Context, peerData *trackerModel.Peer, event string) (*AnnounceResult, error) {
	if t.rds == nil || t.db == nil {
		return nil, fmt.Errorf("tracker storage unavailable")
	}
	if peerData == nil {
		return nil, fmt.Errorf("missing peer data")
	}
	if peerData.UserID == "" || peerData.TorrentID == "" || peerData.PeerID == "" {
		return nil, fmt.Errorf("missing peer identity")
	}

	if err := t.ensureUserTotalsInitialized(ctx, peerData.UserID); err != nil {
		return nil, err
	}
	if err := t.ensureSiteTotalsInitialized(ctx); err != nil {
		return nil, err
	}
	if err := t.ensureTorrentTotalsInitialized(ctx, peerData.TorrentID); err != nil {
		return nil, err
	}
	if err := t.ensurePairTotalsInitialized(ctx, peerData.UserID, peerData.TorrentID); err != nil {
		return nil, err
	}

	now := peerData.LastAction
	if now.IsZero() {
		now = time.Now()
		peerData.LastAction = now
	}

	realtimeBucket := alignRealtimeBucket(now.Unix())
	expiredBucket := realtimeBucket - realtimeWindowSeconds - realtimeBucketSpanSeconds
	result, err := applyAnnounceScript.Run(
		ctx,
		t.rds,
		[]string{
			peerKey(peerData.TorrentID, peerData.PeerID),
			torrentPeersKey(peerData.TorrentID),
			userTrafficKey(peerData.UserID),
			torrentTrafficKey(peerData.TorrentID),
			pairTrafficKey(peerData.UserID, peerData.TorrentID),
			trafficDirtyUsersKey(),
			trafficDirtyTorrentsKey(),
			trafficDirtyPairsKey(),
			userRealtimeUploadKey(peerData.UserID),
			userRealtimeDownloadKey(peerData.UserID),
			siteRealtimeUploadKey(),
			siteRealtimeDownloadKey(),
			torrentSeederKey(peerData.TorrentID),
			torrentLeecherKey(peerData.TorrentID),
			userSeedingKey(peerData.UserID),
			userDownloadingKey(peerData.UserID),
			siteTrafficTotalsKey(),
			torrentPromoKey(peerData.TorrentID),
			globalFreeKey(),
			seedingPairsKey(),
		},
		event,
		now.Unix(),
		int64(peerTTL/time.Second),
		nonNegative(peerData.Uploaded),
		nonNegative(peerData.Downloaded),
		nonNegative(peerData.Left),
		peerData.ID,
		peerData.TorrentID,
		peerData.UserID,
		peerData.PeerID,
		peerData.IP,
		peerData.LanIP,
		peerData.Port,
		peerData.Agent,
		strconv.FormatInt(realtimeBucket, 10),
		strconv.FormatInt(expiredBucket, 10),
		int64(realtimeBucketTTL/time.Second),
		pairTrafficMember(peerData.UserID, peerData.TorrentID),
		boolToInt(event != "stopped"),
	).Result()
	if err != nil {
		return nil, err
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 8 {
		return nil, fmt.Errorf("unexpected announce result: %T", result)
	}

	res := &AnnounceResult{
		UploadDelta:    parseAnyInt64(values[0]),
		DownloadDelta:  parseAnyInt64(values[1]),
		SeedTimeDelta:  parseAnyInt64(values[2]),
		StoredUpload:   parseAnyInt64(values[4]),
		StoredDownload: parseAnyInt64(values[5]),
		StoredLeft:     parseAnyInt64(values[6]),
	}
	if startedAtUnix := parseAnyInt64(values[3]); startedAtUnix > 0 {
		res.StartedAt = time.Unix(startedAtUnix, 0)
	}
	if finishedAtUnix := parseAnyInt64(values[7]); finishedAtUnix > 0 {
		res.FinishedAt = time.Unix(finishedAtUnix, 0)
	}
	return res, nil
}
