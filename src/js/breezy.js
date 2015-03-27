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

breezy.controller("BreezyLoginController", function($scope, $http){
	$scope.loginCredentials = {"username":"", "password":""}
	
	$scope.submitLoginInfo = function(){
		console.log($scope.loginCredentials)
		$http.post("/checkcredentials", $scope.loginCredentials).success(function(data){
			
		})
	}
})