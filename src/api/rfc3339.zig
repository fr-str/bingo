const std = @import("std");

pub const DateTime = struct {
    year: u16,
    month: u8,
    day: u8,
    hour: u8,
    minute: u8,
    second: u8,
};

pub fn fromTimestamp(ts: i64) DateTime {
    const SECONDS_PER_DAY = 86400;
    const DAYS_PER_YEAR = 365;
    const DAYS_IN_4YEARS = 1461;
    const DAYS_IN_100YEARS = 36524;
    const DAYS_IN_400YEARS = 146097;
    const DAYS_BEFORE_EPOCH = 719468;

    const seconds_since_midnight: i64 = @rem(ts, SECONDS_PER_DAY);
    var day_n: i64 = DAYS_BEFORE_EPOCH + @divTrunc(ts, SECONDS_PER_DAY);
    var temp: i64 = 0;

    // Calculate year
    temp = 4 * @divTrunc((day_n + DAYS_IN_100YEARS + 1), DAYS_IN_400YEARS) - 1;
    var year: u16 = @intCast(100 * temp);
    day_n -= DAYS_IN_100YEARS * temp + @divTrunc(temp, 4);

    temp = 4 * @divTrunc((day_n + DAYS_PER_YEAR + 1), DAYS_IN_4YEARS - 1);
    year += @intCast(temp);
    day_n -= DAYS_PER_YEAR * temp + @divTrunc(temp, 4);

    // Calculate month and day
    var month: u8 = @intCast(@divTrunc((5 * day_n + 2), 153));
    const day: u8 = @intCast(day_n - @divTrunc((@as(i64, @intCast(month)) * 153 + 2), 5 + 1));

    month += 3;
    if (month > 12) {
        month -= 12;
        year += 1;
    }

    return DateTime{ .year = year, .month = month, .day = day, .hour = @intCast(@divTrunc(seconds_since_midnight, 3600)), .minute = @intCast(@divTrunc(@mod(seconds_since_midnight, 3600), 60)), .second = @intCast(@mod(seconds_since_midnight, 60)) };
}

pub fn toRFC3339(dt: DateTime) [20]u8 {
    var buf: [20]u8 = undefined;

    // Format year
    _ = std.fmt.formatIntBuf(buf[0..4], dt.year, 10, .lower, .{ .width = 4, .fill = '0' });
    buf[4] = '-';

    // Format month
    paddingTwoDigits(buf[5..7], dt.month);
    buf[7] = '-';

    // Format day
    paddingTwoDigits(buf[8..10], dt.day);
    buf[10] = 'T';

    // Format time
    paddingTwoDigits(buf[11..13], dt.hour);
    buf[13] = ':';
    paddingTwoDigits(buf[14..16], dt.minute);
    buf[16] = ':';
    paddingTwoDigits(buf[17..19], dt.second);
    buf[19] = 'Z';

    return buf;
}

fn paddingTwoDigits(buf: *[2]u8, value: u8) void {
    switch (value) {
        0 => buf.* = "00".*,
        1 => buf.* = "01".*,
        2 => buf.* = "02".*,
        3 => buf.* = "03".*,
        4 => buf.* = "04".*,
        5 => buf.* = "05".*,
        6 => buf.* = "06".*,
        7 => buf.* = "07".*,
        8 => buf.* = "08".*,
        9 => buf.* = "09".*,
        else => _ = std.fmt.formatIntBuf(buf, value, 10, .lower, .{}),
    }
}

test "rfc3339 formatting" {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const ts: u64 = 1640995200; // January 1, 2022 00:00:00 UTC
    const dt = fromTimestamp(ts);
    const rfc3399_str = toRFC3339(dt);

    try std.testing.expectEqualStrings("2022-01-01T00:00:00Z", rfc3399_str);
}
