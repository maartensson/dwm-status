{
  description = "DWM statusbar service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {nixpkgs, flake-utils, ...}:
  flake-utils.lib.eachDefaultSystem (system: let 
    pkgs = nixpkgs.legacyPackages.${system};

    statusbar = pkgs.buildGoModule {
      pname = "statusbar";
      version = "0.0.1";
      src = ./.;
      subPackages = [
        "cmd/statusbar"
        "cmd/wg"
      ];
      vendorHash = "sha256-KrYjCyPHbQxv+e/FObxwabmgzNmmAxFExNuRxX5rqL0=";
    };

    wgctl = pkgs.buildGoModule {
      pname = "wgctl";
      version = "0.0.1";
      src = ./.;
      subPackages = [
        "cmd/wgctl"
      ];
      vendorHash = "sha256-KrYjCyPHbQxv+e/FObxwabmgzNmmAxFExNuRxX5rqL0=";
    };

  in {
    packages.default = statusbar;
    packages.wgctl = wgctl;

    apps.statusbar = {
      type = "app";
      program = "${statusbar}/bin/statusbar";
    };

    apps.wgctl = {
      type = "app";
      program = "${wgctl}/bin/wgctl";
    };

    nixosModules.default = {config, lib, pkgs, ...}: {
      options.services.statusbar = {
        enable = lib.mkEnableOption "Enable dwm-statusbar";
      };

      config = lib.mkIf config.services.statusbar.enable {

        environment.systemPackages = [ wgctl ];

        systemd.services.wg-helper = {
          description = "WireGuard helper service";
          after = [ "network.target" ];
          serviceConfig = {
            ExecStart = "${statusbar}/bin/wg";
            Restart = "always";

             # Where the socket goes (shared for everyone)
            RuntimeDirectory = "wg-helper";
            RuntimeDirectoryMode = "0777";           
            Type = "simple";

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
            ExecStart = "${statusbar}/bin/statusbar";
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
  });
}

