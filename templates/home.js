{{define "angular-script"}}
var homeApp = angular.module('homeApp', []);
var allowedLanguages = {};
var langCodeToLabel = {};

homeApp.controller("HomeCtrl", function($scope, $http) {
	$scope.ExtractTypes = {{.ExtractTypes}};
	$scope.AllLanguages = {{.Data.LanguageOptions}};
	$scope.languageAllowed = function(lang) {
		return (lang.Value in allowedLanguages);
	}
	angular.forEach($scope.AllLanguages, function(lang) {
		langCodeToLabel[lang.Value] = lang.Text;
	});
	
	$http({method: 'GET', url: '/api/extract/languages'}).
    success(function(data, status, headers, config) {
			allowedLanguages = {};
			for (var i in data) {
				allowedLanguages[data[i]] = true;
			}
    }).
    error(function(data, status, headers, config) {
			console.log(data, status, headers, config);
    });
		
	$scope.fetchExtracts = function(langA, langB, eType) {
		$scope.showResults = true;
		var query = "";
		if (langA) {
			query += "&langA=" + langA.Value;
		}
		if (langB) {
			query += "&langB=" + langB.Value;
		}
		if (eType) {
			query += "&type=" + eType.Value;
		}
		if (query.length !== 0) {
			query = "?" + query.substr(1);
		}
		$http.get("/api/extract/search" + query).
			success(function(data, status, headers, config) {
				$scope.ResultCount = data.ExtractCount;
				$scope.Results = reorder($scope, data.Results);
			}).
			error(function(data) {
				console.log(data);
			});
	};
});

function reorder($scope, results) {
	angular.forEach(results, function(r) {
		var sums = [];
		var lang1, lang2;
		angular.forEach(r.Summaries, function(s) {
			if ($scope.LanguageA && (s.Language === $scope.LanguageA.Value)) {
				lang1 = s;
			} else if ($scope.LanguageB && (s.Language === $scope.LanguageB.Value)) {
				lang2 = s;
			} else {
				sums.push(s);
			}
		});
		if (lang2) {
			sums.unshift(lang2);
		}
		if (lang1) {
			sums.unshift(lang1);
		}
		r.Summaries = sums;
	});
	return results;
}

homeApp.filter('nonEmpty', function () {
	return function (items, search) {
		var list = [];
		angular.forEach(items, function (value) {
			if (value.Title || value.Summary) {
				list.push(value);
			}
		});
		return list;
	}
});

homeApp.filter('languageLabel', function () {
	return function (langCode) {
		return langCodeToLabel[langCode];
	}
});
{{end}}
