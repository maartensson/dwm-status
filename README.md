# 📈 Statusbar for DWM

dwm-status is a statusbar for dwm, written in golang. 

The state of the bar can be modified in runtime over http. For instance show date instead of time or toggle text stream.

The bar can show some standard information about the computer, but it also has an infinite rolling text field which is used to show answers from LLM

# Code architechture

The project is divided into modules, so you can add and remove the modules you want in the statusbar.

## ⚡️Some features

- Modules
  - [x] Show ip of default route (with special permissions: `sudo setcap cap_net_admin=+ep ./statusbar`)
  - [x] Volume from pulseaudio
  - [x] Battery
  - [x] Weather from Wttr (if internet)
  - [x] Brightness
  - [x] Disk usage
  - [x] Date and time
  - [x] Text stream 
- [x] Toggle datetime field state
- [x] Toggle text stream visibility

## 🔧 TODO's

### Statusbar
- [ ] Make moving visible warning thingy when low battery
- [ ] Try to not use library for pulseaudio (make it yourself)
- [ ] Change speed of text in runtime
- [ ] Test in on something else than DWM
- [ ] Write tests and coverage
- [ ] Make UI Run function receive a dynamic 

### AI Service
- [ ] Remove code formatting 
- [ ] Limit text length
- [ ] Make ai know system information
- [ ] Make ai access local database of relevant stuff
- [ ] Write tests (find out what it the best way to write tests when you are interacting with the os)
