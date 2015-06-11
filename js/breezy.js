var breezy = angular.module('breezyApp',['ngSanitize'])

breezy.directive('contenteditable', ['$sce', function($sce) {
  return {
    restrict: 'A', // only activate on element attribute
    require: '?ngModel', // get a hold of NgModelController
    link: function(scope, element, attrs, ngModel) {
      if (!ngModel) return; // do nothing if no ng-model

      // Specify how UI should be updated
      ngModel.$render = function() {
        element.html($sce.getTrustedHtml(ngModel.$viewValue || ''));
      };

      // Listen for change events to enable binding
      element.on('blur keyup change', function() {
        scope.$evalAsync(read);
		if(!scope.contentDirty){
			scope.contentDirty = true
		}
      });
      read(); // initialize

      // Write data to the model
      function read() {
        var html = element.html();
        // When we clear the content editable the browser leaves a <br> behind
        // If strip-br attribute is provided then we strip this out
        if ( attrs.stripBr && html == '<br>' ) {
          html = '';
        }
		
        ngModel.$setViewValue(html);
      }
    }
  };
}]);

breezy.directive('infiniteScroll', function() {
	return function(scope, elemen, attr){
		console.log(elemen)
		var container = elemen[0]
	  	elemen.bind("scroll", function() {
			console.log("scrolling")
		})	
	}	
})

breezy.controller("BreezyController", function($scope,$http){
	//Controller for the blog portion of Breezy
	$scope.Title = ""
	$scope.DisplayPostCreator = false
	$scope.postlist =[]
	$scope.searchEnabled= true
	$scope.searchInput = ""
	$scope.lower = 0
	$scope.higher = 4
	$http.get("/get_blog_display").success(function(data){
		$scope.Title = data.Title
		if (data.UserCount > 1 ){
			$scope.DisplayPostCreator = true
		}
	})
	$http.get("/get_posts/between/"+$scope.lower+"-"+$scope.higher).success(function(data){
		$scope.postlist = data
	})
	$scope.loadMorePosts = function() {
		$scope.lower=$scope.higher
		$scope.higher=$scope.higher+2
		$http.get("/get_posts/between/"+$scope.lower+"-"+$scope.higher).success(function(data){
			for (var i = 0; i < data.length; i++ ){
				$scope.postlist.push(data[i])
			}
		})

	}
	/*
 	$http.get("/get_all_posts").success(function(data){
		//console.log(data)
		$scope.postlist = data
	})
	*/
	$scope.formatDate = function(dateString){
		var dateHalf = dateString.split("T")
		var dateParts = dateHalf[0].split("-")
		//console.log(dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0])
		return dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0]
	}
	
	// Add in search functionality of the blog.
	$scope.submitSearch = function(data){
		//console.log(data)
		$scope.searchInput = data.searchInput
		console.log($scope.searchInput)
		$http.post("/getsearch/", {searchtext: $scope.searchInput}).success(function(data){
			
		})
	}
	$scope.postsByDate = function(option) {
		console.log(option.post.Created)
		$http.get("/get_posts/created/"+option.post.Created).success(function(data){
			console.log(data)
		})
	}
	$scope.postsByCreator = function(option) {
		console.log(option.post.Creator)
		$http.get("/get_posts/creator/"+option.post.Creator).success(function(data){
			console.log(data)
		})
	}

})

breezy.controller("BreezyNavigationController", function($scope, $http, $window){
	$scope.openBlog = function(){
		//Open the users blog in a separate page/tab
	}
	$scope.gotoPage = function(option){
		//Sends user to the page suggested by the option
		console.log(option)
		switch (option){
			case "settings":
				$window.location.href="/settings"
				break;
			case "media":
				$window.location.href="/medialist"
				break;
			case "posts":
				$window.location.href="/postlist"
				break;
			case "new post":
				$window.location.href="/edit"
				break;
			default:
				console.log("No locaiton to go to.")
				break;
		}
	}
})

breezy.controller("BreezyDashboardController", function($scope, $http, $window){
	
})

