package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gobuffalo/envy"
	"github.com/rs/cors"
	"github.com/shawntoffel/darksky"
)

// Status rgrerg
type Status struct {
	Intensity   darksky.Measurement `json:"intensity"`
	Probability darksky.Measurement `json:"probability"`
	Message     string              `json:"message"`
	Icon        string              `json:"icon"`
}

func main() {
	// Load .env file
	envy.Load(".env", "$GOROOT/src/github.com/lasseh/regnerdetibergen/.env")

	// Create a new router
	r := chi.NewRouter()

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Middleware
	r.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.RealIP,          // Logs the real ip from nginx
		middleware.Logger,          // Log API request calls
		middleware.Recoverer,       // Recover from panics without crashing server
		middleware.RequestID,       // Injects a request ID into the context of each request
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		cors.Handler,               // Handle CORS headers
	)

	// /api/status
	r.Get("/api/status", RainStatus)

	// Serve Vue
	workDir, _ := os.Getwd()
	webDir := filepath.Join(workDir, "web/dist")
	fileServer(r, "/", http.Dir(webDir))

	// Start the router
	fmt.Println("Starting server..")
	http.ListenAndServe(":9003", r)
}

// RainStatus does everything, gets from darksky, calculate the message
func RainStatus(w http.ResponseWriter, r *http.Request) {
	var s Status
	c := darksky.New(envy.Get("DARK_KEY", ""))
	req := darksky.ForecastRequest{}
	req.Latitude = 60.3919
	req.Longitude = 5.3228
	req.Options = darksky.ForecastRequestOptions{Exclude: "hourly,minutely", Lang: "no", Units: "ca"}

	resp, err := c.Forecast(req)
	if err != nil {
		fmt.Println(err.Error())
		s.Message = "Skyen svarer ikke, se ut vinduet."
		render.JSON(w, r, s)
		return
	}
	s.Intensity = resp.Currently.PrecipIntensity
	s.Probability = resp.Currently.PrecipProbability
	s.Icon = resp.Currently.Icon

	// Calculate current weather
	switch {
	case s.Intensity < 0.0 && s.Probability < 0.0:
		s.Message = "Nei, faktisk ikke!"
	case s.Intensity < 0.1 && s.Probability > 0.15:
		s.Message = "Nei, men ser ut som det skal!"
	case s.Intensity > 0.1 && s.Intensity < 0.5:
		s.Message = "Ja for faen!"
	case s.Intensity > 0.5:
		s.Message = "Omg ja, hold deg inne!"
	default:
		s.Message = "Usikker, se ut vinduet!"
	}

	render.JSON(w, r, s)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
