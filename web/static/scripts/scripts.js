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
    plane.setAttribute("material", "opacity: 0.0; transparent: true")
    let innerPlaneU = document.createElement("a-plane")
    innerPlaneU.setAttribute("position",  "0 " + ((antennaData.I_offset) / 2 + 1.6).toString() + " -30")
    innerPlaneU.setAttribute("color", "gray")
    innerPlaneU.setAttribute("height", offset_i)
    innerPlaneU.setAttribute("width", antennaData.Size_J - 2)
    let innerPlaneD = document.createElement("a-plane")
    innerPlaneD.setAttribute("position",  "0 " + ((-antennaData.I_offset) / 2 + 1.6).toString() + " -30")
    innerPlaneD.setAttribute("color", "gray")
    innerPlaneD.setAttribute("height", offset_i)
    innerPlaneD.setAttribute("width", antennaData.Size_J - 2)
    let innerPlaneR = document.createElement("a-plane")
    innerPlaneR.setAttribute("position",  ((-antennaData.J_offset - antennaData.SlotSize + 1) / 2).toString() + " 1.6 -30")
    innerPlaneR.setAttribute("color", "gray")
    innerPlaneR.setAttribute("height", 1)
    innerPlaneR.setAttribute("width", antennaData.J_offset - 1)
    let innerPlaneL = document.createElement("a-plane")
    innerPlaneL.setAttribute("position",  ((antennaData.J_offset + antennaData.SlotSize - 1) / 2).toString() + " 1.6 -30")
    innerPlaneL.setAttribute("color", "gray")
    innerPlaneL.setAttribute("height", 1)
    innerPlaneL.setAttribute("width", antennaData.J_offset - 1)
    plane.appendChild(innerPlaneU)
    plane.appendChild(innerPlaneD)
    plane.appendChild(innerPlaneR)
    plane.appendChild(innerPlaneL)
    for (i = 0; i < antennaData.Size_I - 2; i++) {
        for (j = 0; j < antennaData.Size_J - 2; j++) {
            if (j - offset_j < -(antennaData.SlotSize - 1) / 2 || j - offset_j > (antennaData.SlotSize - 1) / 2 || i - offset_i !== 0) {
                let vec = document.createElement("a-entity")
                vec.setAttribute("position",  (j - offset_j).toString() + " " + (i - offset_i + 1.6).toString() + " -29.9")
                if (m === 0) {
                    vec.setAttribute("line", "start: 0 0.4 0; end: 0 -0.4 0; color: red")
                    vec.setAttribute("rotation", "0 0 " + (90 + antennaData.MD[antennaData.Size_I - 3 - i][j]).toString())
                } else {
                    vec.setAttribute("line", "start: 0 0.4 0; end: 0 -0.4 0; color: blue")
                    vec.setAttribute("rotation", "0 0 " + (90 + antennaData.ED[antennaData.Size_I - 3 - i][j]).toString())
                }
                plane.appendChild(vec)
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