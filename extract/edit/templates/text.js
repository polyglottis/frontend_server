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
		{{if .FocusOnA}}
			var pre = "b";
			var q = "?a={{.Focus.Language}}&at={{.Focus.Id}}&focus=a";
		{{else}}
			var pre = "a";
			var q = "?b={{.Focus.Language}}&bt={{.Focus.Id}}&focus=b";
		{{end}}
		var loc = "/extract/edit/text/" + slug;
		if ($scope.language) {
			q += "&" + pre + "=" + $scope.language.Code + "&" + pre + "t=" + $scope.text.Id;
		}
		$window.location = loc + q;
	};
}
{{end}}
