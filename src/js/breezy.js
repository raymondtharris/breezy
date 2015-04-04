var breezy = angular.module('breezyApp',['ngSanitize'])

breezy.directive("contenteditable", function() {
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
			$scope.markupContent = $scope.markdownContent;
		}
	}
})


