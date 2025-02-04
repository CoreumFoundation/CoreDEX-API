package ohlc

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Test the periods lookup
func Test_AssociatedPeriods(t *testing.T) {
	// Test the lookup for 1 minute:
	period := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 1}
	cmp.Equal(AssociatedPeriods[period.ToString()], period)
	period3m := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 3}
	cmp.Equal(AssociatedPeriods[period3m.ToString()], period)
	period5m := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 5}
	cmp.Equal(AssociatedPeriods[period5m.ToString()], period)
	period15m := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 15}
	cmp.Equal(AssociatedPeriods[period15m.ToString()], period5m)
	period30m := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 30}
	cmp.Equal(AssociatedPeriods[period30m.ToString()], period15m)
	period1h := &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 1}
	cmp.Equal(AssociatedPeriods[period1h.ToString()], period30m)
	period3h := &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 3}
	cmp.Equal(AssociatedPeriods[period3h.ToString()], period1h)
	period6h := &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 6}
	cmp.Equal(AssociatedPeriods[period6h.ToString()], period3h)
	period12h := &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 12}
	cmp.Equal(AssociatedPeriods[period12h.ToString()], period6h)
	period1d := &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 1}
	cmp.Equal(AssociatedPeriods[period1d.ToString()], period12h)
	period3d := &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 3}
	cmp.Equal(AssociatedPeriods[period3d.ToString()], period1d)
	period1w := &Period{PeriodType: PeriodType_PERIOD_TYPE_WEEK, Duration: 1}
	cmp.Equal(AssociatedPeriods[period1w.ToString()], period1d)
}

func Test_ToString(t *testing.T) {
	period := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 1}
	cmp.Equal(period.ToString(), "1minute")
}

func Test_offset(t *testing.T) {
	period := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 1}
	cmp.Equal(period.offset(), 0)
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 1}
	cmp.Equal(period.offset(), 0)
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 1}
	cmp.Equal(period.offset(), 0)
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_WEEK, Duration: 1}
	i := 4 * int64(time.Minute) * 24 * 60
	cmp.Equal(period.offset(), i)
}

func Test_ToOHLCKeyTimestamp(t *testing.T) {
	period := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 1}
	var tInput int64 = 1677081282173787000
	var tMinute int64 = 1677081240000000000
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), tMinute) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), tMinute))
	}
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 3}
	var t3Minute int64 = 1677081240000000000
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t3Minute) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t3Minute))
	}
	var t5Minute int64 = 1677081000000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 5}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t5Minute) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t5Minute))
	}
	var t15Minute int64 = 1677080700000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 15}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t15Minute) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t15Minute))
	}
	var t30Minute int64 = 1677079800000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 30}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t30Minute) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t30Minute))
	}
	var tHour int64 = 1677078000000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 1}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), tHour) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), tHour))
	}

	var t3Hour int64 = 1677078000000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 3}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t3Hour) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t3Hour))
	}
	var t6Hour int64 = 1677067200000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 6}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t6Hour) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t6Hour))
	}
	var t12Hour int64 = 1677067200000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 12}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t12Hour) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t12Hour))
	}
	var tDay int64 = 1677024000000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 1}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), tDay) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), tDay))
	}
	var t3Day int64 = 1677024000000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 3}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), t3Day) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), t3Day))
	}
	var tWeek int64 = 1676851200000000000
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_WEEK, Duration: 1}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), tWeek) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), tWeek))
	}
	// week has an offset, check the offset to be correct by checking 2 specific values in the week: 1 before the offset and 1 after the offset
	tInput = tMinute + 4*int64(time.Minute)*24*60
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_WEEK, Duration: 1}
	if !cmp.Equal(period.ToOHLCKeyTimestamp(tInput), tWeek) {
		t.Error(cmp.Diff(period.ToOHLCKeyTimestamp(tInput), tWeek))
	}
}

func Test_ToMinute(t *testing.T) {
	period := &Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 1}
	if !cmp.Equal(*period.ToMinute(), *period) {
		t.Error(cmp.Diff(period.ToMinute(), 1))
	}
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_HOUR, Duration: 1}
	if !cmp.Equal(*period.ToMinute(), Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 60}) {
		t.Error(cmp.Diff(period.ToMinute(), 60))
	}
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_DAY, Duration: 1}
	if !cmp.Equal(*period.ToMinute(), Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 60 * 24}) {
		t.Error(cmp.Diff(period.ToMinute(), 60*24))
	}
	period = &Period{PeriodType: PeriodType_PERIOD_TYPE_WEEK, Duration: 1}
	if !cmp.Equal(*period.ToMinute(), Period{PeriodType: PeriodType_PERIOD_TYPE_MINUTE, Duration: 60 * 24 * 7}) {
		t.Error(cmp.Diff(period.ToMinute(), 60*24*7))
	}
}
