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
      vendorHash = "sha256-GXLUqP0n1QoXuBZs/0mPQhqNG7YJJhEg5R+g3CTEfcc=";
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
          };
        };
      };
    };
  };
}

