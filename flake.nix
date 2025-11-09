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

        systemd.services.wg-helper = {
          description = "WireGuard helper service";
          after = [ "network.target" ];
          serviceConfig = {
            ExecStart = "${self.packages.${pkgs.system}.default}/bin/wg";
            Restart = "always";

             # Where the socket goes (shared for everyone)
            RuntimeDirectory = "wg-helper";
            RuntimeDirectoryMode = "0777";           Type = "simple";

            DynamicUser = true;
            AmbientCapabilities = [ "CAP_NET_ADMIN" "CAP_NET_RAW" ];
            CapabilityBoundingSet = [ "CAP_NET_ADMIN" "CAP_NET_RAW" ];
            NoNewPrivileges = false;

            # Optional hardening
            PrivateTmp = true;
            ProtectSystem = "strict";
            ProtectHome = true;

            Environment = [
              "SOCKET_PATH=/run/wg-helper/wg-helper.sock"
            ];
          };
        };
        systemd.user.services.statusbar = {
          description = "DWM statusbar";
          wantedBy = ["default.target"];
          after = ["graphical-session.target"];
          serviceConfig = {
            ExecStart = "${self.packages.${pkgs.system}.default}/bin/statusbar";
            Restart = "always";
            RestartSec = "5s";
            Type = "simple";
            Environment = [
              "SOCKET_PATH=/run/wg-helper/wg-helper.sock"
            ];
          };
        };
      };
    };
  };
}

