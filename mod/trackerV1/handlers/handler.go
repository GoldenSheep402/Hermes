package handlers

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/GoldenSheep402/Hermes/conf"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	trackerV1Dao "github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	"github.com/juanjiTech/jin"
	"github.com/zeebo/bencode"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Registry(jinE *jin.Engine) {
	trackerGroup := jinE.Group("/api/tracker")

	trackerGroup.Group("/announce").
		GET("/key/:key", AnnounceWithKey)
}

func AnnounceWithKey(c *jin.Context) {
	ctx := context.Background()
	type AnnounceRequest struct {
		InfoHash      string
		PeerID        string
		Port          int
		Uploaded      int
		Downloaded    int
		Left          int
		Event         string // omitempty
		IP            string // omitempty
		NumWant       int    // omitempty
		Key           string // omitempty
		Compact       int    // omitempty
		NoPeerID      int    // omitempty
		TrackerID     string // omitempty
		Corrupt       int    // not standard
		SupportCrypto int    // not standard
		Redundant     int    // not standard
	}

	type AnnounceResponse struct {
		FailureReason  string      `bencode:"failure reason,omitempty"`
		WarningMessage string      `bencode:"warning message,omitempty"`
		Interval       int         `bencode:"interval"`
		MinInterval    int         `bencode:"min interval,omitempty"`
		TrackerID      string      `bencode:"tracker id,omitempty"`
		Complete       int         `bencode:"complete"`
		Incomplete     int         `bencode:"incomplete"`
		Peers          interface{} `bencode:"peers"`
	}

	type Peer struct {
		PeerID string `bencode:"peer id"`
		IP     string `bencode:"ip"`
		Port   int    `bencode:"port"`
	}
	// TODO: KEY check
	key, ok := c.Params.Get("key")
	if !ok {
		message := AnnounceResponse{
			WarningMessage: "Unauthorized access",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}

	status, UID, err := trackerV1Dao.TrackerV1.CheckKey(ctx, key)
	if err != nil {
		message := AnnounceResponse{
			WarningMessage: "Unauthorized access",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}

	if status != trackerV1Values.OK {
		message := AnnounceResponse{
			WarningMessage: "You are banned",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}

	req := &AnnounceRequest{}
	infoHashDecode := c.Request.URL.Query().Get("info_hash")
	req.PeerID = c.Request.URL.Query().Get("peer_id")
	portStr := c.Request.URL.Query().Get("port")
	uploadedStr := c.Request.URL.Query().Get("uploaded")
	downloadedStr := c.Request.URL.Query().Get("downloaded")
	leftStr := c.Request.URL.Query().Get("left")
	req.Event = c.Request.URL.Query().Get("event")
	// TODO: take IP from header
	req.IP = c.Request.URL.Query().Get("ip")
	numWantStr := c.Request.URL.Query().Get("numwant")
	req.Key = c.Request.URL.Query().Get("key")
	compactStr := c.Request.URL.Query().Get("compact")
	noPeerIDStr := c.Request.URL.Query().Get("no_peer_id")
	req.TrackerID = c.Request.URL.Query().Get("trackerid")
	corruptStr := c.Request.URL.Query().Get("corrupt")
	supportCryptoStr := c.Request.URL.Query().Get("supportcrypto")
	redundantStr := c.Request.URL.Query().Get("redundant")

	infoHashBytes := []byte(infoHashDecode)
	hexString := fmt.Sprintf("%x", infoHashBytes)
	// Set hash -> TorrentID in redis
	torrentID, err := trackerV1Dao.TrackerV1.GetTorrentID(ctx, hexString)
	if err != nil {
		message := AnnounceResponse{
			WarningMessage: "Unknown info_hash",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}

	req.InfoHash = hexString

	if req.Port, err = strconv.Atoi(portStr); err != nil {
		message := AnnounceResponse{
			WarningMessage: "Invalid Args",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}
	if req.Uploaded, err = strconv.Atoi(uploadedStr); err != nil {
		message := AnnounceResponse{
			WarningMessage: "Invalid Args",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}
	if req.Downloaded, err = strconv.Atoi(downloadedStr); err != nil {
		message := AnnounceResponse{
			WarningMessage: "Invalid Args",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
		return
	}
	if req.Left, err = strconv.Atoi(leftStr); err != nil {
		message := AnnounceResponse{
			WarningMessage: "Invalid Args",
		}

		encodedResp, err := bencode.EncodeBytes(message)
		if err != nil {
			c.Writer.WriteString("Failed to encode response")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(encodedResp)
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

	// take ip from x-forwarded-for
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	var host string

	// TODO: xff has multiple ips
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		host = strings.TrimSpace(strings.Split(ips[0], ":")[0])
	} else {
		host, _, err = net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			host = c.Request.RemoteAddr
		}
	}

	allowedSubnets := conf.Get().TrackerV1.AllowedSubnets
	fmt.Printf("Allowed subnets: %v\n", allowedSubnets)
	var allowed bool

	for _, subnet := range allowedSubnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			continue
		}

		if ipNet.Contains(net.ParseIP(host)) {
			allowed = true
			break
		}
	}

	if !allowed {
		_host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			host = c.Request.RemoteAddr
		}
		req.IP = _host
	} else {
		req.IP = host
	}

	var peerStatus int
	switch req.Event {
	case "started", "":
		if req.Left == 0 && req.Uploaded == 0 && req.Left == 0 {
			peerStatus = trackerV1Values.ReadySeeding
		} else if req.Left == 0 && req.Uploaded > 0 {
			peerStatus = trackerV1Values.Seeding
		} else {
			peerStatus = trackerV1Values.Downloading
		}

		peer := &model.Peer{
			PeerID: req.PeerID,
			IP:     req.IP,
			Port:   req.Port,
		}

		fmt.Printf("Peer: %v\n", peer)
		err := dao.Peer.AddPeer(ctx, torrentID, peer, UID)

		if err != nil {
			c.Writer.WriteString("Failed to add peer")
			return
		}
	case "stopped":
		peerStatus = trackerV1Values.Stopped
		err := dao.Peer.RemovePeer(ctx, torrentID, req.PeerID, UID)
		if err != nil {
			c.Writer.WriteString("Failed to remove peer")
			return
		}
	case "completed":
		peerStatus = trackerV1Values.Completed
		peer := &model.Peer{
			PeerID:   req.PeerID,
			IP:       req.IP,
			Port:     req.Port,
			LastSeen: time.Now(),
			Status:   0,
		}

		err := dao.Peer.AddPeer(ctx, torrentID, peer, UID)
		if err != nil {
			c.Writer.WriteString("Failed to add peer")
			return
		}
	default:
		c.Writer.WriteString("Invalid event")
		return
	}

	peers, err := dao.Peer.GetPeers(ctx, torrentID, req.NumWant)
	var completed, incompleted int
	for _, peer := range peers {
		fmt.Printf("Peer: %v\n", peer)
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

	uploadBytes, err := strconv.ParseInt(uploadedStr, 10, 64)
	if err != nil {
		c.Writer.WriteString("Failed to parse uploaded")
		return
	}
	downloadBytes, err := strconv.ParseInt(downloadedStr, 10, 64)
	if err != nil {
		c.Writer.WriteString("Failed to parse downloaded")
		return
	}

	uploadMB := uploadBytes >> 20
	downloadMB := downloadBytes >> 20

	// TODO: lock on redis to avoid race condition
	go func() {
		if err := trackerV1Dao.TrackerV1.HandelDownloadAndUpload(ctx, torrentID, UID, peerStatus, uploadMB, downloadMB); err != nil {
			log.Printf("Failed to update download and upload: %v", err)
		}
	}()

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(encodedResp)
	return
}
