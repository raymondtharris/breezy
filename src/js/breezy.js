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


breezy.controller("BreezyController", function($scope){
	
})

breezy.controller("BreezyLoginController", function($scope, $http){
	$scope.loginCredentials = {"username":"", "password":""}
	
	$scope.submitLoginInfo = function(){
		//submitLoginInfo function sends loginCredentials to see if there is a match
		//and if it works will transfer user to dashboard. 
		console.log($scope.loginCredentials)
			//Send loginCredentials to server to be checked against database
		$http.post("/checkcredentials", $scope.loginCredentials).success(function(data){
			console.log(data)
			
			// if all good go to dashboard
		})
	}
})


breezy.controller("BreezyEditorController", function($scope, $http){
	$scope.preview=false
	$scope.contentDirty=true
	$scope.postData={title:"", dateCreated:"", dateModified:""}
	
	$scope.togglePreview = function(newValue){
		$scope.preview=newValue
		console.log($scope.markdownContent)
		
		if($scope.contentDirty){
			console.log("send to server for translation to markup.")
			var temp = String($scope.markdownContent).replace(/<[^>]+>/gm, '\n');
			var contentToConvert = {"markdownContent": temp, "markupContent":"", "postData": $scope.postData}
			console.log(contentToConvert)
			$http.post("/mdowntomup", contentToConvert).success(function(data){
				console.log(data)
				$scope.markupContent =  data.MarkupContent
				$scope.contentDirty=false	
			});
			
		}
	}
	$scope.savePost = function(){
		console.log($scope)
		//if contentDirty send markdown to change
		//on return send data to be saved in database 			
		var postToSave = {}
		$http.post("savepost",postToSave).success(function(){

		});
	}
})

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

breezy.service('$files', function($rootScope,$http){
	var currentProgress;
	var uploadSize;
	var currentFileSize;
	this.upload = function(fileList){
		currentProgress= 0;
		uploadSize=0;
		for(var i = 0; i < fileList.length; i++){
			uploadSize +=fileList[i].size;
		}
		//console.log(uploadSize);
		for(var i = 0; i < fileList.length; i++){
			var xmlHttpReq = new XMLHttpRequest();
			currentFileSize = fileList[i].size;
			xmlHttpReq.open("POST", "/uploadfile");
			
			xmlHttpReq.setRequestHeader('X_FILE_NAME', fileList[i].name);
			xmlHttpReq.setRequestHeader('X_FILE_SIZE', fileList[i].size);
			xmlHttpReq.setRequestHeader('X-Requested-With', true);
			xmlHttpReq.setRequestHeader('Content-Type', fileList[i].type);
			
			xmlHttpReq.upload.addEventListener("progress", this.uploadProgress, false);
			xmlHttpReq.upload.addEventListener("loadstart", this.uploadStart, false);
			//xmlHttpReq.upload.addEventListener("loadend", this.uploadEnd(xmlHttpReq.upload, fileList[i]), false);
			xmlHttpReq.addEventListener("load", this.uploadComplete(xmlHttpReq, fileList[i]), false);
			//xmlHttpReq.addEventListener("error", this.uploadFailed, false);
			//xhr.addEventListener("abort", uploadCanceled, false)
			xmlHttpReq.send(fileList[i]);
			if(i == fileList.length-1){
				this.allFiles(fileList);
			}
		}
	}
	var getUploadSize = function(){
		return uploadSize;
	}
	this.uploadStart = function(){
		$rootScope.$broadcast('started');
	}
	this.uploadProgress = function(evt){
		
		currentProgress += Math.round(evt.loaded * 100 / uploadSize);
		
		$rootScope.$broadcast('uploadProgress', currentProgress);
	}
	this.uploadComplete = function(req, file){
		//console.log(req);
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


