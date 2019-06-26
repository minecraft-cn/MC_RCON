var addr = document.getElementById("addr");
var pswd = document.getElementById("pswd");
var cmd = document.getElementById("cmd");
var resp = document.getElementById("resp");

function runCmd(host, addr, pswd, cmd) {
    console.log(host, addr, pswd, cmd);
    let xhr = new XMLHttpRequest();

    xhr.addEventListener("readystatechange", function () {
        // if (this.readyState === 4) {
        console.log(this.responseText);
        resp.value = this.responseText;
        // }
    });

    xhr.open("GET", `http://${host}/cmd?addr=${addr}&pswd=${pswd}&cmd=${cmd}`);
    xhr.send(null);
}

function RunCmd() {

    // runCmd(window.location.host, addr.value, pswd.value, cmd.value);
    runCmd("127.0.0.1:25576", addr.value, pswd.value, cmd.value);
}

