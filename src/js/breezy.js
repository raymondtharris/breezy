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
		console.log($scope.loginCredentials)
		$http.post("/checkcredentials", $scope.loginCredentials).success(function(data){
			console.log(data)
			
			// if all good go to dashboard
		})
	}
})


breezy.controller("BreezyEditorController", function($scope, $http){
	$scope.preview=false
	$scope.contentDirty=true
	//$scope.markdownContent = "# Title fjansd"
	$scope.togglePreview = function(newValue){
		$scope.preview=newValue
		console.log($scope.markdownContent)
		
		if($scope.contentDirty){
			console.log("send to server for translation to markup.")
			 //String($scope.markdownContent).replace("<br>", "\n");
			//console.log(temp)
			var temp = String($scope.markdownContent).replace(/<[^>]+>/gm, '\n');
			console.log(temp)
			var contentToConvert = {"markdownContent": temp, "markupContent":""}
			console.log(contentToConvert)
			$http.post("/mdowntomup", contentToConvert).success(function(data){
				console.log(data)
				$scope.markupContent =  data.MarkupContent
				$scope.contentDirty=false	
			});
			
		}
	}
})


