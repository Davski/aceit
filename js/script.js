function checkBrowserSupport() {
    // Check if the browser supports these HTML DOM methods.
    return document.getElementById &&
        document.getElementsByTagName &&
        document.createElement &&
        document.createTextNode;
}

function getUrlValue(key) {
    var url = window.location.href;
    var keyIndex = url.indexOf(key) + key.length;
    var valueUrl = url.substring(keyIndex);
    var ampIndex = keyIndex + valueUrl.indexOf("&");
    if (ampIndex > keyIndex) {
        endIndex = ampIndex;
    } else {
        endIndex = url.length;
    }
    var value = url.substring(keyIndex, endIndex);
    return value;
}

function goTo(page, courseID, moduleID) {
    switch (page) {
        case "module":
            window.location.href = "module?courseID=" + courseID;
            break;
        case "lobby":
            window.location.href = "lobby?courseID=" + courseID + "&moduleID=" + moduleID;
            break;
        case "game":
            window.location.href = "game?courseID=" + courseID + "&moduleID=" + moduleID;
            break;
        case "score":
            window.location.href = "score?courseID=" + courseID + "&moduleID=" + moduleID;
            break;
        default:
            window.location.href = page + "?courseID=" + courseID + "&moduleID=" + moduleID;
            break;
    }
}

function exit() {
    console.log("exiting");
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/exit", true);
    
    xhr.send();
}

function logout() {
    console.log("logging out");
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/logout", true);
    xhr.send();
}

function exitTo(choice) {
    exit();

    switch (choice) {
        case "modules":
            console.log("exit to modules");
            var courseID = getUrlValue("courseID=");
            window.location.href = "module?courseID=" + courseID;
            break;
        case "courses":
            console.log("exit to courses")
            window.location.href = "/";
            break;
        case "logout":
            console.log("exit to logout")
            logout();
            window.location.href = "/";
            break;
        default:
            window.location.href = "#";
            break;
    }
}

function writeCourseName(courseID, courseNamePlaceID) {
    var xmlhttp = new XMLHttpRequest();
    var params = "courseID=" + courseID;
    xmlhttp.open("POST", "/getCourseName", true);
    xmlhttp.responseType = "json";

    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var courseNamePlace = document.getElementById(courseNamePlaceID);

            if (this.response === null) {
                courseName = "nonexisting course";
            } else {
                var courseName = this.response.Course;
            }

            courseNamePlace.innerHTML = courseName;
        }
    };
    xmlhttp.send(params);    
}

function writeModuleName(courseID, moduleID, moduleNamePlaceID) {
    var xmlhttp = new XMLHttpRequest();
    var params = "courseID=" + courseID + "&moduleID=" + moduleID;
    xmlhttp.open("POST", "/getModuleName", true);
    xmlhttp.responseType = "json";

    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var moduleNamePlace = document.getElementById(moduleNamePlaceID);

            if (this.response === null) {
                var moduleName = "nonexisting module";
            } else {
                moduleName = this.response.Module;
            }

            moduleNamePlace.innerHTML = moduleName;
        }
    };
    xmlhttp.send(params);
}

/*
 THIS IS A GO FUNCTION FOR DISPATCHER.GO

 the struct is this

 // Might need a better name

func submitAnswer(res http.ResponseWriter, req *http.Request) {
	// Replace with new way to get the active lobby, from here
	//
	//	session := checkCookie(res, req)
	//	lobbyID := session.Values["lobbyID"].(int)
	//	lobbyActive := master.lobbyFind(lobbyID, master.activeLobbys) //ska vara activeGames, inte lobbys
	//	// to here
	//	qID := lobbyActive.questionID
	

	var answer submittedAnswer
	answer.Answer = req.FormValue("answer")
	answer.Correct = true // checkAnswer(qID, answer)
	output, err := json.Marshal(answer)
	//
	//	if checkAnswer(qID, answer) {
	//		log.Println("rätt")
	//	} else {
	//		log.Println("fel")
		//}
	//
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)

}
*/