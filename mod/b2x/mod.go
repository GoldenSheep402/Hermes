package b2x

import (
	"context"
	"errors"
	"github.com/Backblaze/blazer/b2"
	"github.com/GoldenSheep402/Hermes/conf"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule

	cancel context.CancelFunc
}

func (m *Mod) Name() string {
	return "b2x"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	b2Conf := conf.Get().B2
	b2Client, err := b2.NewClient(ctx, b2Conf.BucketKeyID, b2Conf.BucketKey)
	if err != nil {
		hub.Log.Error(err)
		cancel()
		return err
	}
	b2Bucket, err := b2Client.Bucket(ctx, b2Conf.BucketName)
	if err != nil {
		hub.Log.Error(err)
		cancel()
		return err
	}
	hub.Map(&b2Client, &b2Bucket)
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	var b2Client *b2.Client
	if hub.Load(&b2Client) != nil {
		return errors.New("can't load b2 client from kernel")
	}

	var b2Bucket *b2.Bucket
	if hub.Load(&b2Bucket) != nil {
		return errors.New("can't load b2 bucket from kernel")
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	m.cancel()
	return nil
}
