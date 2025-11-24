package datatype

import "time"

// GetCurrentTimestamp 获取当前Unix时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTime 获取当前时间字符串，格式：2006-01-02 15:04:05
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatTimestamp 将时间戳格式化为指定格式
func FormatTimestamp(timestamp int64, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Unix(timestamp, 0).Format(layout)
}

// ParseTimeString 将时间字符串解析为时间戳
func ParseTimeString(timeStr, layout string) (int64, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// 快捷查询时间
func GetQueryTime(t ...string) (int, int) {
	queryType := "today"
	if len(t) > 0 {
		queryType = t[0]
	}

	switch queryType {
	case "today":
		today := time.Now()
		startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		return int(startOfDay.Unix()), int(endOfDay.Unix())
	case "yesterday":
		yesterday := time.Now().AddDate(0, 0, -1)
		startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		return int(startOfDay.Unix()), int(endOfDay.Unix())
	case "week":
		now := time.Now()
		weekday := now.Weekday()
		if weekday == 0 { // 周日
			weekday = 7
		}
		startOfWeek := now.AddDate(0, 0, -int(weekday-1))
		startOfDay := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
		endOfDay := startOfDay.AddDate(0, 0, 7)
		return int(startOfDay.Unix()), int(endOfDay.Unix())
	default:
		// 默认返回今天
		today := time.Now()
		startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		return int(startOfDay.Unix()), int(endOfDay.Unix())
	}
}

// 获取递减时间
func GetDayRange(daysAgo int) (int, int) {
	// 获取指定天数前的日期
	targetDate := time.Now().AddDate(0, 0, -daysAgo)

	// 计算当天的开始时间戳 (0点)
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	startTimestamp := int(startOfDay.Unix())

	// 计算当天的结束时间戳 (23:59:59)
	endOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, targetDate.Location())
	endTimestamp := int(endOfDay.Unix())

	return startTimestamp, endTimestamp
}
