let m = 0
let k = 0
let antennaData
let num = 0

let loadAntenna = function () {
    let offset_i = (antennaData.Size_I - 2 - 1) / 2
    let offset_j = (antennaData.Size_J - 2 - 1) / 2
    let scene = document.getElementById("scene")
    if (k !== 0) {
        scene.removeChild(scene.lastChild)
    }
    k = 1
    let plane = document.createElement("a-plane")
    plane.setAttribute("id", "plane")
    plane.setAttribute("position", "0 1.6 -30")
    plane.setAttribute("rotation", "0 0 0")
    plane.setAttribute("height", (antennaData.Size_I - 2).toString())
    plane.setAttribute("width", (antennaData.Size_J - 2).toString())
    plane.setAttribute("material", "opacity: 0.0; transparent: true")
    for (i = 0; i < antennaData.Size_I - 2; i++) {
        for (j = 0; j < antennaData.Size_J - 2; j++) {
            if (j - offset_j < -(antennaData.SlotSize - 1) / 2 || j - offset_j > (antennaData.SlotSize - 1) / 2 || i - offset_i !== 0) {
                let innerPlane = document.createElement("a-plane")
                innerPlane.setAttribute("position",  (j - offset_j).toString() + " " + (i - offset_i + 1.6).toString() + " -30")
                innerPlane.setAttribute("rotation", "0 0 0")
                innerPlane.setAttribute("color", "gray")
                innerPlane.setAttribute("height", 1)
                innerPlane.setAttribute("width", 1)
                let vec = document.createElement("a-plane")
                vec.setAttribute("position",  "0 0 0.1")
                vec.setAttribute("height", 0.8)
                vec.setAttribute("width", 0.1)
                if (m === 0) {
                    vec.setAttribute("color", "red")
                    vec.setAttribute("rotation", "0 0 " + (90 + antennaData.MD[antennaData.Size_I - 3 - i][j]).toString())
                } else {
                    vec.setAttribute("color", "blue")
                    vec.setAttribute("rotation", "0 0 " + (90 + antennaData.ED[antennaData.Size_I - 3 - i][j]).toString())
                }
                innerPlane.appendChild(vec)
                plane.appendChild(innerPlane)
            }
        }
    }
    scene.appendChild(plane)
}

let nextAntenna = function () {
    num += 1
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/get_antenna/' + num.toString(), true);
    xhr.onreadystatechange = function() {
        if (this.status === 200) {
            antennaData = JSON.parse(this.responseText)
            loadAntenna()
        }
    }
    xhr.send()
}

let prevAntenna = function () {
    num -= 1
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/get_antenna/' + num.toString(), true);
    xhr.onreadystatechange = function() {
        if (this.status === 200) {
            antennaData = JSON.parse(this.responseText)
            loadAntenna()
        }
    }
    xhr.send()
}

let switchMode = function () {
    m = (m + 1) % 2
    loadAntenna()
}

document.onkeydown = function(e){
    e = e || window.event;
    var key = e.which || e.keyCode;
    if (key === 81) { //q
        prevAntenna()
    } else if (key === 67) { //c
        switchMode()
    } else if (key === 69) { //e
        nextAntenna()
    }
}