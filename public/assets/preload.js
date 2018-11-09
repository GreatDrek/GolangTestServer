document.body.style.overflow = 'hidden';

document.body.onload = function(){
	setTimeout(function(){
		
		var progressbar = $('[role="progressbar"]');
		progressbar.width = 30;
		
		var preloader = document.getElementById("page-preloader");
		if(!preloader.classList.contains("done")){
			preloader.classList.add("done");
			//preloader.classList.remove("preloader");
			document.body.style.overflow = 'auto';
		}
		
		var test = document.getElementById("myNavBar");
		test.classList.remove("unVisible");
	}, 0);
}





const red = document.querySelector('#loadTest');

function setProperty(duration) {
  red.style.setProperty('--animation-time', duration +'s');
}

function changeAnimationTime() {
  const animationDuration = Math.random();
  setProperty(animationDuration);
}

setInterval(changeAnimationTime, 100);




const red0 = document.querySelector('#loadTest');

function setProperty(duration) {
  red0.style.setProperty('--pos-top', duration +'%');
}

function changeAnimationTime() {
  const animationDuration = Math.random();
  setProperty(animationDuration);
}

setInterval(changeAnimationTime, 100);





$('.carousel0').carousel({
  interval: 10000
})