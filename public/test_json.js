var services = {
	name: "test",
	type: "one"
};

var out = document.getElementById('out');
for (var key in services){
	out.innerHTML += key + "  " + services[key] + ' ';
}