var breezy = angular.module("Breezy",['ngRoute'])
/*
breezy.config(function($routeProvider , $locationProvider){
	$routeProvider.when('/', {
		templateUrl:'/views/dashboard.html',
		controller: 'BreezyController'
	}).when('/admin',{

		controller:'BreezyLoginController'
	})
})
*/
breezy.controller("BreezyController", function($scope){
	
})

breezy.controller("BreezyLoginController", function($scope, $http, $location){
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
	$scope.markdownContent = "# Title fjansd"
	$scope.togglePreview = function(newValue){
		$scope.preview=newValue
		
		
		if($scope.contentDirty){
			console.log("send to server for translation to markup.")
			
		}
	}
})


