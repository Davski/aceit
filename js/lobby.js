function getLobbyState(courseID, moduleID) {
    
    var requester = setInterval(function () {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/lobbyx", true);
        xhr.responseType = "json";
        
        xhr.onreadystatechange = function () {
            if (this.readyState === 4 && this.status === 200) {
                var lobby = (this.response);
                if (lobby === null) { // Ugly ugly fix, lobby==null is probably a bad check(but it works)
                    clearInterval(requester);
                    goTo("game", courseID, moduleID);
                    return
                }
                checkLobbyState(lobby);
            }
        };
        xhr.send();
    }, 950);
}

function checkLobbyState(lobby) {
    addPlayers(lobby.Players, lobby.Course.Course, lobby.Module.Module);
    
    var timer = document.getElementById("countDown");
    if ((timer.innerHTML) > 1) {
        timer.innerHTML = lobby.TimeLeft;
    }
    else {
        timer.innerHTML = 0; // Dirty hack because of chaos with TimeLeft near 0 in back-end.
    }
}

function addPlayers(players, course, module) {
    $('#moduleName').empty();
    $('#lobbyList').empty();
    $('#moduleName').text(course + " : " + module);
    for (var i = 0; i < players.length; i++) {
        var playerName = $("<li></li>").text(players[i].name).attr("id", "user" + (i + 1));
        $('#lobbyList').append(playerName);
    }
}

// With this at the end of the file, the functions will be called when the DOM is ready.
$(document).ready(function () {
    var courseID = getUrlValue("courseID=");
    var moduleID = getUrlValue("moduleID=");
    
    //writeCourseName(courseID, "courseName"); <---- Onöjdig eftersom vi redan har fått course name och module name i lobby uppe
    //writeModuleName(courseID, moduleID, "moduleName");
    getLobbyState(courseID, moduleID);
});