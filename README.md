# Terminal Local WebService


## Installation
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


## Description
Go webservice that allows user to restart and shutdown terminal, make screenshot and setup terminal from user interface

## Additional information
* password is 2011
* autologin to windows is preferred
* webservice is running on port 8000
* special page /RestartBrowser restarts only browser (not accessible from user interface, use remotely)
* timer is set to 20 seconds (right top corner)
    * if server ip address has running webservice, counter starts to decrease every second and the displays the page from server
    * if no webservice is running on server, 20 seconds remains

    
www.zapsi.eu © 2020
