package setting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/mod/system/dao"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
)

// EndpointList is kept to preserve typed env parsing for endpoint array settings.
type EndpointList []string

var mock = false

func MockInit() {
	mock = true
}

func MockReset() {
	mock = false
}

var ItemMap = make(map[string]updateableItem)

type updateableItem interface {
	ValueBytes() []byte
	UpdateValueBytes([]byte) error
}

type Item[T any] struct {
	Label        string `json:"label"`
	Key          string `json:"key"`
	group        string
	ItemOrder    int `json:"order"`
	CurrValue    T   `json:"value"`
	DefaultValue T   `json:"defaultValue"`

	afterUpdate []func()
}

func (o *Item[T]) Init() {
	_ = o.ensureDefaultInDB(context.Background())
	o.Value()
}

func (o *Item[T]) Order() int {
	return o.ItemOrder
}

func (o *Item[T]) UpdateValueBytes(v []byte) error {
	dst := new(T)
	if err := json.Unmarshal(v, dst); err != nil {
		return err
	}
	return o.Update(*dst)
}

func (o *Item[T]) Update(v T) error {
	vb, err := json.Marshal(v)
	if err != nil {
		return err
	}

	settingCache.Delete(o.Key)
	o.CurrValue = v
	if mock {
		o.triggerAfterUpdateHook()
		return nil
	}

	s := &model.Setting{
		Key:   o.Key,
		Value: string(vb),
		Type:  inferSettingType(any(v)),
		Desc:  o.Label,
	}
	if err := dao.Setting.UpdateOrCreate(context.Background(), s); err != nil {
		return err
	}

	settingCache.Delete(o.Key)
	o.triggerAfterUpdateHook()

	group := findGroup(o.group)
	if group != nil {
		group.triggerAfterUpdateDebounced(1 * time.Second)
	}
	return nil
}

func (o *Item[T]) value(data datatypes.JSON) (T, error) {
	if len(data) == 0 {
		o.CurrValue = o.DefaultValue
		return o.CurrValue, nil
	}

	if err := json.Unmarshal(data, &o.CurrValue); err == nil {
		return o.CurrValue, nil
	}

	legacyValue, legacyErr := o.parseLegacyValue(strings.TrimSpace(string(data)))
	if legacyErr != nil {
		o.CurrValue = o.DefaultValue
		return o.CurrValue, legacyErr
	}
	o.CurrValue = legacyValue
	return o.CurrValue, nil
}

func (o *Item[T]) parseLegacyValue(raw string) (T, error) {
	switch any(o.DefaultValue).(type) {
	case string:
		return any(raw).(T), nil
	case bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(v).(T), nil
	case int:
		v, err := strconv.Atoi(raw)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(v).(T), nil
	case int32:
		v, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(int32(v)).(T), nil
	case int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(v).(T), nil
	case float32:
		v, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(float32(v)).(T), nil
	case float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return o.DefaultValue, err
		}
		return any(v).(T), nil
	default:
		return o.DefaultValue, errors.New("legacy value parse failed")
	}
}

func (o *Item[T]) UpdateValue() {
	o.Value()
	o.triggerAfterUpdateHook()
	group := findGroup(o.group)
	if group != nil {
		group.triggerAfterUpdateDebounced(1 * time.Second)
	}
}

func (o *Item[T]) ValueBytes() (vb []byte) {
	defer func() {
		settingCache.Set(o.Key, vb, cache.NoExpiration)
	}()

	if v, ok := settingCache.Get(o.Key); ok {
		if vv, ok2 := v.([]byte); ok2 {
			return vv
		}
	}

	if mock {
		vb, _ = json.Marshal(o.DefaultValue)
		return vb
	}

	s, err := dao.Setting.GetByKey(context.Background(), o.Key)
	if err != nil {
		vb, _ = json.Marshal(o.DefaultValue)
		return vb
	}

	raw := strings.TrimSpace(s.Value)
	if raw == "" {
		vb, _ = json.Marshal(o.DefaultValue)
		return vb
	}
	return []byte(raw)
}

func (o *Item[T]) Value() T {
	value, _ := o.value(o.ValueBytes())
	return value
}

func (o *Item[T]) AfterUpdateHook(f func()) {
	o.afterUpdate = append(o.afterUpdate, f)
}

