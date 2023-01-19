local jf = import "github.com/zeet-dev/jsonnet-filer-lib/main.libsonnet";
local klibs = {
  "1.24": import "github.com/jsonnet-libs/k8s-libsonnet/1.24/main.libsonnet",
};

{
  local k = klibs[$.config_.kube.versionShort],

  config_: {
    kubernetes: {
      // this probably needs to be passed in from the outside
      architecture: "arm64",
      clusterName: "personal",
      // we specify the patch version because it's used for other tooling
      // like ASDF to get the kubectl version to install
      versionLong: '1.24.7',
      versionShort: std.splitLimitR(self.versionLong, ".", 1)[0]
    },
  },
  deployment: jf.File("gen/kube/deployment.yml", {}),

  k3dConfig: jf.File("gen/local/k3d.yml", {
    apiVersion: "k3d.io/v1alpha4",
    kind: "Simple",
    metadata+: {
      name: $.config_.kubernetes.clusterName,
    },
    imageSuffix_:: (
    if $.config_.kubernetes.architecture == "arm64"
    then "-arm64"
    else ""
    ),
    image: "rancher/k3s:v" + $.config_.kubernetes.versionLong + "-k3s1" + self.imageSuffix_,
  }),

  taskFile: jf.File("Taskfile.yml", {
    version: "3",
    tasks+: {
      local clusterName = $.config_.kubernetes.clusterName,
      local k3d = {
        // args could be checked for whether it receives a string or an array and act accordingly
        // commands in general could be abstracted into a separate library
        New(name, args=[]):: std.rstripChars("k3d cluster %s %s" % [name, std.join(" ", args)], " "),
        WithConfig():: "--config %s" % [$.k3dConfig.metadata.name],
      },

      "k3d:up": {
        status: [k3d.New("get", [clusterName])],
        cmds: [k3d.New("create", [k3d.WithConfig()])],
      },

      "k3d:down": {
        cmds: [k3d.New("delete", [k3d.WithConfig()])],
      },
    },
  })
}
