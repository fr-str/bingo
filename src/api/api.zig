const std = @import("std");
const httpz = @import("httpz");
const zqlite = @import("zqlite");
const html = @import("html/html.zig");
const rfc3339 = @import("rfc3339.zig");
const uuid = @import("uuid.zig").UUID;

pub const Api = struct {
    alloc: std.mem.Allocator,
    db: zqlite.Conn,
    pub fn dispatch(self: *Api, action: httpz.Action(*Api), req: *httpz.Request, res: *httpz.Response) !void {
        var timer = try std.time.Timer.start();
        // if cookies are not set, set them
        var session: [36]u8 = undefined;
        if (req.cookies().get("session") == null) {
            const s = try uuid.init();
            session = s.to_string();
            try res.setCookie("session", &session, .{});
        }
        // disable all caching
        res.header("Cache-Control", "no-cache, no-store, must-revalidate");
        res.header("Pragma", "no-cache");
        res.header("Expires", "0");

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

const cell = struct {
    field: []const u8 = "",
    is_set: bool = false,
};
const fields = [_]cell{
    cell{ .field = "you've forgot charger", .is_set = false },
    cell{ .field = "Feature only works on prod", .is_set = false },
    cell{ .field = "new jFrog token expired", .is_set = false },
    cell{ .field = "challenges and oportunities", .is_set = false },
    cell{ .field = "bug or a new feature", .is_set = false },
    cell{ .field = "unjustified PD call", .is_set = false },
    cell{ .field = "workstation WOL crash", .is_set = false },
    cell{ .field = "timesheets crash", .is_set = false },
    cell{ .field = "slack connection problems", .is_set = false },
    cell{ .field = "SVPN won't connect", .is_set = false },
    cell{ .field = "bugged migration", .is_set = false },
    cell{ .field = "Random exception expired", .is_set = false },
    cell{ .field = "VS code removed", .is_set = false },
    cell{ .field = "Feature works everywhere but prod", .is_set = false },
    cell{ .field = "work planned without details", .is_set = false },
    cell{ .field = "sandwiches guy already gone", .is_set = false },
};

const squereTemplate =
    \\  <a href="/api/square/click?field={s}" style="text-decoration: none;"><div class="bingo-cell">{s}</div></a>
    \\
;
const squereTemplateSelected =
    \\  <a href="/api/square/click?field={s}" style="text-decoration: none;"><div class="bingo-cell selected">{s}</div></a>
    \\
;

const selectBingoCells = "SELECT field, is_set FROM bingo_history WHERE session = ? AND is_set IS NOT NULL;";
fn isValid(c: u8) bool {
    return switch (c) {
        'A'...'Z', 'a'...'z', '0'...'9', '-', '_', '.' => true,
        else => false,
    };
}

pub fn index(api: *Api, req: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    var body: std.ArrayListUnmanaged(u8) = .empty;
    defer body.deinit(res.arena);

    const session = req.cookies().get("session") orelse {
        res.status = 308;
        res.header("Location", "/");
        return;
    };

    const session2 = try bingoSession(req.arena, session);

    // copy fields
    var tmpFields = fields;
    var rows = try api.db.rows(selectBingoCells, .{session2});
    defer rows.deinit();
    // for row set is_set
    while (rows.next()) |row| {
        for (&tmpFields) |*field| {
            if (std.mem.eql(u8, field.field, row.text(0))) {
                field.is_set = row.int(1) != 0;
            }
        }
    }
    try body.appendSlice(res.arena, html.index);

    // write bingo cells
    var encodedHref: std.ArrayListUnmanaged(u8) = .empty;
    defer encodedHref.deinit(res.arena);
    for (tmpFields) |field| {
        encodedHref.clearRetainingCapacity();
        try std.Uri.Component.percentEncode(
            encodedHref.writer(res.arena),
            field.field,
            isValid,
        );
        if (!field.is_set) {
            try std.fmt.format(body.writer(res.arena), squereTemplate, .{ encodedHref.items, field.field });
            continue;
        }
        try std.fmt.format(body.writer(res.arena), squereTemplateSelected, .{ encodedHref.items, field.field });
    }

    try body.appendSlice(res.arena, html.index2);
    res.body = body.items;
    try res.write();
}

const squereInsert = "INSERT INTO bingo_history (id, field,session, is_set, created_at, updated_at) VALUES (:id, :field,:session, :is_set, :created_at, :updated_at) ON CONFLICT (id) DO UPDATE SET is_set = excluded.is_set, updated_at = excluded.updated_at;";

const squereSelect = "SELECT field, is_set FROM bingo_history WHERE id = :id AND session = :session;";
pub fn squareClick(api: *Api, req: *httpz.Request, res: *httpz.Response) !void {
    const query = try req.query();
    const field = query.get("field") orelse {
        std.debug.print("[dupa] field is null\n", .{});
        res.status = 400;
        return;
    };

    const session = req.cookies().get("session") orelse {
        std.debug.print("[dupa] session is null\n", .{});
        res.status = 400;
        return;
    };

    const id = try bingoIDFormat(req.arena, field, session);
    const session2 = try bingoSession(req.arena, session);
    const row = api.db.row(squereSelect, .{
        id,
        session2,
    }) catch |err| {
        std.debug.print("[dupa] err: {s}\n", .{@errorName(err)});
        return;
    };

    var cell2: cell = .{};
    if (row) |r| {
        defer r.deinit();
        cell2.field = try std.mem.Allocator.dupe(req.arena, u8, r.text(0));
        cell2.is_set = r.int(1) != 0;
    }

    const rfc3339_str = rfc3339.toRFC3339(rfc3339.fromTimestamp(std.time.timestamp()));

    try api.db.exec(squereInsert, .{
        id,
        field,
        session2,
        !cell2.is_set,
        &rfc3339_str,
        &rfc3339_str,
    });
    res.status = 308;
    res.header("Location", "/");
}

fn bingoIDFormat(alloc: std.mem.Allocator, field: []const u8, session: []const u8) ![]const u8 {
    return std.fmt.allocPrint(alloc, "{s}/{s}/{d}", .{ session, field, dayStamp(std.time.timestamp()) });
}

fn dayStamp(t: i64) i64 {
    const timestamp = @divTrunc(t, 86400);
    return timestamp * 86400;
}

fn bingoSession(alloc: std.mem.Allocator, session: []const u8) ![]const u8 {
    return std.fmt.allocPrint(alloc, "{s}/{d}", .{ session, dayStamp(std.time.timestamp()) });
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
