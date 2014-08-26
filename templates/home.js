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
		console.log(langA, langB, eType);
	};
});
{{end}}
