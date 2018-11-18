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