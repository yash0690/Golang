document.getElementById('btn').onclick = function() {

	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		
		if (xhttp.readyState == 4 && xhttp.status == 200) {
			var message;
			if (xhttp.responseText.includes('true')) {
				message = 'User already exists';
			} else {
				message = 'New user is registered';
			}
			
			document.getElementById("errorMessage").innerHTML = message;
		}
	};
	xhttp.open("POST", "isUser", true);
	var data = new FormData();

	data.append('new-word', document.getElementById("entry").value);
	xhttp.send(data);
}
