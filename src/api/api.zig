const std = @import("std");
const httpz = @import("httpz");
const zqlite = @import("zqlite");
const html = @import("html/html.zig");

pub const Api = struct {
    alloc: std.mem.Allocator,
    db: zqlite.Conn,
    pub fn dispatch(self: *Api, action: httpz.Action(*Api), req: *httpz.Request, res: *httpz.Response) !void {
        var timer = try std.time.Timer.start();

        try action(self, req, res);

        const elapsed = timer.lap();
        std.log.info("{} {s} {s}", .{ req.method, req.url.path, try timeFormatter(self.alloc, elapsed) });
    }
};

fn timeFormatter(alloc: std.mem.Allocator, t: u64) ![]const u8 {
    const minute = 1_000_000_000;
    const milisecond = 1_000_000;
    const microsecond = 1_000;
    if (t >= minute) {
        const ft: f64 = @floatFromInt(t);
        const fm: f64 = @floatFromInt(minute);
        return std.fmt.allocPrint(alloc, "{d}m", .{ft / fm}) catch unreachable;
    }
    if (t >= milisecond) {
        const ft: f64 = @floatFromInt(t);
        const fm: f64 = @floatFromInt(milisecond);
        return std.fmt.allocPrint(alloc, "{d}ms", .{ft / fm}) catch unreachable;
    }
    if (t >= microsecond) {
        return std.fmt.allocPrint(alloc, "{d}us", .{t / microsecond}) catch unreachable;
    }
    return std.fmt.allocPrint(alloc, "{d}ns", .{t}) catch unreachable;
}

const fields = [_][]const u8{
    "you've forgot charger",
    "Feature only works on prod",
    "new jFrog token expired",
    "challenges and oportunities",
    "bug or a new feature",
    "unjustified PD call",
    "workstation WOL crash",
    "timesheets crash",
    "slack connection problems",
    "SVPN won't connect",
    "bugged migration",
    "Random exception expired",
    "VS code removed",
    "Feature works everywhere but prod",
    "work planned without details",
    "sandwiches guy already gone",
};

// <a
//href={ templ.URL(fmt.Sprintf("/api/square/click?field=%s", cell.Field)) }
//style="text-decoration: none;"
// >
//if cell.IsSet {
//<div class="bingo-cell selected">
//{ cell.Field }
//</div>
//} else {
//<div class="bingo-cell">
//{ cell.Field }
//</div>
//}
// </a>
const squereTemplate =
    \\  <a style="text-decoration: none;" 
    \\ >
    \\<div class="bingo-cell">{s}</div>
    \\</a>
;

const selectBingoCells = "SELECT * FROM bingo_history WHERE session = :session AND is_set IS NOT NULL;";
pub fn index(_: *Api, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    var body = std.ArrayList(u8).init(res.arena);
    defer body.deinit();

    try body.appendSlice(html.index);
    for (fields) |field| {
        try std.fmt.format(body.writer(), squereTemplate, .{field});
    }
    try body.appendSlice(html.index2);
    res.body = body.items;
    try res.write();

    // var rows = try api.db.rows(selectBingoCells, .{});
    // defer rows.deinit();
    // try res.chunk(html.index);
    // while (rows.next()) |row| {
    // try std.fmt.format(res.writer(), squereTemplate, .{row.text(1)});
    // std.debug.print("[dupa] row.text(1): {any}\n", .{row.text(1)});
    // }
    // try res.chunk(html.index2);
}

const selectStats = "SELECT field, count(*) as count, date(created_at) as date FROM bingo_history WHERE is_set IS NOT NULL GROUP BY field, date(created_at) ORDER BY date(created_at) DESC; ";
pub fn stats(api: *Api, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    var rows = try api.db.rows(selectStats, .{});
    defer rows.deinit();

    const writer = res.writer();
    try writer.writeByte('[');
    var dupa = false;
    while (rows.next()) |row| {
        if (dupa) {
            try writer.writeByte(',');
        }
        dupa = true;
        try res.json(.{
            .field = row.text(0),
            .count = row.int(1),
            .date = row.text(2),
        }, .{});
    }
    try writer.writeByte(']');
}
