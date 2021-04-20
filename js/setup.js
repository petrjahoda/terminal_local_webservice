const dhcpSlider = document.getElementById("dhcp-slider")
const leftButton = document.getElementById("left-button")
const middleButton = document.getElementById("middle-button")
const rightButton = document.getElementById("right-button")
const passwordField = document.getElementById("password")
const ipaddress = document.getElementById("ipaddress")
const gateway = document.getElementById("gateway")
const server = document.getElementById("server")
const mask = document.getElementById("mask")

dhcpSlider.addEventListener('change', function (e) {
    if (dhcpSlider.checked) {
        document.getElementById("ipaddress").disabled = true
        document.getElementById("gateway").disabled = true
        document.getElementById("mask").disabled = true
        Keyboard.close()
    } else {
        document.getElementById("ipaddress").disabled = false
        document.getElementById("gateway").disabled = false
        document.getElementById("mask").disabled = false
    }
}, false);

leftButton.addEventListener('click', function () {
    leftButton.style.border = "2px solid red"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid white"
    window.open("/", "_self")
}, false);

leftButton.addEventListener('touchstart', function () {
    leftButton.style.border = "2px solid red"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid white"
    window.open("/", "_self")
}, false);


middleButton.addEventListener('touchend', function () {
    middleButton.style.border = "2px solid white"
    middleButton.blur()
    setTimeout(() => {
        middleButton.blur()
    }, 10);
})

rightButton.addEventListener('touchend', function () {
    rightButton.style.border = "2px solid white"
    rightButton.blur()
    setTimeout(() => {
        rightButton.blur()
    }, 10);
})

middleButton.addEventListener('touchstart', function () {
    middleButton.blur()
    if (!dhcpSlider.checked) {
        leftButton.style.border = "2px solid white"
        middleButton.style.border = "2px solid red"
        rightButton.style.border = "2px solid white"
        ipaddress.value = ""
        mask.value = ""
        gateway.value = ""
        server.value = ""
    } else {
        leftButton.style.border = "2px solid white"
        middleButton.style.border = "2px solid red"
        rightButton.style.border = "2px solid white"
        server.value = ""
    }
    middleButton.blur()
}, false);

rightButton.addEventListener('touchstart', function () {
    leftButton.style.border = "2px solid white"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid red"
    rightButton.blur()
}, false);

passwordField.addEventListener('touchstart', function () {
    sessionStorage.setItem("selection", "password")
}, false);

ipaddress.addEventListener('touchstart', function () {
    sessionStorage.setItem("selection", "ipaddress")
}, false);

server.addEventListener('touchstart', function () {
    sessionStorage.setItem("selection", "server")
}, false);

gateway.addEventListener('touchstart', function () {
    sessionStorage.setItem("selection", "gateway")
}, false);

mask.addEventListener('touchstart', function () {
    sessionStorage.setItem("selection", "mask")
}, false);