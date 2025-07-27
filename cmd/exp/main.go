package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/slmkb/weblensgo/models"
)

func parameterHandler(w http.ResponseWriter, r *http.Request) {
	fetchParam := chi.URLParam(r, "*")
	if fetchParam != "" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Congratz!!!"))
		fmt.Fprintf(w, "%+v", r)

		return
	}
	fmt.Fprintf(w, "%+v", r)
	w.Write([]byte("parameter required"))
}

func dummyHanlder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dummy endpoint"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		panic(err)
	}
	u := struct {
		Name   string
		Age    int
		Nested struct {
			NestedMap   map[string]int
			NestedSlice []string
		}
	}{
		Name: "Kabekaes",
		Age:  9001,
		Nested: struct {
			NestedMap   map[string]int
			NestedSlice []string
		}{
			NestedMap: map[string]int{
				"Value1": 33,
				"Value2": 44,
			},
			NestedSlice: []string{
				"SliceString1",
				"SliceString2",
			},
		},
	}
	if err := tpl.Execute(w, u); err != nil {
		panic(err)
	}

}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	u := struct {
		Name   string
		Age    int
		Nested struct {
			NestedMap   map[string]int
			NestedSlice []string
		}
	}{
		Name: "Kabekaes",
		Age:  9001,
		Nested: struct {
			NestedMap   map[string]int
			NestedSlice []string
		}{
			NestedMap: map[string]int{
				"Value1": 33,
				"Value2": 44,
			},
			NestedSlice: []string{
				"SliceString1",
				"SliceString2",
			},
		},
	}

	if err = t.Execute(os.Stdout, u); err != nil {
		panic(err)
	}
	r := chi.NewRouter()

	r.With(middleware.Logger).Get("/param/{asdf}", parameterHandler)
	r.With(middleware.Logger).Get("/", homeHandler)

	r.Get("/dummy", dummyHanlder)

	r.Get("/lissa", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	})

	if err := CreateUser(); err != nil {
		log.Println(err)
	}

	if err := CreateOrg(); err != nil {
		log.Println(err)
	}

	testDBCreation()

	http.ListenAndServe(":4000", r)

}

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func testDBCreation() {

	cfg := models.DefaultPostgresConfig()
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	us := models.UserService{
		DB: db,
	}
	user, err := us.Create("bob2@bob.com", "bob123")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

// r.With(middleware.Logger).Post("/signin", MiddHandlerFunc(usersCtrl.ExecuteSignIn))
// http.ListenAndServe(":3000", MiddHandlerFunc(r.ServeHTTP))
func MiddHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Middleware-Header", "Value")
		h(w, r)
		log.Printf("Custom header middleware...")
	}
}

// IDIOMATIC
// r.Use(MiddHandler)
// r.With(middleware.Logger, MiddHandler).Post("/signin", usersCtrl.ExecuteSignIn)
// http.ListenAndServe(":3000", MiddHandler(r))
func MiddHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Middleware-Header", "Value")
		next.ServeHTTP(w, r)
		log.Printf("Custom header middleware...")
	})
}
