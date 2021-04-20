const Keyboard = {
    elements: {
        main: null,
        keysContainer: null,
        keys: []
    },

    eventHandlers: {
        oninput: null,
        onclose: null
    },

    properties: {
        value: "",
        capsLock: false
    },

    init() {
        this.elements.main = document.createElement("div");
        this.elements.keysContainer = document.createElement("div");

        // Setup main elements
        this.elements.main.classList.add("keyboard", "keyboard--hidden");
        this.elements.keysContainer.classList.add("keyboard__keys");
        this.elements.keysContainer.appendChild(this._createKeys());

        this.elements.keys = this.elements.keysContainer.querySelectorAll(".keyboard__key");

        // Add to DOM
        this.elements.main.appendChild(this.elements.keysContainer);
        document.body.appendChild(this.elements.main);

        // Automatically use keyboard for elements with .use-keyboard-input
        document.querySelectorAll(".use-keyboard-input").forEach(element => {
            element.addEventListener("focus", () => {
                this.open(element.value, currentValue => {
                    element.value = currentValue;
                });
            });
        });
    },
    _createKeys() {
        const fragment = document.createDocumentFragment();
        const keyLayout = [
            "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace",
            "q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
            "a", "s", "d", "f", "g", "h", "j", "k", "l", "enter",
            "z", "x", "c", "v", "b", "n", "m", ".", ":", "/",
        ];
        keyLayout.forEach(key => {
            const keyElement = document.createElement("button");
            const insertLineBreak = ["backspace", "p", "enter", "/"].indexOf(key) !== -1;
            keyElement.setAttribute("type", "button");
            keyElement.classList.add("keyboard__key");

            switch (key) {
                case "backspace":
                    keyElement.classList.add("keyboard__key--wide");
                    keyElement.innerHTML = "⌫";
                    keyElement.addEventListener("touchstart", () => {
                        this.properties.value = this.properties.value.substring(0, this.properties.value.length - 1);
                        this._triggerEvent("oninput");
                    });
                    break;
                case "enter":
                    keyElement.classList.add("keyboard__key--wide");
                    keyElement.innerHTML = "↵";
                    keyElement.addEventListener("touchstart", () => {
                        if (sessionStorage.getItem("selection") === "password") {
                            let password = this.properties.value
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
                                            document.getElementById("ipaddress").disabled = false
                                            document.getElementById("gateway").disabled = false
                                            document.getElementById("mask").disabled = false

                                        } else {
                                            Keyboard.close()
                                        }
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
                        this._triggerEvent("oninput");
                    });
                    break;
                default:
                    keyElement.textContent = key.toLowerCase();
                    keyElement.addEventListener("touchstart", () => {
                        this.properties.value += this.properties.capsLock ? key.toUpperCase() : key.toLowerCase();
                        this._triggerEvent("oninput");
                    });
                    break;
            }

            fragment.appendChild(keyElement);

            if (insertLineBreak) {
                fragment.appendChild(document.createElement("br"));
            }
        });

        return fragment;
    },

    _triggerEvent(handlerName) {
        if (typeof this.eventHandlers[handlerName] == "function") {
            this.eventHandlers[handlerName](this.properties.value);
        }
        let elem = document.getElementById('server');
        elem.scrollLeft = elem.scrollWidth;
    },

    open(initialValue, oninput, onclose) {
        this.properties.value = initialValue || "";
        this.eventHandlers.oninput = oninput;
        this.eventHandlers.onclose = onclose;
        this.elements.main.classList.remove("keyboard--hidden");
    },

    close() {
        this.properties.value = "";
        this.eventHandlers.oninput = oninput;
        this.eventHandlers.onclose = onclose;
        this.elements.main.classList.add("keyboard--hidden");
    }
};

window.addEventListener("DOMContentLoaded", function () {
    Keyboard.init();
});