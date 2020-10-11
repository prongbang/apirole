package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/innotechdevops/mgo-driver/pkg/mgodriver"
	"github.com/prongbang/apirole/pkg/apirole"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
)

func main() {
	a, _ := mongodbadapter.NewAdapter("mongodb://root:admin@127.0.0.1:27017/roledb?authSource=admin&ssl=false")
	e, _ := casbin.NewEnforcer("model.conf", a)
	driver := mgodriver.New(mgodriver.Config{
		User:         "root",
		Pass:         "admin",
		Host:         "127.0.0.1",
		DatabaseName: "roledb",
		Port:         mgodriver.DefaultPort,
	})

	_ = e.LoadPolicy()

	_, _ = e.AddPolicy("5f82de37aacb828dc9466173", "/*", "(GET)|(POST)|(PUT)|(DELETE)")
	//_, err := e.RemovePolicy("5f82de37aacb828dc9466173", "/*", "(GET)|(POST)|(PUT)|(DELETE)")
	//log.Println(err)

	app := fiber.New()
	app.Use(fibercasbinrest.NewDefault(e, "secret"))

	// Router
	source := apirole.NewDataSource(driver.Connect())
	repo := apirole.NewRepository(e, source)
	uc := apirole.NewUseCase(repo)
	handle := apirole.NewHandler(uc)
	route := apirole.NewRouter(handle)
	route.Initial(app)

	log.Fatal(app.Listen(":3000"))
}
