package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/placeHolder143032/CodeChallengeHub/web/routes"
    "github.com/placeHolder143032/CodeChallengeHub/middleware"

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

    fs := http.FileServer(http.Dir("./ui"))
    http.Handle("/css/", fs)

    // Frontend Routes (HTML pages)
    mux.HandleFunc("/", routes.GoLandingPage)      // Landing page

    mux.HandleFunc("/login-user", routes.GoLoginUser)   // Login page for regular users
    mux.HandleFunc("/login-admin", routes.GoLoginAdmin) // Login page for admin users
    mux.HandleFunc("/register-user", routes.GoSignupUser)   // Login page for regular users

    
    mux.HandleFunc("/api/auth/login-user", routes.LoginUser)    // User login
    mux.HandleFunc("/api/auth/login-admin", routes.LoginAdmin)    // Admin login
    mux.HandleFunc("/api/auth/register-user", routes.SignupUser)      // User registration
    


    // Protected routes
    mux.HandleFunc("/profile", middleware.RequireAuth(routes.GoProfilePage))

    mux.HandleFunc("/allproblems-user", middleware.RequireAuth(routes.GoProblemsListPageUser))
    mux.HandleFunc("/allproblems-admin", middleware.RequireAuth(routes.GoProblemsListPageAdmin))
    mux.HandleFunc("/add-problem", middleware.RequireAuth(routes.AddProblem)) // add problem page
    mux.HandleFunc("/problem", middleware.RequireAuth(routes.GoProblemPage) )


    mux.HandleFunc("/my_submissions",  middleware.RequireAuth(routes.ViewSubmissionsByUser)) // my submissions page
    mux.HandleFunc("/submit_answer",  middleware.RequireAuth(routes.GoSubmitAnswer)) // Go to submit answer page
    mux.HandleFunc("/api/submit_problem",  middleware.RequireAuth(routes.SubmitAnswer)) 
    mux.HandleFunc("/publish-problem",  middleware.RequireAuth(routes.PublishProblem)) // Go to submit answer page


    return mux
}