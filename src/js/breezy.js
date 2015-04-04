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

/*directive("contenteditable", function() {
    return {
        restrict: "A",
        require: "ngModel",
        link: function(scope, element, attrs, ngModel) {

            function read() {
                ngModel.$setViewValue(element.html());
            }
            ngModel.$render = function() {
                element.html(ngModel.$viewValue || "");
            };
            element.bind("blur keyup change", function() {
                scope.$apply(read);
            });


            element.css('outline','none');
            element.bind("keydown keypress", function (event) {
              console.log("aoidfja")
                if(event.which === 13) {
                    element[0].blur();
                    event.preventDefault();
                  
                }
            });

        }
    };
});
*/

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
			var contentToConvert = {"markdownContent":$scope.markdownContent, "markupContent":""}
			//console.log(contentToConvert)
			$http.post("/mdowntomup", contentToConvert).success(function(data){
				//console.log(data)
				$scope.markupContent =  data.markupContent
				$scope.contentDirty=false	
			});
			
		}
	}
})


