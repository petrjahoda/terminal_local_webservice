let timeleft = 15;
const downloadTimer = setInterval(function () {
    let serverName = document.getElementById("server").innerText
    if (serverName.includes("offline")) {
        timeleft = 15
    } else {
        document.getElementById("server-info").innerText = "server page will be loaded in " + timeleft +" seconds";
        if (timeleft <= 0) {
            clearInterval(downloadTimer);
            console.log(serverName)
            window.open("http://" + serverName, "_self")
        }
        timeleft -= 1;

    }
}, 1000);

const networkDataSource = new EventSource('/networkdata');
networkDataSource.addEventListener('networkdata', (e) => {
    const networkdata = e.data.split(";");
    document.getElementById("ipaddress").innerHTML = networkdata[0];
    document.getElementById("mask").innerHTML = networkdata[1];
    document.getElementById("gateway").innerHTML = networkdata[2];
    document.getElementById("dhcp").innerHTML = networkdata[3];
    document.getElementById("server").innerHTML = networkdata[6];
}, false);

const leftButton = document.getElementById("left-button")
const middleButton = document.getElementById("middle-button")
const rightButton = document.getElementById("right-button")


leftButton.addEventListener('touchstart', function () {
    leftButton.style.border = "2px solid red"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid white"
    let data = {
        password: "3600"
    };
    fetch("/restart", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch(() => {
    });
}, false);

middleButton.addEventListener('touchstart', function () {
    leftButton.style.border = "2px solid white"
    middleButton.style.border = "2px solid red"
    rightButton.style.border = "2px solid white"
    window.open("/setup", "_self")
}, false);

rightButton.addEventListener('touchstart', function () {
    leftButton.style.border = "2px solid white"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid red"
    let data = {
        password: "3600"
    };
    fetch("/shutdown", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch(() => {
    });
}, false);