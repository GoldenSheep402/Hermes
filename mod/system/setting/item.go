package setting

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

const (
	SettingKeySiteName                     = "site.name"
	SettingKeySiteMaintenanceMode          = "site.maintenance_mode"
	SettingKeySiteMarquee                  = "site.marquee"
	SettingKeyInviteOpenRegistration       = "invite.open_registration"
	SettingKeyInviteOnly                   = "invite.only"
	SettingKeyInviteGlobalMessage          = "invite.global_message"
	SettingKeyInviteEmailVerify            = "invite.email_verification_required"
	SettingKeyTrackerAnnounceInterval      = "tracker.announce_interval"
	SettingKeyTrackerFlushInterval         = "tracker.flush_interval"
	SettingKeyTrackerFlushBatchSize        = "tracker.flush_batch_size"
	SettingKeyTrackerGlobalFreeleech       = "tracker.global_freeleech"
	SettingKeyTrackerFreeleechCountdownHrs = "tracker.freeleech_countdown_hours"
	SettingKeyTrackerAnnounceURL           = "tracker.announce_url"
	SettingKeyBonusEnabled                 = "bonus.enabled"
	SettingKeyBonusMultiplier              = "bonus.multiplier"
	SettingKeyBonusUploadPointsPerGiB      = "bonus.upload_points_per_gib"
	SettingKeyBonusInvitePoints            = "bonus.invite_points"
	SettingKeyAuthSMTPEnable               = "auth.smtp_enable"
)

var (
	SiteName            *Item[string]
	SiteMaintenanceMode *Item[bool]
	SiteMarquee         *Item[string]
)

var (
	InviteOpenRegistration          *Item[bool]
	InviteOnly                      *Item[bool]
	InviteGlobalMessage             *Item[string]
	InviteEmailVerificationRequired *Item[bool]
)

var (
	TrackerAnnounceInterval       *Item[int]
	TrackerFlushInterval          *Item[int]
	TrackerFlushBatchSize         *Item[int]
	TrackerGlobalFreeleech        *Item[bool]
	TrackerFreeleechCountdownHour *Item[int]
	TrackerAnnounceURL            *Item[string]
)

var (
	BonusEnabled            *Item[bool]
	BonusMultiplier         *Item[float64]
	BonusUploadPointsPerGiB *Item[int]
	BonusInvitePoints       *Item[int]
)

var (
	AuthSMTPEnable *Item[bool]
)

func init() {
	Group("site", "站点设置", func() {
		SiteName = NewProjectItem(SettingKeySiteName, "站点名称", "Hermes")
		SiteMaintenanceMode = NewProjectItem(SettingKeySiteMaintenanceMode, "维护模式", false)
		SiteMarquee = NewProjectItem(SettingKeySiteMarquee, "全站公告", "")
	})

	Group("invite", "注册与邀请", func() {
		InviteOpenRegistration = NewProjectItem(SettingKeyInviteOpenRegistration, "允许开放注册", true)
		InviteOnly = NewProjectItem(SettingKeyInviteOnly, "仅邀请码注册", false)
		InviteGlobalMessage = NewProjectItem(SettingKeyInviteGlobalMessage, "注册页提示", "")
		InviteEmailVerificationRequired = NewProjectItem(SettingKeyInviteEmailVerify, "注册需要邮箱验证", false)
	})

	Group("tracker", "Tracker", func() {
		TrackerAnnounceInterval = NewProjectItem(SettingKeyTrackerAnnounceInterval, "Announce间隔(秒)", 1800)
		TrackerFlushInterval = NewProjectItem(SettingKeyTrackerFlushInterval, "流量落库间隔(秒)", 60)
		TrackerFlushBatchSize = NewProjectItem(SettingKeyTrackerFlushBatchSize, "流量落库批次大小", 200)
		TrackerGlobalFreeleech = NewProjectItem(SettingKeyTrackerGlobalFreeleech, "全站Freeleech", false)
		TrackerFreeleechCountdownHour = NewProjectItem(SettingKeyTrackerFreeleechCountdownHrs, "限时Freeleech倒计时(小时)", 0)
		TrackerAnnounceURL = NewProjectItem(SettingKeyTrackerAnnounceURL, "Tracker Announce URL", "")
	})

	Group("bonus", "魔力", func() {
		BonusEnabled = NewProjectItem(SettingKeyBonusEnabled, "启用做种魔力", true)
		BonusMultiplier = NewProjectItem(SettingKeyBonusMultiplier, "魔力系数", 1.0)
		BonusUploadPointsPerGiB = NewProjectItem(SettingKeyBonusUploadPointsPerGiB, "兑换上传量单价(魔力/GiB)", 300)
		BonusInvitePoints = NewProjectItem(SettingKeyBonusInvitePoints, "兑换邀请单价(魔力)", 50000)
	})

	Group("auth", "认证", func() {
		AuthSMTPEnable = NewProjectItem(SettingKeyAuthSMTPEnable, "启用SMTP发信", false)
	})
}

func Init(ctx context.Context) error {
	_ = ctx
	rootGroup.SortChildren()
	rootGroup.Init()
	return nil
}