breezy.controller("BreezyMediaLibraryController", function($scope, $http, $window, $element){
	$scope.MediaList = [];
	$http.get("/get_all_media").success(function(data){
		$scope.MediaList = data;
	})
	$scope.removeMedia = function(){
		//console.log(this.$index)
		var indDelete = this.$index
		$http.get("/deletemedia/"+this.item.ID).success(function(data){
			//console.log("deleted " + indDelete)
			$scope.MediaList.splice(indDelete, 1)
		});
	}
	$scope.getIndex = function(){
		console.log(this.$index)
	}
	$scope.compareDimensions = function(el){
		//console.log(this);
		//console.log(angular.element($element[0].querySelector('#thumbnail'+this.$index))[0].offsetHeight)	
		var elm = angular.element($element[0].querySelector('#thumbnail'+this.$index))[0]
		if ((elm.offsetHeight && elm.offsetWidth ) < 160){
			//console.log("small")
			return false
		}
		if (elm.offsetHeight > elm.offsetWidth ){
			//console.log("taller")
			return true
		} else {
			//console.log("wider or equal")
			return true
		}
		return false

	}
	$scope.displayMedia =  function(){
		console.log(this.item.ID)
	}
})

breezy.controller("BreezySetupController",function($scope, $http, $window){
	$scope.setupConfig={"username":"","password":"","name":"","blogname":""}
	$scope.submitSetupConfig = function(){
		//send password to be setup to sent to server
		$http.post("/setup_config", $scope.setupConfig).success(function(data){
			console.log(data)
				//send user to dashboard
				$window.location.href="/admin"
		})
	}
})

breezy.controller("BreezyLoginController", function($scope, $http, $window){
	$scope.loginCredentials = {"username":"", "password":""}
	$scope.incorrect = false
	var loginSent = {"username":"", "password":""}	
	$scope.submitLoginInfo = function(){
		//submitLoginInfo function sends loginCredentials to see if there is a match
		//and if it works will transfer user to dashboard. 
		$scope.incorrect = false
		$http.post("/checkcredentials", $scope.loginCredentials).success(function(data){
			if(data == "true"){
				//Correct loginCredentials go to dashboard
				$window.location.href="/dashboard"
			}else{
				console.log("Send Error to user about loginCredentials")
				$scope.incorrect = true
			}			
			// if all good go to dashboard
		})
		
	}
	$scope.incorrectCredentials = function(){

	}
})

breezy.controller("BreezyPostsController", function($scope, $http, $dateformat){
	$scope.postlist =[]
	$http.get("/get_all_posts").success(function(data){
		console.log(data)
		$scope.postlist = data
	})
	$scope.deletePost = function(){
		console.log(this.post.ID)
		var postDelete = this.$index;
		$http.get("/deletepost/"+this.post.ID).success(function(data){
			$scope.postlist.splice(postDelete, 1)
		})
	}
	$scope.formatDate = function(dateString){
		return $dateformat.formatDate(dateString);
	}
})


breezy.controller("BreezyEditorController", function($scope, $http, $timeout, $window){
	$scope.preview=false
	$scope.contentDirty=true
	$scope.postData={title:"", dateCreated:"", dateModified:""}
	$scope.mediaData={"Links":[], "Images":[], "Audio":[], "Video":[]}
	$scope.recentlyAdded = [];
	$scope.MediaList = [];
	$scope.hasRecently = true;
	$scope.preview=true;
	$scope.showPreview = false;
	$http.get("/get_all_media").success( function(data){
		$scope.MediaList = data
	})
	$scope.togglePreview = function(newValue){
		$scope.preview=newValue
		console.log($scope.markdownContent)
		
		if($scope.contentDirty){
			console.log("send to server for translation to markup.")
			var temp = String($scope.markdownContent).replace(/<[^>]+>/gm, '\n');
			var contentToConvert = {"markdownContent": temp, "markupContent":"", "postData": $scope.postData}
			console.log(contentToConvert)
			$http.post("/mdowntomup", contentToConvert).success(function(data){
				$scope.showPreview = true;
				console.log(data)
				console.log($scope.MarkupContent)
				$scope.postData.title = data.PostData.Title
				$scope.mediaData = data.MediaData
				$scope.markupContent =  data.MarkupContent
				$scope.contentDirty=false	
			});
			
		}
	}
	$scope.toggleEditor = function(){
		$scope.preview = true;
		$scope.showPreview = false;
	}
	$scope.savePost = function(){
		console.log($scope)
		//if contentDirty send markdown to change
		//on return send data to be saved in database 			
		var temp = String($scope.markdownContent).replace(/<[^>]+>/gm, '\n');
		$scope.postData.title = String($scope.postData.title).replace(/<[^>]+>/gm, '');
		var postToSave = {"title": $scope.postData.title,"content":{"markdown":temp, "markup":$scope.markupContent},"media": $scope.mediaData ,"contentDirty": $scope.contentDirty}
		$http.post("/savepost",postToSave).success(function(){
			console.log("Post Saved")
			$timeout(function(){
				$window.location.href="/dashboard"
			}, 2000)
		});
	}
	$scope.$on("uploadComplete", function(evt, file){
		$http.get("/newest_media").success(function(data){
			$scope.recentlyAdded.push(data)
		})
	});
})

