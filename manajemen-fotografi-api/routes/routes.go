package routes

import (
	"manajemen-fotografi-api/handlers"


	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	
	user := app.Group("/api/users")
	user.Post("/register", handlers.RegisterUser)
	user.Post("/login", handlers.LoginUser)
	user.Post("/logout", handlers.LogoutUser)

	transaction := app.Group("/api/transaction")
	transaction.Post("/transactions", handlers.CreateDummyTransaction)
	transaction.Get("/transactions", handlers.GetAllTransactions)

	

	photographer := app.Group("/photographers")
	photographer.Post("/", handlers.CreatePhotographer)                 // Create new photographer
	photographer.Get("/:id", handlers.GetPhotographerByID)              // Get photographer by ID
	photographer.Get("/user/:user_id", handlers.GetPhotographerByUserID) // Get photographer by user ID
	photographer.Put("/:id", handlers.UpdatePhotographer)               // Update photographer (dengan upload file)
	photographer.Delete("/:id", handlers.DeletePhotographer)            // Delete photographer
        

	// Client routes
	client := app.Group("/api/clients")
	client.Post("/", handlers.CreateClient)
    client.Get("/", handlers.GetAllClients)          // <<< Jangan lupa ini
    client.Get("/:id", handlers.GetClientByID)
    client.Get("/user/:user_id", handlers.GetClientByUserID)
    client.Put("/:id", handlers.UpdateClient)
    client.Delete("/:id", handlers.DeleteClient)


	// Booking Routes
	booking := app.Group("/api/bookings")
	booking.Get("/", handlers.GetAllBookings)
	booking.Get("/:id", handlers.GetBookingByID)
	booking.Post("/", handlers.CreateBooking)
	booking.Put("/:id", handlers.UpdateBooking)
	booking.Delete("/:id", handlers.DeleteBooking)

	// Gallery Routes
	gallery := app.Group("/api/galleries")
	gallery.Get("/", handlers.GetAllGalleries)
	gallery.Get("/:id", handlers.GetGalleryByID)
	gallery.Post("/", handlers.CreateGallery)
	gallery.Put("/:id", handlers.UpdateGallery)
	gallery.Delete("/:id", handlers.DeleteGallery)
}
