module github.com/step/angmar

go 1.13

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/step/saurontypes v0.0.0-20191127114135-1c7b69a4e64f
	github.com/step/uruk v0.0.0-20191127114036-eb84283fad8d
)

replace github.com/step/sauron_go => ../sauron_go/

replace github.com/step/uruk => ../uruk/

replace github.com/step/saurontypes => ../saurontypes/