func (o *Item[T]) triggerAfterUpdateHook() {
	var wg sync.WaitGroup
	for _, f := range o.afterUpdate {
		wg.Add(1)
		go func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					zap.S().Errorf("trigger after update hook of item %s error: %v", o.Key, err)
				}
			}()
			defer wg.Done()
			f()
		}(f)
	}
	wg.Wait()
}

func (o *Item[T]) ensureDefaultInDB(ctx context.Context) error {
	if mock || dao.Setting.DB() == nil {
		return nil
	}
	_, err := dao.Setting.GetByKey(ctx, o.Key)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	vb, marshalErr := json.Marshal(o.DefaultValue)
	if marshalErr != nil {
		return marshalErr
	}

	return dao.Setting.UpdateOrCreate(ctx, &model.Setting{
		Key:   o.Key,
		Value: string(vb),
		Type:  inferSettingType(any(o.DefaultValue)),
		Desc:  o.Label,
	})
}

func inferSettingType(v any) string {
	switch v.(type) {
	case bool:
		return "bool"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "int"
	case float32, float64:
		return "float"
	case string:
		return "string"
	default:
		return "json"
	}
}

func NewItem[T any](key, label string, defaultValue T) *Item[T] {
	if key == "" {
		panic("empty key")
	}
	defer func() { currGroup.currChildrenOrder++ }()
	return NewItemWithGroupAndOrder(key, label, currGroup.Path, currGroup.currChildrenOrder, defaultValue)
}

func NewProjectItem[T any](key, label string, defaultValue T) *Item[T] {
	if key == "" {
		panic("empty key")
	}
	defer func() { currGroup.currChildrenOrder++ }()

	envKey := buildEnvKey(key)
	defaultValue = GetDefaultValue(envKey, defaultValue)

	item := &Item[T]{
		Label:        label,
		Key:          key,
		group:        currGroup.Path,
		ItemOrder:    currGroup.currChildrenOrder,
		DefaultValue: defaultValue,
	}
	currGroup.Children = append(currGroup.Children, item)
	if _, ok := ItemMap[item.Key]; ok {
		panic("setting item key conflict")
	}
	ItemMap[item.Key] = item
	return item
}

func GetDefaultValue[T any](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result any
	var err error

	switch any(defaultValue).(type) {
	case EndpointList:
		result = EndpointList(strings.Split(value, ","))
	case []string:
		result = strings.Split(value, ",")
	case string:
		result = value
	case bool:
		result, err = strconv.ParseBool(value)
	case int:
		var tmp int64
		tmp, err = strconv.ParseInt(value, 10, 0)
		result = int(tmp)
	case int8:
		var tmp int64
		tmp, err = strconv.ParseInt(value, 10, 8)
		result = int8(tmp)
	case int16:
		var tmp int64
		tmp, err = strconv.ParseInt(value, 10, 16)
		result = int16(tmp)
	case int32:
		var tmp int64
		tmp, err = strconv.ParseInt(value, 10, 32)
		result = int32(tmp)
	case int64:
		result, err = strconv.ParseInt(value, 10, 64)
	case uint:
		var tmp uint64
		tmp, err = strconv.ParseUint(value, 10, 0)
		result = uint(tmp)
	case uint8:
		var tmp uint64
		tmp, err = strconv.ParseUint(value, 10, 8)
		result = uint8(tmp)
	case uint16:
		var tmp uint64
		tmp, err = strconv.ParseUint(value, 10, 16)
		result = uint16(tmp)
	case uint32:
		var tmp uint64
		tmp, err = strconv.ParseUint(value, 10, 32)
		result = uint32(tmp)
	case uint64:
		result, err = strconv.ParseUint(value, 10, 64)
	case float32:
		var tmp float64
		tmp, err = strconv.ParseFloat(value, 32)
		result = float32(tmp)
	case float64:
		result, err = strconv.ParseFloat(value, 64)
	default:
		zap.S().Warn("found unsupported env var type", zap.String("envKey", key), zap.String("envValue", value))
		if _, ok := any(defaultValue).(sql.Scanner); ok {
			newResult := new(T)
			newResultAsScanner := any(newResult).(sql.Scanner)
			if err := newResultAsScanner.Scan(value); err == nil {
				return *newResult
			}
		}
		return defaultValue
	}

	if err != nil {
		return defaultValue
	}
	return result.(T)
}

