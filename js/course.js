function getUsername() {
    if (!checkBrowserSupport()) {
        return false;
    }

    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST", "/getUsername", true);
    xmlhttp.responseType = "json";

    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            document.getElementById("username").innerHTML = this.response.name;
        }
    };
    xmlhttp.send();
}

function createCourseList() {
    if (!checkBrowserSupport()) {
        return false;
    }

    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST", "/getCourses", true);
    xmlhttp.responseType = "json";

    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var courseList = (this.response);
            var courseButtons = document.getElementById("courseForm");

            for (var i = 0; i < courseList.length; i++) {
                createButton(courseButtons, courseList[i].Course, courseList[i].CourseID);
            }
        }
    };
    xmlhttp.send();
}

function createButton(body, name, courseID) {
    var btn = document.createElement("button");
    var text = document.createTextNode(name);

    btn.setAttribute("type", "button");
    btn.setAttribute("id", courseID);
    btn.setAttribute("onclick", "goTo('module', '" + courseID + "')");
    btn.appendChild(text);
    body.appendChild(btn);
}

function toggleExplainGame() {
    $('#explainGame').slideToggle(400);
}

// With this at the end of the file, the functions will be called when the DOM is ready.
$(document).ready(function () {
    getUsername();
    createCourseList();
});