breezy.controller('BreezySettingsController', function($scope, $http) {
	$scope.enabledBackups = true;
	$scope.newUser = {"username":"", "password":"", "name":""}
	$scope.submitedNewUser = false;
	$scope.userExists = false;
	$scope.blogname = "";
	$http.get("/blog_info").success(function(data){
		$scope.blogname = data.Name;	
	})
	$scope.updateBlogInfo = function(){
		var dataToUpdate = {"name": $scope.blogname, "searchEnabled": $scope.searchEnabled, "postGroupSize": $scope.PostGroupSize}
		$http.post("/blog_info_update", dataToUpdate).success(function(data){

		})	
	}
	$scope.updateScheduledBackup = function(){
		console.log($scope)
	}
	$scope.backup = function(){
		$http.get("/backup").success(function(data){
			
		});
	}
	$scope.submitNewUser = function(){
		console.log($scope.newUser)
		$scope.submitedNewUser = false;
		$http.post("/newuser", $scope.newUser).success(function(data){
			$scope.submitedNewUser = true;
			console.log(data)
			if(data == "New user created."){
				$scope.userExists = false;
			} else{
				$scope.userExists = true;
			}
		})
	}

});

breezy.directive('breezyActivity', function(){
	return{
		restrict:'E',
		controller: function($scope, $element){

		},
		template: function(){
			return "<div>Activity</div>"
		}
	}
});

breezy.directive('droppable', function(){
	//droppable directive enables an element to have the droppable behavior
	return{
		restrict:'A',
		scope:{
			progress:'='
		},
		controller: function($scope, $element, $files){
			$scope.progress = $files.currentProgress;
			
			$element.on('dragover', function(evt){
				//DragOver State
				evt.dataTransfer.dropEffect ="all";
				if(evt.preventDefault) evt.preventDefault();
			});
			$element.on('dragenter', function(evt){
				//DragEnter State	
			});
			$element.on('dragleave', function(evt){
				//DragLeave State	
			});
			$element.on('drop', function(evt){
				//Drop State
				if(evt.preventDefault) evt.preventDefault();
				if(evt.dataTransfer.files.length >0){
					var filesList = evt.dataTransfer.files;
					$files.upload(filesList);
					
				}
				else{
					$element.append(evt.target.getData("text/html"));
				}
			});
			$scope.$on('started', function(){
				//console.log('upload started');
			});
			$scope.$on('uploadProgress', function(evt, progress){
				//console.log(progress);
			});
			$scope.$on('uploadComplete', function(evt, file){
				$scope.$emit('sortFile', file);	
			});	
			$scope.$on('uploadCompleted', function(){
				$scope.$emit('addElements');
			});
		}
	}
});

breezy.directive('draggable', function(){
	//draggle directive enables an element to have a draggable behavior 
	return{
		restrict:'A',
		scope:{
			sourceElement:'='
		},
		controller: function($scope, $element){
			$scope.sourceElement = $element[0];
			$scope.sourceElement.draggable = true;
			
			$element.on('dragstart', function(evt){
				evt.dataTransfer.effectAllowed ="all";
				evt.dataTransfer.setData('text/html', this.innerHTML);
			});
			$element.on('dragend', function(evt){
				
			});
		}
	}
});

breezy.service('$dateformat', function($rootScope){
	this.formatDate = function(dateString){
		var dateHalf = dateString.split("T")
		var dateParts = dateHalf[0].split("-")
		//console.log(dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0])
		return dateParts[2] + "/" + dateParts[1] + "/" + dateParts[0]
	}
});

breezy.service('$files', function($rootScope,$http){
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


