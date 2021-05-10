const dhcpSlider = document.getElementById("dhcp-slider")
const leftButton = document.getElementById("left-button")
const middleButton = document.getElementById("middle-button")
const rightButton = document.getElementById("right-button")
const passwordField = document.getElementById("password")
const ipaddress = document.getElementById("ipaddress")
const gateway = document.getElementById("gateway")
const server = document.getElementById("server")
const mask = document.getElementById("mask")
var position = 0


dhcpSlider.addEventListener('change', function (event) {
    if (dhcpSlider.checked) {
        document.getElementById("ipaddress").disabled = true
        document.getElementById("gateway").disabled = true
        document.getElementById("mask").disabled = true
    } else {
        document.getElementById("ipaddress").disabled = false
        document.getElementById("gateway").disabled = false
        document.getElementById("mask").disabled = false
    }
}, false);


leftButton.addEventListener('click', function (event) {
    callRpiIndex();
}, false);

function callRpiIndex() {
    leftButton.style.border = "2px solid red"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid white"
    window.open("/", "_self")
}



middleButton.addEventListener('click', function (event) {
    callRpiResetAll();
}, false);


rightButton.addEventListener('click', function (event) {
    callRpiSaveNetworkChange();
}, false);

function callRpiResetAll() {
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
}

function checkInputData() {
    let result = false
    let ipResult = false
    let maskResult = false
    let gatewayResult = false
    if (/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(ipaddress.value)) {
        ipResult = true
        ipaddress.style.border = "1px solid white"
    } else {
        ipaddress.style.border = "1px solid red"
    }
    if (/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(mask.value)) {
        maskResult = true
        mask.style.border = "1px solid white"
    } else {
        mask.style.border = "1px solid red"
    }
    if (/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(gateway.value)) {
        gatewayResult = true
        mask.style.border = "1px solid white"
    } else {
        gateway.style.border = "1px solid red"
    }
    if (ipResult && maskResult && gatewayResult) {
        ipaddress.style.border = "1px solid white"
        gateway.style.border = "1px solid white"
        mask.style.border = "1px solid white"
        result = true
    }
    return result;
}

function callRpiSaveNetworkChange() {
    leftButton.style.border = "2px solid white"
    middleButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid red"
    rightButton.blur()
    if (dhcpSlider.checked) {
        document.getElementById("ipaddress").disabled = true
        document.getElementById("gateway").disabled = true
        document.getElementById("mask").disabled = true
        let data = {
            password: "3600",
            server: server.value,
        };
        fetch("/dhcp", {
            method: "POST",
            body: JSON.stringify(data)
        }).then(() => {
            window.open("/", "_self")
        }).catch(() => {
        });
    } else {
        let resultOk = checkInputData();
        if (resultOk) {
            let data = {
                password: "3600",
                ipaddress: ipaddress.value,
                mask: mask.value,
                gateway: gateway.value,
                server: server.value,
            };
            fetch("/static", {
                method: "POST",
                body: JSON.stringify(data)
            }).then(() => {
                window.open("/", "_self")
            }).catch(() => {
            });
        }
    }
}

passwordField.addEventListener('click', function (event) {
    sessionStorage.setItem("selection", "password")
}, false);

ipaddress.addEventListener('click', function (event) {
    sessionStorage.setItem("selection", "ipaddress")
}, false);

server.addEventListener('click', function (event) {
    sessionStorage.setItem("selection", "server")
}, false);

gateway.addEventListener('click', function (event) {
    sessionStorage.setItem("selection", "gateway")
}, false);

mask.addEventListener('click', function (event) {
    sessionStorage.setItem("selection", "mask")
}, false);

passwordField.addEventListener("keyup", function(event) {
    if (event.key === "Enter") {
        let password = passwordField.value
        let data = {
            password: password
        };
        fetch("/password", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function (data) {
                let result = JSON.parse(data);
                if (result["Result"] === "ok") {
                    document.getElementById("password").hidden = true
                    if (!dhcpSlider.checked) {
                        console.log("disabling elements")
                        document.getElementById("ipaddress").disabled = false
                        document.getElementById("gateway").disabled = false
                        document.getElementById("mask").disabled = false

                    }
                    console.log("disabling default")
                    middleButton.disabled = false
                    middleButton.style.pointerEvents = "auto"
                    rightButton.disabled = false
                    rightButton.style.pointerEvents = "auto"
                    document.getElementById("server").disabled = false
                    document.getElementById("dhcp-slider").disabled = false
                }
            })
        }).catch(() => {
        });
    }
});