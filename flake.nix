{
  description = "DWM statusbar service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, nixpkgs, flake-utils, ...}:
  flake-utils.lib.eachDefaultSystem (system: {
    packages.default = nixpkgs.legacyPackages.${system}.buildGoModule {
      pname = "dwm-status";
      version = "0.0.1";
      src = ./.;
      vendorHash = "sha256-KrYjCyPHbQxv+e/FObxwabmgzNmmAxFExNuRxX5rqL0=";
    };

    apps.default = {
      type = "app";
      program = "${self.packages.${system}.default}/bin/dwm-status";
    };
  }) // {
    nixosModules.default = {config, lib, pkgs, ...}: {
      options.services.statusbar = {
        enable = lib.mkEnableOption "Enable dwm-statusbar";
      };

      config = lib.mkIf config.services.statusbar.enable {
        systemd.user.services.statusbar = {
          description = "DWM statusbar";
          wantedBy = ["default.target"];
          after = ["graphical-session.target"];
          serviceConfig = {
            ExecStart = "${self.packages.${pkgs.system}.default}/bin/statusbar";
            Restart = "always";
            RestartSec = "5s";
            Type = "simple";

            # Use dynamic, unprivileged user
            User = "statusbar";
            Group = "statusbar";

            # Give only the capabilities you need
            AmbientCapabilities = [ "CAP_NET_ADMIN" "CAP_NET_RAW" ];
            CapabilityBoundingSet = [ "CAP_NET_ADMIN" "CAP_NET_RAW" ];

            # Protect the service
            NoNewPrivileges = true;
            PrivateTmp = true;
            ProtectSystem = "strict";
            ProtectHome = true;
            ProtectKernelModules = true;
            ProtectKernelTunables = true;
            PrivateDevices = true;
            RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
          };
        };
      };
    };
  };
}

