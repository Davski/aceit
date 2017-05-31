function getGameState() {
    var timer = document.getElementById("countDown");
    var requester = setInterval(function () {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/gamex", true);
        xhr.responseType = "json";
        xhr.onreadystatechange = function () {
            if (this.readyState === 4 && this.status === 200) {
                var game = (this.response);
                if (timer.innerHTML > 1 && game.TimeLeft < 15) {
                    timer.innerHTML = game.TimeLeft;
                } else if (game.TimeLeft > 15) {
                    timer.innerHTML = 0;
                } else {
                    getQuestionAndAlts(requester);
                    timer.innerHTML = 15;
                }
            }
        };
        xhr.send();
    }, 950);
}

function getQuestionAndAlts(intervalHandle) {
    if (!checkBrowserSupport()) {
        return false;
    }
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/gameServe", true);
    xhr.responseType = "json";
    xhr.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var obj = (this.response);
            if (obj === null) {
                clearInterval(intervalHandle);
                goTo("score", "1", "1"); // TODO: Replace the "1":s with courseID and moduleID.
            } else {
                $('input[name="answer"]').prop("disabled", false);
                $("#answerSubmit").prop("disabled", false);
                createQuestion(obj);
            }
        }
    };
    xhr.send();
}

//Get JSON object and put the Question to website by JQuery
function createQuestion(obj) {
    $('#quzzNo').text(obj.Question.Question);
    randomAlt(obj.Alts);
}

//Random the number 0,1,2 without repeating
function randomAlt(obj) {
    var altArr = [];
    for(var i in obj) {
        if(i !== "QuestionID"){
            altArr.push(obj[i]);
        }
    }
    
    var arr = [];
    while (arr.length < altArr.length) {
        var randomNumber = Math.floor((Math.random() * altArr.length));
        
        if (arr.indexOf(randomNumber) > -1) {
            continue;
        }
        
        arr[arr.length] = randomNumber;
    }
    
    createAlt(arr, altArr);
}

//Layout the alternatives that are related to the question by using Jquery
function createAlt(arr, altArr){
    $('#alternatives').empty();
    $('#showAnswer').empty();
    for(var i in arr){
        var alt = $("<input>").attr("type", "radio").attr("name", "answer")
        .attr("value", "alternative" + (arr[i] + 1)).attr("id", "alt" + (arr[i] + 1));
        var span = $("<span></span>").text(altArr[arr[i]]).attr("class", "spanText");
        $('#alternatives').append(alt).append(span).append($("<br>"));
    }
    
    var btn = $("<button></button>").text("Answer").attr("type","button")
    .attr("id","answerSubmit").attr("onclick","testAnswer()");
    $('#alternatives').append(btn);
    
}


/*
// This is a clientside countdown.
function timeCount() {
    var tzero = 14;
    var requester = setInterval(function () {
        $('#countDown').empty();
        $('#countDown').append(tzero);
        
        if (tzero === 0) {
            getQuestionAndAlts(requester);
            tzero = 15;
        } else {
            tzero--;
        }
    }, 1000);
}*/


// Compares the answer with the server and returns correct or not.
function testAnswer() {
    var answer = document.querySelector('input[name="answer"]:checked').id;
    var params = "answer=" + answer;
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST", "/submitAnswer", true);
    xmlhttp.responseType = "json";
    
    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    
    xmlhttp.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            var correct = this.response;
            if (correct) {
                $('#showAnswer').empty();
                $('#showAnswer').append("Correct");
                $('#showAnswer').css('color', 'green');
            } else {
                $('#showAnswer').empty();
                $('#showAnswer').append("Incorrect");
                $('#showAnswer').css('color', 'red');
            }
            if (answer !== null) $('input[name="answer"]').prop("disabled", true);
            if (answer !== null) $("#answerSubmit").prop("disabled", true);
        }
    };
    xmlhttp.send(params);   
}

// With this at the end of the file, the functions will be called when the DOM is ready.
$(document).ready(function () {
    getQuestionAndAlts();
    getGameState();
});

/*
$(document).ready(function() {
    var tzero = 14;
    var myVar = setInterval(function() {
        if (tzero === 0) console.log("timer is at:" + tzero);
        $('#countDown').empty();
        $('#countDown').append(tzero);
        if (tzero > 0) {
            tzero--;
        }
        else {
            clearInterval(myVar);
            testAnswer();
        }
    },1000);
});
*/