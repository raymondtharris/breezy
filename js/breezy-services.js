var breezyService = angular.module('breezyService', [])

breezyService.service('$dateformat', function($rootScope){
	this.formatDate = function(dateString){
		var dateHalf = dateString.split("T")
		var dateParts = dateHalf[0].split("-")
		//console.log(dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0])
		return dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0]
	}
});

breezyService.service('$datastorage', function($rootScope, $http){
	//Service that interacts with the web backened to get and post data to the server
	this.GetAllPosts = function(){//Function returns all posts stored in the database
		$http.get("/get_all_posts").success(function(data){
			//console.log(data)
			return data
		})

	}
	this.GetAllMedia = function(mediaFlag){//This function gets all media in the database and can be culled by type
		$http.get("/getallmedia/"+mediaFlag).success(function(data){
			return data
		})
	}
	this.GetPostsBetween = function(lower,upper){
		$http.get("/get_posts/between/"+lower+"-"+upper).success(function(data){
			return data
		})	
	}
	this.GetPosts = function(option, filter){
		switch(option){
			case "created":
				$http.get("/get_posts/created/"+filter).success(function(data){
					return data
				});
				break
			case "creator":
				$http.get("/get_posts/creator/"+filter).success(function(data){
					return data
				});
				break
			case "single":
				$http.get("/get_posts/single/"+filter).success(function(data){
					return data
				});
				break
		}
	}
	this.GetMedia = function(option, filter){
		switch(option){
			case "single":
				$http.get("/get_media/single/"+filter).success(function(data){
					return data
				});
				break
			case "multi":
				$http.post("/get_media/multi/", filter).success(function(data){
					return data
				});
				break
			case "added":
				$http.get("/get_media/added/"+filter).success(function(data){
					return data
				});
				break

		}
	}
});

breezyService.service('$files', function($rootScope,$http){
	var currentProgress;
	var uploadSize;
	var currentFileSize;
	this.upload = function(fileList){
		console.log(fileList)
		currentProgress= 0;
		uploadSize=0;
		for(var i = 0; i < fileList.length; i++){
			uploadSize +=fileList[i].size;
		}
		//console.log(uploadSize);
		for(var i = 0; i < fileList.length; i++){
			var xmlHttpReq = new XMLHttpRequest();
			currentFileSize = fileList[i].size;
			var formdata = new FormData();
			formdata.append("file", fileList[i]);
			formdata.append("filetype", fileList[i].type);
			xmlHttpReq.open("POST", "/uploadfile");
			xmlHttpReq.upload.addEventListener("progress", this.uploadProgress, false);
			xmlHttpReq.addEventListener("load", this.uploadComplete(xmlHttpReq, fileList[i]), false);
			xmlHttpReq.send(formdata);
			/*
			xmlHttpReq.setRequestHeader('X_FILE_NAME', fileList[i].name);
			xmlHttpReq.setRequestHeader('X_FILE_SIZE', fileList[i].size);
			xmlHttpReq.setRequestHeader('X-Requested-With', true);
			*/
	//		xmlHttpReq.setRequestHeader('Content-Type', fileList[i].type);
		//	xmlHttpReq.setRequestHeader('Content-Type', 'multipart/form-data');
			console.log(fileList[i].type);
		/*	xmlHttpReq.upload.addEventListener("progress", this.uploadProgress, false);
			xmlHttpReq.upload.addEventListener("loadstart", this.uploadStart, false);
			//xmlHttpReq.upload.addEventListener("loadend", this.uploadEnd(xmlHttpReq.upload, fileList[i]), false);
			xmlHttpReq.addEventListener("load", this.uploadComplete(xmlHttpReq, fileList[i]), false);
			//xmlHttpReq.addEventListener("error", this.uploadFailed, false);
			//xhr.addEventListener("abort", uploadCanceled, false)
		*/
		//	xmlHttpReq.send(fileList[i]);
			if(i == fileList.length-1){
				this.allFiles(fileList);
			}
		}
	}
	var getUploadSize = function(){
		return uploadSize;
	}
	this.uploadStart = function(){
	//	$rootScope.$broadcast('started');
	}
	this.uploadProgress = function(evt){

		currentProgress += Math.round(evt.loaded * 100 / uploadSize);

	//	$rootScope.$broadcast('uploadProgress', currentProgress);
	}
	this.uploadComplete = function(req, file){
		console.log(req);
		$rootScope.$broadcast('uploadComplete',file);
	}
	this.uploadFailed = function(){

	}
	this.uploadEnd = function(req, file){
		console.log('fi');
		//$rootScope.$broadcast('uploadComplete', file);
	 }
	this.allFiles = function(fileList){
		$rootScope.$broadcast('uploadCompleted',fileList);
	}

});
