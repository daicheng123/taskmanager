package utils

import (
	"fmt"
	"taskmanager/internal/models/common"
	"time"
)

const (
	secondPerMinute = 60
	secondPerHour   = secondPerMinute * 60
	secondPerDay    = secondPerHour * 24
)

func CalcDuration(s int) string {
	day := s / secondPerDay
	hour := (s - day*secondPerDay) / secondPerHour
	minute := (s - day*secondPerDay - hour*secondPerHour) / secondPerMinute
	second := s - day*secondPerDay - hour*secondPerHour - minute*60
	if day == 0 {
		if hour == 0 {
			if minute == 0 {
				return fmt.Sprintf("%.2d秒", second)
			}
			return fmt.Sprintf("%.2d分%.2d秒", minute, second)
		}
		return fmt.Sprintf("%.2d小时%.2d分%.2d秒", hour, minute, second)
	}
	return fmt.Sprintf("%.2d天%.2d小时%.2d分%.2d秒", day, hour, minute, second)
}

func CalcTaskDelta(start common.CustomTime, end time.Time) (delta string) {
	t, _ := time.ParseInLocation(common.TimeFormat, start.String(), time.Local)
	delta = CalcDuration(int(end.Sub(t).Seconds()))
	return
}
