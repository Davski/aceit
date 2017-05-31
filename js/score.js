function getLocalScore() {
    if (!checkBrowserSupport()) {
        return false;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/getLocalScore", true);
    xhr.responseType = "json";
    
    xhr.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var obj = (this.response);
            createList(obj);
        }
    };
    xhr.send();
}

function createList(obj) {
    //var nameList = [];
    //var myScore = [];
    //var score = [];

    /*for (var i = 0; i < obj.length; i++) {
        nameList[i] = $.map(obj[i], function (el) { return el }); // Convert object to array.
        score[i] = nameList[i][0];
    }*/

    //JUST MAKE UP TO BE PLAYER SCORE
    //myScore[0] = nameList[0][0];
    //myScore[1] = nameList[0][1];
    //console.log(obj.length);
    //console.log(obj[0].Name);
    var tbodyElement = document.getElementById("scoreList");
    //var playerTbodyElement = document.getElementById("scoreOfPlayer");

    //score = score.sort(function (a, b) { return b - a }); // Sort an array.

    /*for (var i = 0; i < score.length; i++) {
        for (var j = 0; j < nameList.length; j++) {
            if (score[i] === nameList[j][0]) {
                createTrAndTd(tbodyElement, nameList[j], i);
                if (nameList[j][1] === myScore[1]) {
                    createTrAndTd(playerTbodyElement, nameList[j], i);
                }
                // Function splice(deleted elem, antal deleting) used for deleting an specific element.
                nameList.splice(j, 1);
                break;
            }
        }
    }*/

    for(var i = 0; i < obj.length; i++) {
        createTrAndTd(tbodyElement, obj[i].name, obj[i].score, i);
    }

    //showTopFive();
    // getHighScore();
}

function createTrAndTd(tbody, name, score, i) {
    var trElement = document.createElement("tr");
    var tdElementRange = document.createElement("td");
    var tdElementInfo = document.createTextNode(i + 1);

    tdElementRange.appendChild(tdElementInfo);
    trElement.appendChild(tdElementRange);

    /*for (var x = 0; x < elem.length; x++) {
        var tdInfor = document.createTextNode(name);
        var tdElement = document.createElement("td");

        tdElement.appendChild(tdInfor);
        trElement.appendChild(tdElement);
    }*/
    
    createTd(name, trElement);
    createTd(score, trElement);

    tbody.appendChild(trElement);
}

function createTd(elem, trElement) {
    var tdInfor = document.createTextNode(elem);
    var tdElement = document.createElement("td");
    tdElement.appendChild(tdInfor);
    trElement.appendChild(tdElement);
}
/*
function getHighScore() {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', "/getCourses", true);
    xhr.responseType = 'json';

    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var obj = (this.response);
            showTopFive(obj);
        }
    };
    xhr.send();
}

function showTopFive() {
    var btn = document.getElementById("showScoreBtn");

    btn.onclick = function () {
        var tbodyLobby = document.getElementById("scoreList");
        var trLobby = tbodyLobby.getElementsByTagName("tr");

        var tbodyPlayer = document.getElementById("scoreOfPlayer");
        var trPlayer = tbodyPlayer.getElementsByTagName("tr");

        deleteRow(tbodyLobby, trLobby);
        deleteRow(tbodyPlayer, trPlayer);

        createHighScoreList();
    }
}

function deleteRow(tbodyElem, trElem) {
    var rowCount = trElem.length;
    for (var x = rowCount - 1; x >= 0; x--) {
        tbodyElem.removeChild(trElem[x]);
    }
}

function createHighScoreList() {
    // These values are made up. Have to find out how to get real value.
    //--------------------- HAVE TO THINK HOW TO GET DATA FROM DATABASE---------------------
    var nameList = Array(["Jonathan", 78], ["Johan", 70], ["Karl", 68], ["Malin", 68], ["Elin", 60]);
    var myScore = ["Erik", 47];
    var myScoreRange = 13;
    //--------------------------------------------------------------------------------------

    var tbodyElement = document.getElementById("scoreList");
    var playerTbodyElement = document.getElementById("scoreOfPlayer");

    for (var i = 0; i < nameList.length; i++) {
        createTrAndTd(tbodyElement, nameList[i], i);
    }

    createTrAndTd(playerTbodyElement, myScore, myScoreRange - 1);
}
*/

function goBackToMainMenu() {
    var btn = document.getElementById("exitBtn");
    btn.onclick = function () {
        exitTo("courses");
    };
}

// With this at the end of the file, the functions will be called when the DOM is ready.
$(document).ready(function () {
    getLocalScore();
    goBackToMainMenu();
});



//function LocalScore(obj) {
    /*var xhr = new XMLHttpRequest();
    xhr.open('GET', "/getCourses", true);
    //xhr.responseType = 'json';
    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            obj = JSON.parse(this.responseText);
            console.log(obj[0].Course + obj[0].CourseID);
            //document.getElementById("demo").innerHTML = response.name;
        }
    };
    xhr.send();
    */
        /*-------------------------------------------------------------------
        var getJSON = function(url, successHandler, errorHandler) {
            var xhr = new XMLHttpRequest();
            xhr.open('get', url, true);
            xhr.responseType = 'json';
            xhr.onload = function() {
                var status = xhr.status;
                if (status == 200) {
                    successHandler && successHandler(xhr.response);
                } else {
                    errorHandler && errorHandler(status);
                }
            };
            xhr.send();
        };

        getJSON('https://mathiasbynens.be/demo/ip', function(data) {
            alert('Your public IP address is: ' + data.ip);
        }, function(status) {
            alert('Something went wrong.');
        });
        ---------------------------------------------------------------------*/
//}

/*EXAMPLE---------------------DELETE LATER-------------------------
    aElement.appendChild(imgElement);
    var value = array.splice( index, 1 )[0];
    var a = Math.max(5, 10);

    var y = document.createElement("TR");
    y.setAttribute("id", "myTr");
    document.getElementById("myTable").appendChild(y);

    var z = document.createElement("TD");
    var t = document.createTextNode("cell");
    z.appendChild(t);
    document.getElementById("myTr").appendChild(z);
    --------------------------------------------------FINISH*/