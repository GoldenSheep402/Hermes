package service

import (
	"context"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"github.com/zeebo/bencode"
	"go.uber.org/zap"
	"net/url"
	"time"
)

var _ trackerV1.TrackerServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trackerV1.UnimplementedTrackerServiceServer
}

type trackerResp struct {
	FailureReason  string      `bencode:"failure reason,omitempty"`  // 可选
	WarningMessage string      `bencode:"warning message,omitempty"` // 可选
	Interval       int         `bencode:"interval"`                  // 必需
	MinInterval    int         `bencode:"min interval,omitempty"`    // 可选
	TrackerID      string      `bencode:"tracker id,omitempty"`      // 可选
	Complete       int         `bencode:"complete"`                  // 必需
	Incomplete     int         `bencode:"incomplete"`                // 必需
	Peers          interface{} `bencode:"peers"`                     // 必需
}

func (s *S) GetTracker(ctx context.Context, req *trackerV1.GetTrackerRequest) (*trackerV1.GetTrackerResponse, error) {
	key := req.Key
	fmt.Printf("GetTracker: %s\n", key)

	infoHash := req.InfoHash
	peerId := req.PeerId
	port := req.Port
	// uploaded := req.Uploaded
	// downloaded := req.Downloaded
	// left := req.Left
	event := req.Event
	ip := req.Ip
	// numWant := req.NumWant
	// compact := req.Compact
	// noPeerId := req.NoPeerId
	// corrupt := req.Corrupt
	// supportCrypto := req.SupportCrypto
	// redundant := req.Redundant

	decodedHash, err := url.QueryUnescape(infoHash)
	if err != nil {
		return nil, err
	}
	infoHashBytes := []byte(decodedHash)
	hexString := fmt.Sprintf("%x", infoHashBytes)
	infoHashDecoded := hexString

	switch event {
	case "started", "":
		if ip == "" || port == 0 || peerId == "" {
			return nil, fmt.Errorf("invalid request")
		}

		peer := &model.Peer{
			PeerID:   peerId,
			IP:       ip,
			Port:     int(port),
			LastSeen: time.Now(),
			Status:   1,
		}

		err := dao.Peer.AddPeer(ctx, infoHashDecoded, peer)
		if err != nil {
			respDetail := &trackerResp{
				FailureReason: "Failed to add peer",
			}

			_, err := bencode.EncodeBytes(respDetail)
			if err != nil {
				return nil, err
			}

			return &trackerV1.GetTrackerResponse{
				Response: "test",
			}, nil
		}
	case "stopped":
		err := dao.Peer.RemovePeer(ctx, infoHashDecoded, peerId)
		if err != nil {
			return &trackerV1.GetTrackerResponse{
				Response: "test",
			}, nil
		}
		return &trackerV1.GetTrackerResponse{}, nil
	case "completed":
		peer := &model.Peer{
			PeerID:   peerId,
			IP:       ip,
			Port:     int(port),
			LastSeen: time.Now(),
			Status:   0,
		}

		err := dao.Peer.AddPeer(ctx, infoHashDecoded, peer)
		if err != nil {
			return nil, err
		}

		return &trackerV1.GetTrackerResponse{}, nil
	}

	// peers, err := dao.Peer.GetPeers(ctx, infoHashDecoded, int(numWant))
	// var completed, incompleted int32
	// for _, peer := range peers {
	// 	if peer.Status == 0 {
	// 		completed++
	// 	}
	// }
	//
	// incompleted = int32(len(peers)) - completed
	//
	// peersToReturn := make([]*trackerV1.PeerInfo, 0, len(peers))
	// for _, peer := range peers {
	// 	var peersData []byte
	// 	ip := net.ParseIP(peer.IP).To4()
	// 	if ip == nil {
	// 		continue
	// 	}
	// 	peersToReturn = append(peersToReturn, &trackerV1.PeerInfo{
	// 		PeerId: []byte(peer.PeerID),
	// 		Ip:     ip,
	// 		Port:   int32(peer.Port),
	// 	})
	// }

	// if compact == 1 {
	// 	var peersData []byte
	// 	for _, peer := range peers {
	// 		ip := net.ParseIP(peer.IP).To4()
	// 		if ip == nil {
	// 			continue
	// 		}
	// 		portBytes := make([]byte, 2)
	// 		binary.BigEndian.PutUint16(portBytes, uint16(peer.Port))
	// 		peersData = append(peersData, ip...)
	// 		peersData = append(peersData, portBytes...)
	// 	}
	// 	resp.Peers = peersData
	// } else {
	// 	var peersList []Peer
	// 	for _, peer := range peers {
	// 		if req.NoPeerID == 1 {
	// 			peer.PeerID = ""
	// 		}
	// 		peersList = append(peersList, *peer)
	// 	}
	// 	resp.Peers = peersList
	// }
	//
	// respDeail := &trackerV1.GetTrackerResponseDetail{
	// 	FailureReason:  "",
	// 	WarningMessage: "",
	// 	Interval:       1800,
	// 	MinInterval:    0,
	// 	TrackerId:      "beta",
	// 	Complete:       completed,
	// 	Incomplete:     incompleted,
	// 	Peers:          peersToReturn,
	// 	PeersCompact:   nil,
	// }
	return &trackerV1.GetTrackerResponse{
		Response: "test",
	}, nil
}
