var breezy = angular.module("Breezy",['ngRoute'])

breezy.config(function($routeProvider , $locationProvider){
	$routeProvider.when('/', {
		templateUrl:'dashboard.html',
		controller: 'BreezyController'
	}).when('/admin',{
		templateUrl:'login.html'
		controller:'BreezyLoginController'
	})
})