func TrackerFlushIntervalValue(ctx context.Context) int {
	_ = ctx
	if TrackerFlushInterval == nil {
		return 60
	}
	v := TrackerFlushInterval.Value()
	if v <= 0 {
		return 60
	}
	return v
}

func TrackerFlushBatchSizeValue(ctx context.Context) int {
	_ = ctx
	if TrackerFlushBatchSize == nil {
		return 200
	}
	v := TrackerFlushBatchSize.Value()
	if v <= 0 {
		return 200
	}
	return v
}

func InviteEmailVerifyRequiredValue(ctx context.Context) bool {
	_ = ctx
	if InviteEmailVerificationRequired == nil {
		return false
	}
	return InviteEmailVerificationRequired.Value()
}

func AuthSMTPEnableValue(ctx context.Context) bool {
	_ = ctx
	if AuthSMTPEnable == nil {
		return false
	}
	return AuthSMTPEnable.Value()
}

func TrackerAnnounceURLValue(ctx context.Context) string {
	_ = ctx
	if TrackerAnnounceURL == nil {
		return ""
	}
	return strings.TrimSpace(TrackerAnnounceURL.Value())
}

func TrackerAnnounceIntervalValue(ctx context.Context) int {
	_ = ctx
	if TrackerAnnounceInterval == nil {
		return 1800
	}
	v := TrackerAnnounceInterval.Value()
	if v < 60 {
		return 1800
	}
	return v
}

func TrackerGlobalFreeleechValue(ctx context.Context) bool {
	_ = ctx
	if TrackerGlobalFreeleech == nil {
		return false
	}
	return TrackerGlobalFreeleech.Value()
}

func TrackerFreeleechCountdownHoursValue(ctx context.Context) int {
	_ = ctx
	if TrackerFreeleechCountdownHour == nil {
		return 0
	}
	v := TrackerFreeleechCountdownHour.Value()
	if v < 0 {
		return 0
	}
	return v
}

func BonusEnabledValue(ctx context.Context) bool {
	_ = ctx
	if BonusEnabled == nil {
		return true
	}
	return BonusEnabled.Value()
}

func BonusMultiplierValue(ctx context.Context) float64 {
	_ = ctx
	if BonusMultiplier == nil {
		return 1
	}
	v := BonusMultiplier.Value()
	if v <= 0 {
		return 1
	}
	return v
}

func BonusUploadPointsPerGiBValue(ctx context.Context) int {
	_ = ctx
	if BonusUploadPointsPerGiB == nil {
		return 300
	}
	v := BonusUploadPointsPerGiB.Value()
	if v <= 0 {
		return 300
	}
	return v
}

func BonusInvitePointsValue(ctx context.Context) int {
	_ = ctx
	if BonusInvitePoints == nil {
		return 50000
	}
	v := BonusInvitePoints.Value()
	if v <= 0 {
		return 50000
	}
	return v
}

// normalizeAnnounceURLInput accepts legacy multi-line tracker.list values and keeps the first URL.
func normalizeAnnounceURLInput(raw string) string {
	content := strings.TrimSpace(raw)
	if content == "" {
		return ""
	}

	if strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		var list []string
		if err := json.Unmarshal([]byte(content), &list); err == nil {
			for _, item := range list {
				item = strings.TrimSpace(item)
				if item != "" && !strings.HasPrefix(item, "#") {
					return item
				}
			}
			return ""
		}
	}

	var single string
	if err := json.Unmarshal([]byte(content), &single); err == nil {
		content = strings.TrimSpace(single)
	}

	content = strings.ReplaceAll(content, "\r\n", "\n")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, ","); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if line != "" {
			return line
		}
	}
	return ""
}

func ManagedTypeAndDesc(key string) (string, string, bool) {
	switch key {
	case SettingKeySiteName:
		return "string", SiteName.Label, true
	case SettingKeySiteMaintenanceMode:
		return "bool", SiteMaintenanceMode.Label, true
	case SettingKeySiteMarquee:
		return "string", SiteMarquee.Label, true
	case SettingKeyInviteOpenRegistration:
		return "bool", InviteOpenRegistration.Label, true
	case SettingKeyInviteOnly:
		return "bool", InviteOnly.Label, true
	case SettingKeyInviteGlobalMessage:
		return "string", InviteGlobalMessage.Label, true
	case SettingKeyInviteEmailVerify:
		return "bool", InviteEmailVerificationRequired.Label, true
	case SettingKeyTrackerAnnounceInterval:
		return "int", TrackerAnnounceInterval.Label, true
	case SettingKeyTrackerFlushInterval:
		return "int", TrackerFlushInterval.Label, true
	case SettingKeyTrackerFlushBatchSize:
		return "int", TrackerFlushBatchSize.Label, true
	case SettingKeyTrackerGlobalFreeleech:
		return "bool", TrackerGlobalFreeleech.Label, true
	case SettingKeyTrackerFreeleechCountdownHrs:
		return "int", TrackerFreeleechCountdownHour.Label, true
	case SettingKeyTrackerAnnounceURL:
		return "string", TrackerAnnounceURL.Label, true
	case SettingKeyBonusEnabled:
		return "bool", BonusEnabled.Label, true
	case SettingKeyBonusMultiplier:
		return "float", BonusMultiplier.Label, true
	case SettingKeyBonusUploadPointsPerGiB:
		return "int", BonusUploadPointsPerGiB.Label, true
	case SettingKeyBonusInvitePoints:
		return "int", BonusInvitePoints.Label, true
	case SettingKeyAuthSMTPEnable:
		return "bool", AuthSMTPEnable.Label, true
	default:
		return "", "", false
	}
}

