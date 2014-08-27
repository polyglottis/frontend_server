{{define "angular-script"}}
function LanguageCtrl($scope, $window) {
	$scope.languages = {{.Data.LanguageOptions}};
	$scope.change = function(lang) {
		if (lang) {
			$scope.text = lang.Text[0];
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
	{{if .FocusOnA}}
		$scope.language = $scope.languageB;$scope.text = $scope.textB;
	{{else}}
		$scope.language = $scope.languageA;$scope.text = $scope.textA;
	{{end}}
	var slug = {{.Slug}};
	$scope.submit = function(args) {
		var q = "?a={{.Focus.Language}}&at={{.Focus.Id}}&focus=a";
		if ($scope.language) {
			{{if .FocusOnA}}
				var pre = "b";
			{{else}}
				var pre = "a";
				q = "?b={{.Focus.Language}}&bt={{.Focus.Id}}&focus=b";
			{{end}}
			q += "&" + pre + "=" + $scope.language.Code + "&" + pre + "t=" + $scope.text.Id;
		}
		var loc = "/extract/edit/text/" + slug;
		$window.location = loc + q;
	};
}
{{end}}
