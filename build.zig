const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});
    const httpz = b.dependency("httpz", .{
        .target = target,
        .optimize = optimize,
    });
    const zqlite = b.dependency("zqlite", .{
        .target = target,
        .optimize = optimize,
    });
    const uuid = b.dependency("uuid", .{
        .target = target,
        .optimize = optimize,
    });

    const exe_mod = b.createModule(.{
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const exe = b.addExecutable(.{
        .name = "bingo",
        .root_module = exe_mod,
    });

    exe.root_module.addImport("httpz", httpz.module("httpz"));
    exe.root_module.addImport("uuid", uuid.module("uuid"));
    exe.linkLibC();
    exe.linkSystemLibrary("sqlite3");
    exe.root_module.addImport("zqlite", zqlite.module("zqlite"));

    b.installArtifact(exe);

    const run_cmd = b.addRunArtifact(exe);

    run_cmd.step.dependOn(b.getInstallStep());

    if (b.args) |args| {
        run_cmd.addArgs(args);
    }

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);
}
