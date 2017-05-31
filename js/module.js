function createModuleList(courseID) {

    var xmlhttp = new XMLHttpRequest();
    var params = "courseID=" + courseID;
    xmlhttp.open("POST", "/getModules", true);
    xmlhttp.responseType = "json";

    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var moduleList = this.response;
            var moduleButtons = document.getElementById("moduleForm");
            var courseID = getUrlValue("courseID=");

            for (var i = 0; i < moduleList.length; i++) {
                createButton(moduleButtons, moduleList[i].Module, moduleList[i].ModuleID, courseID);
            }
        }
    };
    xmlhttp.send(params);
}

function createButton(body, name, moduleID, courseID) {
    var btn = document.createElement("button");
    var text = document.createTextNode(name);

    btn.setAttribute("type", "button");
    btn.setAttribute("id", moduleID);
    btn.setAttribute("onclick", "goTo('lobby', '" + courseID + "', '" + moduleID + "')");
    btn.appendChild(text);
    body.appendChild(btn);
}

// With this at the end of the file, the functions will be called when the DOM is ready.
$(document).ready(function () {
    var courseID = getUrlValue("courseID=");
    writeCourseName(courseID, "courseName");
    createModuleList(courseID);
});