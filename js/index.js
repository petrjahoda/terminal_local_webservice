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
    fetch("/restart", {
        method: "POST",
    }).then(() => {
    }).catch(() => {
    });
}, false);

middleButton.addEventListener('click', function () {
    leftButton.style.border = "2px solid white"
    middleButton.style.border = "2px solid red"
    rightButton.style.border = "2px solid white"
    window.open("/setup", "_self")
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
    fetch("/shutdown", {
        method: "POST",
    }).then(() => {
    }).catch(() => {
    });
}, false);