const leftButton = document.getElementById("left-button")
const middleButton = document.getElementById("middle-button")
const rightButton = document.getElementById("right-button")
//
// leftButton.addEventListener('touchstart', function() {
//     document.getElementById("info").innerText = "touched left"
// }, false);
//
//
// middleButton.addEventListener('touchstart', function() {
//     document.getElementById("info").innerText = "touched middle"
// }, false);
//
// rightButton.addEventListener('touchstart', function() {
//     document.getElementById("info").innerText = "touched right"
// }, false);

leftButton.addEventListener('click', function() {
    document.getElementById("info").innerText = "clicked left"
    fetch("/restart", {
        method: "POST",
    }).then(() => {
    }).catch(() => {
    });
}, false);

middleButton.addEventListener('click', function() {
    document.getElementById("info").innerText = "clicked middle"
}, false);

rightButton.addEventListener('click', function() {
    document.getElementById("info").innerText = "clicked right"
    fetch("/shutdown", {
        method: "POST",
    }).then(() => {
    }).catch(() => {
    });
}, false);