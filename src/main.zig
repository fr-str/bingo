const std = @import("std");
const httpz = @import("httpz");
const api = @import("api/api.zig");
// const handlers = @import("api/handlers.zig");
const zqlite = @import("zqlite");

pub fn main() !void {
    std.log.info("starting...", .{});
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    const flags = zqlite.OpenFlags.Create | zqlite.OpenFlags.EXResCode;
    var db = try zqlite.open("/home/user/code/bingo/data/bingo.db", flags);
    defer db.close();

    var a = api.Api{
        .alloc = allocator,
        .db = db,
    };
    var server = try httpz.Server(*api.Api).init(allocator, .{ .port = 5882 }, &a);
    defer {
        // clean shutdown, finishes serving any live request
        server.stop();
        server.deinit();
    }

    var router = try server.router(.{});
    router.get("/", api.index, .{});
    router.get("/api/stats", api.stats, .{});

    // blocks
    try server.listen();
}
