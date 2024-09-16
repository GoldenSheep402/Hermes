package clsLog

import (
	tencentcloud_cls_sdk_go "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"go.uber.org/zap/zapcore"
)

type ZapHook struct {
	client  *tencentcloud_cls_sdk_go.AsyncProducerClient
	topicID string
	level   zapcore.LevelEnabler
}

func (h *ZapHook) Hook(entry zapcore.Entry) error {
	if !h.level.Enabled(entry.Level) {
		return nil
	}

	l := tencentcloud_cls_sdk_go.NewCLSLog(entry.Time.Unix(),
		map[string]string{
			"level":   entry.Level.String(),
			"message": entry.Message,
		},
	)
	err := h.client.SendLog(h.topicID, l, nil)
	return err
}

func NewZapHook(clientConfig *tencentcloud_cls_sdk_go.AsyncProducerClientConfig, logLevel zapcore.LevelEnabler, topicID string) (*ZapHook, error) {
	client, err := tencentcloud_cls_sdk_go.NewAsyncProducerClient(clientConfig)
	if err != nil {
		return nil, err
	}
	client.Start()
	// todo: client close gracefully
	//client.Close(int64(time.Second / time.Millisecond))
	return &ZapHook{
		client:  client,
		topicID: topicID,
		level:   logLevel,
	}, nil
}
