var dgApp = angular.module('DgApp', ['ngRoute']);

dgApp.config(['$routeProvider', function($routeProvider) {
    $routeProvider
        .when('/', {
            templateUrl: 'html/home.html',
            controller: 'TabController'
        })
        .when('/home', {
            templateUrl: 'html/home.html',
            controller: 'TabController'
        })
        .when('/find', {
            templateUrl: 'html/find.html',
            controller: 'TabController'
        })
        .otherwise({
            redirectTo: '/'
        });
}]);

dgApp.controller('TabController', function($scope, $http) {
    $scope.ticker = "";
    $scope.dgappInfo = {
        hideResult: true,
    };
    $scope.stocks = [];
    $scope.findStock = function() {
        if ($scope.ticker == '') {
            $http.get('/api/').
                success(function(data, status, headers, config) {
                    $scope.stocks = data;
                    $scope.dgappInfo.hideResult = false;
                }).
                error(function(data, status, headers, config) {
                    $scope.dgappInfo.hideResult = true;
                });

            return;
        }

        $http.get('/api/'.concat($scope.ticker)).
            success(function(data, status, headers, config) {
                $scope.stocks = [];
                $scope.stocks.push(data);
                $scope.dgappInfo.hideResult = false;
            }).
            error(function(data, status, headers, config) {
                $scope.dgappInfo.hideResult = true;
            });

        $scope.stock = '';
    }
});
