var breezyDirective = angular.module('breezyDirective', ['ngSanitize'])

breezyDirective.directive('breezyActivity', function(){
	return{
		restrict:'E',
		controller: function($scope, $element){

		},
		template: function(){
			return "<div>Activity</div>"
		}
	}
});

breezyDirective.directive('droppable', function(){
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

breezyDirective.directive('draggable', function(){
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
breezyDirective.directive("breezyNavigation", function(){
	return{
		restrict:'E',
		controller: function($scope, $http, $window){
			$scope.openBlog = function(){
				//Open the users blog in a separate page/tab
				$window.open("/")
			}
			$scope.gotoPage = function(option){
				console.log(option)
				switch(option){
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
				case "home":
					$window.location.href="/dashboard"
					break;
				default:
					console.log("No locaiton to go to.")
					break;
	
				}
			}
		},
		templateUrl : 'views/dashboardnavigation.html'		
	}
});
breezyDirective.directive('contenteditable', ['$sce', function($sce) {
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

breezyDirective.directive('infiniteScroll', function() {
	return {
		restrict: 'A',
		link:function($scope, elemen, attr){
			var container = elemen[0]
	  		elemen.bind("onScroll", function() {
				if(container.scrollTop + container.offsetHeight >= container.scrollHeight - 40){
					console.log("add more")
					$scope.loadMorePosts()
					$scope.$apply()
				}
			})
		},
		controller: function($scope, $element, $window){
			var elm = $element[0]
				//console.log(elm)
			angular.element($window).bind("scroll", function() {
				var scrollWindow = angular.element($window)
					//console.log(elm.scrollHeight - $window.scrollY)
				if(elm.scrollHeight - $window.scrollY < elm.clientHeight + 20){
					//console.log("add more")
					$scope.loadMorePosts()
					$scope.$apply()
				}
			})

		}
	}
})

