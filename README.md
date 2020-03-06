# Terminal Local WebService

## Installation for Windows
* copy file to terminal
    * put terminal_local_webservice.exe into c:\Zapsi folder and make it run as administrator
    * put start.bat into c:\Zapsi folder and make it run as administrator
    * put screenshot.exe into c:\Zapsi folder and make it run as administrator
    * put html folder into c:\Zapsi folder
    * put css folder into c:\Zapsi folder
    * put js folder into c:\Zapsi folder
* create shortcut of start.bat in C:\Users\Zapsi\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup folder (in terminal)
* create shortcut to google chrome in C:\Users\Zapsi\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup folder (in terminal)
    * use this "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --kiosk --disable-pinch --app=http://localhost:8000" 


## Installation for Linux, use Solus
* install Solus on AsusPro
* update Solus to latest version 
* add autologin for current user
* disable sleeping, allow dimming, disable screensaver
* install chrome
* install open ssh server and enable it
* install maim (for screenshots)
* set wallpaper to #2B2B2B set taskbar transparent
* remove everything from taskbar
* add chrome to startup `google-chrome-stable http://localhost --window-size=1920,1080 --start-fullscreen --kiosk --incognito --noerrdialogs --disable-translate --no-first-run --fast --fast-start --disable-infobars --disable-features=TranslateUI --disk-cache-dir=/dev/null  --password-store=basic --disable-pinch --overscroll-history-navigation=0` 
* setup nocursor in `/usr/share/lightdm/lightdm.conf.d/` in preferred config file, enable and add`xserver-command = X -nocursor`
* copy everything from folder linux to /home/{user}/ directory
* make it run as a service, according to `https://medium.com/tarkalabs-til/making-your-go-service-systemd-friendly-2ec1c9a702c7`
* test the service
* reboot

## Description
Go webservice that allows user to restart and shutdown terminal, make screenshot and setup terminal from user interface

## Additional information
* password is 2011
* autologin to windows is preferred
* webservice is running on port 80
* timer is set to 20 seconds (right top corner)
    * if server ip address has running webservice, counter starts to decrease every second and the displays the page from server
    * if no webservice is running on server, 20 seconds remains

    
www.zapsi.eu © 2020
