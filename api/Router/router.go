package routers

import (
    "net/http"
    "your_project/api/handlers"
)

func NewRouter() *http.ServeMux {
    router := http.NewServeMux()

    router.HandleFunc("/products", handlers.CreateProduct)

    return router
}
