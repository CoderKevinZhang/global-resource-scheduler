load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "options.go",
        "types.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/go-chassis/go-chassis-config",
    importpath = "github.com/go-chassis/go-chassis-config",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/go-mesh/openlogging:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