func ManagedValueString(ctx context.Context, key string) (string, bool) {
	switch key {
	case SettingKeySiteName:
		return SiteName.Value(), true
	case SettingKeySiteMaintenanceMode:
		return strconv.FormatBool(SiteMaintenanceMode.Value()), true
	case SettingKeySiteMarquee:
		return SiteMarquee.Value(), true
	case SettingKeyInviteOpenRegistration:
		return strconv.FormatBool(InviteOpenRegistration.Value()), true
	case SettingKeyInviteOnly:
		return strconv.FormatBool(InviteOnly.Value()), true
	case SettingKeyInviteGlobalMessage:
		return InviteGlobalMessage.Value(), true
	case SettingKeyInviteEmailVerify:
		return strconv.FormatBool(InviteEmailVerificationRequired.Value()), true
	case SettingKeyTrackerAnnounceInterval:
		return strconv.Itoa(TrackerAnnounceInterval.Value()), true
	case SettingKeyTrackerFlushInterval:
		return strconv.Itoa(TrackerFlushIntervalValue(ctx)), true
	case SettingKeyTrackerFlushBatchSize:
		return strconv.Itoa(TrackerFlushBatchSizeValue(ctx)), true
	case SettingKeyTrackerGlobalFreeleech:
		return strconv.FormatBool(TrackerGlobalFreeleech.Value()), true
	case SettingKeyTrackerFreeleechCountdownHrs:
		return strconv.Itoa(TrackerFreeleechCountdownHour.Value()), true
	case SettingKeyTrackerAnnounceURL:
		return TrackerAnnounceURLValue(ctx), true
	case SettingKeyBonusEnabled:
		return strconv.FormatBool(BonusEnabledValue(ctx)), true
	case SettingKeyBonusMultiplier:
		return strconv.FormatFloat(BonusMultiplierValue(ctx), 'f', -1, 64), true
	case SettingKeyBonusUploadPointsPerGiB:
		return strconv.Itoa(BonusUploadPointsPerGiBValue(ctx)), true
	case SettingKeyBonusInvitePoints:
		return strconv.Itoa(BonusInvitePointsValue(ctx)), true
	case SettingKeyAuthSMTPEnable:
		return strconv.FormatBool(AuthSMTPEnable.Value()), true
	default:
		return "", false
	}
}

func UpdateManagedValueFromString(ctx context.Context, key, value string) (bool, error) {
	_ = ctx
	trimmed := strings.TrimSpace(value)

	switch key {
	case SettingKeySiteName:
		return true, SiteName.Update(value)
	case SettingKeySiteMaintenanceMode:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, SiteMaintenanceMode.Update(v)
	case SettingKeySiteMarquee:
		return true, SiteMarquee.Update(value)
	case SettingKeyInviteOpenRegistration:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, InviteOpenRegistration.Update(v)
	case SettingKeyInviteOnly:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, InviteOnly.Update(v)
	case SettingKeyInviteGlobalMessage:
		return true, InviteGlobalMessage.Update(value)
	case SettingKeyInviteEmailVerify:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, InviteEmailVerificationRequired.Update(v)
	case SettingKeyTrackerAnnounceInterval:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, TrackerAnnounceInterval.Update(v)
	case SettingKeyTrackerFlushInterval:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, TrackerFlushInterval.Update(v)
	case SettingKeyTrackerFlushBatchSize:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, TrackerFlushBatchSize.Update(v)
	case SettingKeyTrackerGlobalFreeleech:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, TrackerGlobalFreeleech.Update(v)
	case SettingKeyTrackerFreeleechCountdownHrs:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, TrackerFreeleechCountdownHour.Update(v)
	case SettingKeyTrackerAnnounceURL:
		return true, TrackerAnnounceURL.Update(normalizeAnnounceURLInput(value))
	case SettingKeyBonusEnabled:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, BonusEnabled.Update(v)
	case SettingKeyBonusMultiplier:
		v, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return true, err
		}
		return true, BonusMultiplier.Update(v)
	case SettingKeyBonusUploadPointsPerGiB:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, BonusUploadPointsPerGiB.Update(v)
	case SettingKeyBonusInvitePoints:
		v, err := strconv.Atoi(trimmed)
		if err != nil {
			return true, err
		}
		return true, BonusInvitePoints.Update(v)
	case SettingKeyAuthSMTPEnable:
		v, err := strconv.ParseBool(trimmed)
		if err != nil {
			return true, err
		}
		return true, AuthSMTPEnable.Update(v)
	default:
		return false, nil
	}
}
