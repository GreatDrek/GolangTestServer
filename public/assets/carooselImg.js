$(document).ready(function() {
	$.ajax({
				url: 'caroosel',
				type: 'POST',
				success:function(results) {
					var mass = JSON.parse(results);
					var lenghtMass = mass.length;
					for(var i = 0; i < lenghtMass; i++){
						addImgCaroosel(mass[i], i);
					}
				}
			});
});

function addImgCaroosel(imgName, first){
	var divCarooselItem = document.createElement('div');
	if(first === 0){ 
		divCarooselItem.className = "carousel-item active";
	}else{
		divCarooselItem.className = "carousel-item";
	}
	
	var imgCaroosel = document.createElement('img');
	imgCaroosel.src = "./assets/img/photo_room/caroosel/" + imgName;
	imgCaroosel.className = "d-block w-100";
	
	divCarooselItem.appendChild(imgCaroosel);
	
	var parentElem = document.getElementById('caroosel_add');
	
	parentElem.appendChild(divCarooselItem);
}