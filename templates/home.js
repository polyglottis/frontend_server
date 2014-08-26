{{define "angular-script"}}
var homeApp = angular.module('homeApp', []);
var allowedLanguages = {};

homeApp.controller("HomeCtrl", function($scope, $http) {
	$scope.ExtractTypes = {{.ExtractTypes}};
	$scope.AllLanguages = {{.Data.LanguageOptions}};
	$scope.languageAllowed = function(lang) {
		return (lang.Value in allowedLanguages);
	}
	
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
				$scope.Results = data.Results;
			}).
			error(function(data) {
				console.log(data);
			});
	};
});

homeApp.filter('nonEmpty', function () {
	return function (items, search) {
		var list = [];
		angular.forEach(items, function (value, key) {
			if (value.Title || value.Summary) {
				list.push(value);
			}
		});
		return list;
	}
});
{{end}}
