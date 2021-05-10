let timeLeft = 15;
const downloadTimer = setInterval(function () {
    let serverActive = document.getElementById("server-active-panel")
    if (serverActive.innerText.includes("server dostupný")) {
        document.getElementById("server-info").innerText = "stránka serveru se načte za " + timeLeft + " vteřin";
        if (timeLeft <= 0) {
            clearInterval(downloadTimer);
            fetch("/stop_stream", {
                method: "POST",
            }).then(() => {
                window.open(document.getElementById("server").innerHTML, "_self")
            }).catch((error) => {
                console.log(error)
            }).catch((error) => {
                console.log(error)
            });
        }
        timeLeft -= 1;
    } else {
        timeLeft = 15
    }
}, 1000);

const networkDataSource = new EventSource('/networkdata');
networkDataSource.addEventListener('networkdata', (e) => {
    const networkData = e.data.split(";");
    document.getElementById("ipaddress").innerHTML = networkData[0];
    document.getElementById("mask").innerHTML = networkData[1];
    document.getElementById("gateway").innerHTML = networkData[2];
    document.getElementById("dhcp").innerHTML = networkData[3];
    document.getElementById("server").innerHTML = networkData[4];
    document.getElementById("active-panel").innerText = networkData[6];
    document.getElementById("server-active-panel").innerText = networkData[7];
    document.getElementById("active-panel").style.color = networkData[8];
    document.getElementById("server-active-panel").style.color = networkData[9];
    document.getElementById("mac-panel").innerText = networkData[5];
}, false);

const middleButton = document.getElementById("middle-button")


middleButton.addEventListener('touchstart', function (event) {
    callRpiSetup();
}, false);


function callRpiSetup() {
    middleButton.style.border = "2px solid red"
    window.open("/setup", "_self")
}