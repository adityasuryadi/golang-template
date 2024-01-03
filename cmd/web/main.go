package main

import (
	"context"
	"fmt"
	"order-service/internal/config"
	"order-service/internal/pkg"
)

// type OrderRequest struct {
// 	ProductId    string  `json:"product_id"`
// 	ProductName  string  `json:"product_name"`
// 	ProductPrice float64 `json:"product_price"`
// 	Qty          int     `json:"qty"`
// }

// type OrderEntity struct {
// 	Id           uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
// 	ProductId    uuid.UUID `gorm:"column:product_id"`
// 	ProductPrice float64   `gorm:"column:product_price"`
// 	Qty          int       `gorm:"column:qty"`
// 	CreatedAt    time.Time `gorm:"column:created_at"`
// 	UpdatedAt    time.Time `gorm:"column:updated_at"`
// }

// func (OrderEntity) TableName() string {
// 	return "orders"
// }

// func (entity *OrderEntity) BeforeCreate(db *gorm.DB) error {
// 	entity.Id = uuid.New()
// 	entity.CreatedAt = time.Now().Local()
// 	return nil
// }

// func (entity *OrderEntity) BeforeUpdate(db *gorm.DB) error {
// 	entity.UpdatedAt = time.Now().Local()
// 	return nil
// }

// func CreateProduct(request OrderRequest) error {
// 	db := InitDB()

// 	entity := OrderEntity{
// 		Id:           uuid.New(),
// 		ProductId:    uuid.MustParse(request.ProductId),
// 		ProductPrice: request.ProductPrice,
// 		Qty:          request.Qty,
// 		CreatedAt:    time.Now().Local(),
// 	}
// 	fmt.Println(entity)

// 	err := db.Create(&entity).Debug().Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

// func InitDB() *gorm.DB {
// 	host := "localhost"
// 	user := "postgres"
// 	password := "postgres"
// 	port := "5433"
// 	db_name := "order_service_db"

// 	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port + " sslmode=disable TimeZone=UTC"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

// // func bodyFrom(args []string) string {
// // 	var s string
// // 	if (len(args) < 2) || os.Args[1] == "" {
// // 		s = "hello"
// // 	} else {
// // 		s = strings.Join(args[1:], " ")
// // 	}
// // 	return s
// // }

// func main() {
// 	app := fiber.New()

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		exchange := "order.created"
// 		queue := "order.create"
// 		routingKey := "create"

// 		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 		failOnError(err, "Failed to connect to RabbitMQ")
// 		defer conn.Close()

// 		ch, err := conn.Channel()
// 		failOnError(err, "Failed to open a channel")
// 		defer ch.Close()

// 		err = ch.ExchangeDeclare(
// 			exchange, // name
// 			"direct", // type
// 			true,     // durable
// 			false,    // auto-deleted
// 			false,    // internal
// 			false,    // no-wait
// 			nil,      // arguments
// 		)
// 		failOnError(err, "Failed to declare an exchange")

// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		q, err := ch.QueueDeclare(
// 			queue, // name
// 			false, // durable
// 			false, // delete when unused
// 			false, // exclusive
// 			false, // no-wait
// 			nil,   // arguments
// 		)
// 		failOnError(err, "Failed to declare a queue")

// 		err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
// 		if err != nil {
// 			failOnError(err, "Failed to declare a queue")
// 		}
// 		log.Print("producer: declaring binding")

// 		// body := bodyFrom(os.Args)
// 		body := "Hello"
// 		err = ch.PublishWithContext(ctx,
// 			exchange,   // exchange
// 			routingKey, // routing key
// 			false,      // mandatory
// 			false,      // immediate
// 			amqp.Publishing{
// 				ContentType: "text/plain",
// 				Body:        []byte(body),
// 			})
// 		failOnError(err, "Failed to publish a message")

// 		log.Printf(" [x] Sent %s", body)
// 		return c.SendString("Order service")
// 	})
// 	app.Post("create", func(c *fiber.Ctx) error {
// 		var request OrderRequest
// 		c.BodyParser(&request)
// 		fmt.Println(request)
// 		err := CreateProduct(request)
// 		if err != nil {
// 			return c.Status(500).JSON(err.Error())
// 		}
// 		return c.Status(200).JSON("ok")
// 	})

// 	app.Listen(":5003")
// }

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	jaegerExporter := config.NewJaegerTracer(viperConfig, log)
	JaegerTracer := pkg.NewJaegerTracer(jaegerExporter)
	tp := JaegerTracer.Trace("order-service")

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	config.Bootstrap(&config.BootstrapConfig{
		DB:             db,
		App:            app,
		Log:            log,
		Validate:       validate,
		Config:         viperConfig,
		JaegerExporter: jaegerExporter,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
