{{define "angular-script"}}
function LanguageCtrl($scope, $window) {
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
	{{if ne .Data.Selection.LanguageB -1}}
	$scope.languageB = $scope.languages[{{.Data.Selection.LanguageB}}];
	$scope.textB = $scope.languageB.Text[{{.Data.Selection.TextB}}];
	{{end}}
	var slug = {{.Slug}};
	$scope.submit = function(args) {
		var loc = "/extract/" + slug + "/" + $scope.languageA.Code + "?ta=" + $scope.textA.Id;
		if ($scope.languageB) {
			loc += "&b=" + $scope.languageB.Code + "&tb=" + $scope.textB.Id;
		}
		$window.location = loc;
	};
}
{{end}}
