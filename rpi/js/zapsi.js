const source = new EventSource('/listen');
source.addEventListener('time', () => {
    document.getElementById("time").innerHTML = event.data;

}, false);

const networkDataSource = new EventSource('/networkdata');
networkDataSource.addEventListener('networkdata', () => {
    const networkdata = event.data.split(";");
    document.getElementById("ipaddress").innerHTML = networkdata[0];
    document.getElementById("mask").innerHTML = networkdata[1];
    document.getElementById("gateway").innerHTML = networkdata[2];
    document.getElementById("dhcp").innerHTML = networkdata[3];
    document.querySelector('meta[name="timer"]').setAttribute("content", networkdata[4] + ";URL='" + networkdata[5] + "'");
    document.getElementById("url").innerHTML = networkdata[6];
    document.getElementById("remainingTime").innerHTML = networkdata[4];
}, false);