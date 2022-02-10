# Terminal Local WebService Rpi

* [1. Passwords](#1-passwords)
* [2. Remote administration](#2-remote-administration)

![actual screenshot](image4.png)

 

## 1. Passwords
* user `pi` with password `3600`
* setup password is `3600`

## 2. Remote administration
* screenshot at `http://<ipaddress>/screenshot`
* remote settings at `http://<ipaddress>/setup-remote`
* remote restart using javascript:
```
let data = {
  password: "3600"
};
fetch("/restart", {
  method: "POST",
  body: JSON.stringify(data)
}).then((result) => {
  console.log(result)
}).catch(() => {
});
```
* remote shutdown using javascript:
```
let data = {
  password: "3600"
};
fetch("/shutdown", {
  method: "POST",
  body: JSON.stringify(data)
}).then((result) => {
  console.log(result)
}).catch(() => {
});
```
* set dhcp using javascript:
```
let data = {
  password: "3600",
  server: server.value,         // server web, example: 192.168.86.100:82/terminal/1
};
fetch("/dhcp", {
  method: "POST",
  body: JSON.stringify(data)
}).then((result) => {
  console.log(result)
}).catch(() => {
});
```

* set static using javascript:
```
let data = {
  password: "3600",     
  ipaddress: ipaddress.value,   // ip address, example: 192.168.86.128
  mask: mask.value,             // mask, example: 255.255.255.0
  gateway: gateway.value,       // gateway, example: 192.168.86.1
  server: server.value,         // server web, example: 192.168.86.100:82/terminal/1
};
fetch("/static", {
  method: "POST",
  body: JSON.stringify(data)
}).then((result) => {
  console.log(result)
}).catch(() => {
});
```
* set only server address using javascript:
```
let data = {
  password: "3600",
  server: server.value,         // server web, example: 192.168.86.100:80/terminal/1
};
fetch("/server", {
  method: "POST",
  body: JSON.stringify(data)
}).then((result) => {
  console.log(result)
}).catch(() => {
});
```
* check if cable is connected:
```
fetch("/checkCable", {
  method: "POST",
}).then((result) => {
  result.text().then(function (data) {
    let connected = JSON.parse(data);
    console.log(connected["Result"])
    });
}).catch(() => {
});
```

© 2021 Petr Jahoda
