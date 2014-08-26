{{define "angular-script"}}
function HomeCtrl($scope, $http) {
	$scope.ExtractTypes = {{.ExtractTypes}};
}
{{end}}