func NewItemWithGroupAndOrder[T any](key, label, group string, order int, defaultValue T) *Item[T] {
	rawEnvKey := strings.TrimPrefix(group+"/"+key, "root/")
	envKey := buildEnvKey(rawEnvKey)
	defaultValue = GetDefaultValue(envKey, defaultValue)

	item := &Item[T]{
		Label:        label,
		Key:          group + ":" + key,
		group:        group,
		ItemOrder:    order,
		DefaultValue: defaultValue,
	}
	currGroup.Children = append(currGroup.Children, item)
	if _, ok := ItemMap[item.Key]; ok {
		panic("setting item key conflict")
	}
	ItemMap[item.Key] = item
	return item
}

func buildEnvKey(raw string) string {
	replacer := strings.NewReplacer("/", "_", ":", "_", ".", "_", "-", "_")
	normalized := strings.ToUpper(replacer.Replace(strings.TrimSpace(raw)))
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "HERMES_SETTING"
	}
	return "HERMES_SETTING_" + normalized
}

var (
	rootGroup = &GroupInfo{Path: "root", currChildrenOrder: 1}
	currGroup = rootGroup
)

func Root() *GroupInfo {
	return rootGroup
}

type GroupInfo struct {
	parent     *GroupInfo
	Path       string `json:"path"`
	Label      string `json:"label"`
	GroupOrder int    `json:"order"`

	currChildrenOrder int
	Children          []any `json:"children"`

	afterUpdate   []func()
	debounceMu    sync.Mutex
	debounceTimer *time.Timer
}

func (o *GroupInfo) Init() {
	for _, child := range o.Children {
		if g, ok := child.(initAble); ok {
			g.Init()
		}
	}
}

func (o *GroupInfo) SortChildren() {
	sort.Slice(o.Children, func(i, j int) bool {
		return o.Children[i].(orderAble).Order() < o.Children[j].(orderAble).Order()
	})
	for _, child := range o.Children {
		if g, ok := child.(*GroupInfo); ok {
			g.SortChildren()
		}
	}
}

func (o *GroupInfo) Order() int {
	return o.GroupOrder
}

func Group(key, name string, f func()) {
	currGroup = &GroupInfo{
		parent:     currGroup,
		Path:       currGroup.Path + "/" + key,
		Label:      name,
		GroupOrder: currGroup.currChildrenOrder,
	}
	currGroup.parent.Children = append(currGroup.parent.Children, currGroup)
	f()
	currGroup = currGroup.parent
}

func findGroup(path string) *GroupInfo {
	var walk func(g *GroupInfo) *GroupInfo
	walk = func(g *GroupInfo) *GroupInfo {
		if g.Path == path {
			return g
		}
		for _, child := range g.Children {
			next, ok := child.(*GroupInfo)
			if !ok {
				continue
			}
			if found := walk(next); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(rootGroup)
}

func FindGroupByPath(path string) *GroupInfo {
	return findGroup(path)
}

type initAble interface {
	Init()
}

type orderAble interface {
	Order() int
}

type valueUpdateAble interface {
	UpdateValue()
}

func (g *GroupInfo) AfterUpdateHook(f func()) {
	g.afterUpdate = append(g.afterUpdate, f)
}

func (g *GroupInfo) triggerAfterUpdateDebounced(timeout time.Duration) {
	g.debounceMu.Lock()
	defer g.debounceMu.Unlock()

	if g.debounceTimer != nil {
		g.debounceTimer.Stop()
	}

	g.debounceTimer = time.AfterFunc(timeout, func() {
		g.debounceMu.Lock()
		defer g.debounceMu.Unlock()
		for _, f := range g.afterUpdate {
			go func(f func()) {
				defer func() {
					if err := recover(); err != nil {
						zap.S().Errorf("trigger after update hook of group %s error: %v", g.Path, err)
					}
				}()
				f()
			}(f)
		}
		g.debounceTimer = nil
	})
}

func (g *GroupInfo) UpdateValue() {
	for _, child := range g.Children {
		switch child := child.(type) {
		case valueUpdateAble:
			child.UpdateValue()
		default:
			zap.S().Errorf("group info update value error: %v", child)
		}
	}
}
