package main

import (
	"encoding/json"
	"net/http"
    "io/ioutil"
	"fmt"
	"os"
	"bufio"
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "medium_go_fiber_swagger/docs"
)


// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @contact.name API Support
// @contact.email youremail@provider.com
// @host localhost:3001
// @BasePath /
func main() {
	// Create new Fiber application
	app := fiber.New()

	// Add endpoint to serve swagger documentation
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// Add endpoint to get rate
	app.Get("/rate", rate)

	// Add endpoint to subscribe email
	app.Post("/subscribe", subscribe)

	// Add endpoint to send rate to emails
	app.Post("/sendEmails", sendEmails)

	// Listen on the port '3001'
	app.Listen(":3001")
}


// rate godoc
// @Summary Отримати поточний курс BTC до UAH
// @Description Запит має повертати поточний курс BTC до UAH використовуючи будь-який third party сервіс з публічним АРІ
// @Accept  json
// @Produce  json
// @Tags rate
// @Success 200 {integer} integer 0
// @Failure 400
// @Router /rate [get]
func rate(c *fiber.Ctx) error {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUAH"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        c.Status(400).SendString(err.Error())
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        c.Status(400).SendString(err.Error())
    }
    defer res.Body.Close()
    body, readErr := ioutil.ReadAll(res.Body)
    if readErr != nil {
        c.Status(400).SendString(err.Error())
    }

	data_obj := Response{}
	jsonErr := json.Unmarshal(body, &data_obj)
    if jsonErr != nil {
		c.Status(400).SendString(jsonErr.Error())
    }

	return c.Status(200).SendString(data_obj.Price)
}

// subscribe godoc
// @Summary Підписати емейл на отримання поточного курсу
// @Description Запит має перевірити, чи немає данної електронної адреси в поточній базі даних (файловій) і, в разі її відсутності, записувати її. Пізніше, за допомогою іншого запиту ми будемо відправляти лист на ті електронні адреси, які будуть в цій базі.
// @Accept  json
// @Produce  json
// @Tags subscription
// @Param email formData string true "Електронна адреса, яку потрібно підписати"
// @Success 200 {string} string
// @Failure 409
// @Router /subscribe [post]
func subscribe(c *fiber.Ctx) error {
	email := c.FormValue("email")

	// Open the file in read-only mode to check if the email exists
	file, err := os.Open("emails.log")
	if err != nil {
		if !os.IsNotExist(err) {
			// Handle other file opening errors
			fmt.Println(err)
			return err
		}
		// File doesn't exist, create a new file
		file, err = os.Create("emails.log")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == email {
			return c.Status(409).SendString("Email already exists")
		}
	}

	// Email doesn't exist, append it to the file
	f, err := os.OpenFile("emails.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(email + "\n"); err != nil {
		fmt.Println(err)
		return err
	}

	return c.Status(200).SendString("Email added")
}

// sendEmails godoc
// @Summary Відправити e-mail з поточним курсом на всі підписані електронні пошти.
// @Description Запит має отримувати актуальний курс BTC до UAH за допомогою third-party сервісу та відправляти його на всі електронні адреси, які були підписані раніше.
// @Accept  json
// @Produce  json
// @Tags subscription
// @Success 200 {string} string
// @Failure 400
// @Router /sendEmails [post]
func sendEmails(c *fiber.Ctx) error {


	return c.Status(200).SendString("E-mailʼи відправлено")
}

type Response struct {
	Symbol string 
	Price string
}
