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
	SettingKeyTrackerBonusFormula          = "tracker.bonus_formula"
	SettingKeyTrackerList                  = "tracker.list"
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
	TrackerBonusFormula           *Item[string]
	TrackerList                   *Item[string]
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
		TrackerBonusFormula = NewProjectItem(SettingKeyTrackerBonusFormula, "魔力公式", "sqrt(uploaded)")
		TrackerList = NewProjectItem(SettingKeyTrackerList, "Tracker 列表", "")
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

func TrackerListValue(ctx context.Context) []string {
	_ = ctx
	if TrackerList == nil {
		return nil
	}

	return parseTrackerListText(TrackerList.Value())
}

func TrackerListTextValue(ctx context.Context) string {
	_ = ctx
	list := TrackerListValue(ctx)
	if len(list) == 0 {
		return ""
	}
	return strings.Join(list, "\n")
}

func parseTrackerListText(raw string) []string {
	content := strings.TrimSpace(raw)
	if content == "" {
		return nil
	}

	if strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		var list []string
		if err := json.Unmarshal([]byte(content), &list); err == nil {
			return normalizeTrackerList(list)
		}
	}

	var single string
	if err := json.Unmarshal([]byte(content), &single); err == nil {
		content = strings.TrimSpace(single)
	}

	if strings.Contains(content, "\n") || strings.Contains(content, "\r") {
		lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
		return normalizeTrackerList(lines)
	}

	return normalizeTrackerList(strings.Split(content, ","))
}

func ParseTrackerList(raw string) []string {
	return parseTrackerListText(raw)
}

func normalizeTrackerList(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		text := strings.TrimSpace(item)
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		result = append(result, text)
	}
	return result
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
	case SettingKeyTrackerBonusFormula:
		return "string", TrackerBonusFormula.Label, true
	case SettingKeyTrackerList:
		return "string", TrackerList.Label, true
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
	case SettingKeyTrackerBonusFormula:
		return TrackerBonusFormula.Value(), true
	case SettingKeyTrackerList:
		return TrackerListTextValue(ctx), true
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
	case SettingKeyTrackerBonusFormula:
		return true, TrackerBonusFormula.Update(value)
	case SettingKeyTrackerList:
		normalized := strings.Join(parseTrackerListText(value), "\n")
		return true, TrackerList.Update(normalized)
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
