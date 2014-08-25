{{define "angular-script"}}
function Ctrl($scope) {
	$scope.languages = {{.Data.LanguageOptions}};
	for (var i in $scope.languages) {
		if ($scope.languages[i].Value === "{{.Data.Language}}") {
			$scope.Language = $scope.languages[i];
			break;
		}
	}
	$scope.flavors = {{.Data.Flavors}};
}
{{end}}
