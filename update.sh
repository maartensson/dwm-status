
git add . && 
  git commit -m update && 
  git push && 
  cd ../dwm-nix && 
  nix flake update && 
  git add ./flake.lock && 
  git commit -m update && 
  git push &&
  cd ../dotfiles && 
  nix flake update && 
  sudo nixos-rebuild switch --flake . && 
  systemctl --user status statusbar.service
