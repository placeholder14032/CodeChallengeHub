package web

import (
    "fmt"
    "net/http"
    "time"

    "github.com/placeHolder143032/CodeChallengeHub/web/routes"

)

type Server struct {
    Host string
    PORT string
    URL string
}

func (a *Application) Listen() error {

    host := fmt.Sprintf("%s:%s", a.Server.Host, a.Server.PORT)

    server := http.Server{
        Handler:     a.routes(),
        Addr:        host,
        ReadTimeout: 300 * time.Second,
    }


    a.InfoLog.Printf("Server listening on :%s\n", host)

    return server.ListenAndServe()
}

func (a *Application) routes() http.Handler {
    mux := http.NewServeMux()

    // Frontend Routes (HTML pages)
    mux.HandleFunc("/", routes.GoLandingPage)      // Landing page

    mux.HandleFunc("/login-user", routes.GoLoginUser)   // Login page for regular users
    mux.HandleFunc("/login-admin", routes.GoLoginAdmin) // Login page for admin users
    mux.HandleFunc("/register-user", routes.GoSignupUser)   // Login page for regular users
    mux.HandleFunc("/register-admin", routes.GoSignupAdmin) // Login page for admin users

    mux.HandleFunc("/profile", routes.GoProfilePage) // Login page for admin users

    mux.HandleFunc("/problem", routes.GoProblemPage) // Login page for admin users
    mux.HandleFunc("/submit_answer", routes.GoSubmitAnswer) // Go to submit answer page


    mux.HandleFunc("/api/auth/login-User", routes.LoginAdmin)    // Admin login
    mux.HandleFunc("/api/auth/login-admin", routes.LoginUser)    // User login
    mux.HandleFunc("/api/auth/register-user", routes.SignupUser)      // User registration
    mux.HandleFunc("/api/auth/register-admin", routes.SignupAdmin) // Admin registration ???
    mux.HandleFunc("/api/submit_answer", routes.SubmitAnswer) // Go to submit answer page



    return mux
}
