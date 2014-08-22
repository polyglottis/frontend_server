{{define "angular-script"}}
function LanguageCtrl($scope) {
	$scope.languages = {{.Data.LanguageOptions}};
	$scope.changeA = function(lang) {
		$scope.textA = lang.Text[0];
	}
	$scope.changeB = function(lang) {
		if (lang) {
			$scope.textB = lang.Text[0];
		}
	}
	{{if ne .Data.Selection.LanguageA -1}}
	$scope.languageA = $scope.languages[{{.Data.Selection.LanguageA}}];
	$scope.textA = $scope.languageA.Text[{{.Data.Selection.TextA}}];
	{{end}}
}
{{end}}
