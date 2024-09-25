package handlers

import (
	"encoding/binary"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/juanjiTech/jin"
	"github.com/zeebo/bencode"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Registry(jinE *jin.Engine) {
	trackerGroup := jinE.Group("/tracker")

	trackerGroup.Group("/announce").
		GET("/key/:key", AnnounceWithKey)
}

func AnnounceWithKey(c *jin.Context) {
	ctx := c.Request.Context()

	// TODO: implement
	// key := c.Params("key")

	type AnnounceRequest struct {
		InfoHash      string // 必需
		PeerID        string // 必需
		Port          int    // 必需
		Uploaded      int    // 必需
		Downloaded    int    // 必需
		Left          int    // 必需
		Event         string // 可选
		IP            string // 可选
		NumWant       int    // 可选
		Key           string // 可选
		Compact       int    // 可选
		NoPeerID      int    // 可选
		TrackerID     string // 可选
		Corrupt       int    // 非标准参数
		SupportCrypto int    // 非标准参数
		Redundant     int    // 非标准参数
	}

	type AnnounceResponse struct {
		FailureReason  string      `bencode:"failure reason,omitempty"`  // 可选
		WarningMessage string      `bencode:"warning message,omitempty"` // 可选
		Interval       int         `bencode:"interval"`                  // 必需
		MinInterval    int         `bencode:"min interval,omitempty"`    // 可选
		TrackerID      string      `bencode:"tracker id,omitempty"`      // 可选
		Complete       int         `bencode:"complete"`                  // 必需
		Incomplete     int         `bencode:"incomplete"`                // 必需
		Peers          interface{} `bencode:"peers"`                     // 必需
	}

	type Peer struct {
		PeerID string `bencode:"peer id"` // 对等节点的唯一标识符
		IP     string `bencode:"ip"`      // 对等节点的 IP 地址
		Port   int    `bencode:"port"`    // 对等节点的端口号
	}
	req := &AnnounceRequest{}

	infoHash := c.Request.URL.Query().Get("info_hash")
	req.PeerID = c.Request.URL.Query().Get("peer_id")
	portStr := c.Request.URL.Query().Get("port")
	uploadedStr := c.Request.URL.Query().Get("uploaded")
	downloadedStr := c.Request.URL.Query().Get("downloaded")
	leftStr := c.Request.URL.Query().Get("left")
	req.Event = c.Request.URL.Query().Get("event")
	req.IP = c.Request.URL.Query().Get("ip")
	numWantStr := c.Request.URL.Query().Get("numwant")
	req.Key = c.Request.URL.Query().Get("key")
	compactStr := c.Request.URL.Query().Get("compact")
	noPeerIDStr := c.Request.URL.Query().Get("no_peer_id")
	req.TrackerID = c.Request.URL.Query().Get("trackerid")
	corruptStr := c.Request.URL.Query().Get("corrupt")
	supportCryptoStr := c.Request.URL.Query().Get("supportcrypto")
	redundantStr := c.Request.URL.Query().Get("redundant")

	var err error
	decodedHash, err := url.QueryUnescape(infoHash)
	if err != nil {
		c.Writer.WriteString("info_hash decode error")
		return
	}

	infoHashBytes := []byte(decodedHash)
	hexString := fmt.Sprintf("%x", infoHashBytes)

	req.InfoHash = hexString

	if req.Port, err = strconv.Atoi(portStr); err != nil {
		c.Writer.WriteString("Invalid port")
		return
	}
	if req.Uploaded, err = strconv.Atoi(uploadedStr); err != nil {
		c.Writer.WriteString("Invalid uploaded")
		return
	}
	if req.Downloaded, err = strconv.Atoi(downloadedStr); err != nil {
		c.Writer.WriteString("Invalid downloaded")
		return
	}
	if req.Left, err = strconv.Atoi(leftStr); err != nil {
		c.Writer.WriteString("Invalid left")
		return
	}

	if numWantStr != "" {
		req.NumWant, _ = strconv.Atoi(numWantStr)
	} else {
		req.NumWant = 50
	}
	if compactStr != "" {
		req.Compact, _ = strconv.Atoi(compactStr)
	} else {
		req.Compact = 1
	}
	if noPeerIDStr != "" {
		req.NoPeerID, _ = strconv.Atoi(noPeerIDStr)
	}
	if corruptStr != "" {
		req.Corrupt, _ = strconv.Atoi(corruptStr)
	}
	if supportCryptoStr != "" {
		req.SupportCrypto, _ = strconv.Atoi(supportCryptoStr)
	}
	if redundantStr != "" {
		req.Redundant, _ = strconv.Atoi(redundantStr)
	}

	if req.IP == "" {
		req.IP = c.Request.RemoteAddr
	}

	switch req.Event {
	case "started", "":
		peer := &model.Peer{
			PeerID: req.PeerID,
			IP:     req.IP,
			Port:   req.Port,
		}

		err := dao.Peer.AddPeer(ctx, hexString, peer)

		if err != nil {
			c.Writer.WriteString("Failed to add peer")
			return
		}
	case "stopped":
		err := dao.Peer.RemovePeer(ctx, hexString, req.PeerID)
		if err != nil {
			c.Writer.WriteString("Failed to remove peer")
			return
		}
	case "completed":
		peer := &model.Peer{
			PeerID:   req.PeerID,
			IP:       req.IP,
			Port:     req.Port,
			LastSeen: time.Now(),
			Status:   0,
		}

		err := dao.Peer.AddPeer(ctx, hexString, peer)
		if err != nil {
			c.Writer.WriteString("Failed to add peer")
			return
		}
	default:
	}

	peers, err := dao.Peer.GetPeers(ctx, hexString, req.NumWant)
	var completed, incompleted int
	for _, peer := range peers {
		if peer.Status == 0 {
			completed++
		}
	}

	responseStruct := &AnnounceResponse{
		Interval:   1800,
		Complete:   completed,
		Incomplete: incompleted,
	}

	if req.Compact == 1 {
		var peersData []byte
		for _, peer := range peers {
			ip := net.ParseIP(peer.IP).To4()
			if ip == nil {
				continue
			}
			portBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(portBytes, uint16(peer.Port))
			peersData = append(peersData, ip...)
			peersData = append(peersData, portBytes...)
		}
		responseStruct.Peers = peersData
	} else {
		var peersList []Peer
		for _, peer := range peers {
			if req.NoPeerID == 1 {
				peer.PeerID = ""
			}
			var _peer = &Peer{
				PeerID: peer.PeerID,
				IP:     peer.IP,
				Port:   peer.Port,
			}

			peersList = append(peersList, *_peer)
		}
		responseStruct.Peers = peersList
	}

	encodedResp, err := bencode.EncodeBytes(responseStruct)
	if err != nil {
		c.Writer.WriteString("Failed to encode response")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(encodedResp)
}
