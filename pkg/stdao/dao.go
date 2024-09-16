package stdao

import (
	"context"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"gorm.io/gorm"
)

func SetTxToCtx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey.DbTransaction, tx)
}

type Std[T any] struct {
	db *gorm.DB
}

func (m *Std[T]) Init(db *gorm.DB) error {
	m.db = db
	return m.db.AutoMigrate(new(T))
}

func (m *Std[T]) DB() *gorm.DB {
	return m.db
}

func (m *Std[T]) Begin() *gorm.DB {
	return m.db.Begin()
}

func (m *Std[T]) SetTxToCtx(ctx context.Context, tx *gorm.DB) context.Context {
	return SetTxToCtx(ctx, tx)
}

func (m *Std[T]) GetTxFromCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(ctxKey.DbTransaction).(*gorm.DB)
	if !ok || tx == nil {
		return m.db
	}
	return tx
}

func (m *Std[T]) Create(ctx context.Context, t T) error {
	return m.GetTxFromCtx(ctx).WithContext(ctx).Create(t).Error
}

func (m *Std[T]) UnscopedList(ctx context.Context) (list []T, err error) {
	err = m.GetTxFromCtx(ctx).WithContext(ctx).Unscoped().Find(&list).Error
	return
}

func (m *Std[T]) List(ctx context.Context) (list []T, err error) {
	err = m.GetTxFromCtx(ctx).WithContext(ctx).Find(&list).Error
	return
}

func (m *Std[T]) UpdateAll(ctx context.Context, t T) error {
	return m.UpdateAllResult(ctx, t).Error
}

func (m *Std[T]) UpdateAllResult(ctx context.Context, t T) *gorm.DB {
	return m.GetTxFromCtx(ctx).WithContext(ctx).Select("*").Updates(t)
}

func (m *Std[T]) Update(ctx context.Context, t T) error {
	return m.GetTxFromCtx(ctx).WithContext(ctx).Save(t).Error
}

func (m *Std[T]) Delete(ctx context.Context, t T) *gorm.DB {
	return m.GetTxFromCtx(ctx).WithContext(ctx).Delete(t)
}
