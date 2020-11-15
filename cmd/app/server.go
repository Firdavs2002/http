package app

import (
	"net/http"
)

// Server представляет собой логический сервер нашего приложения.
type Serer struct {
	mux *http.ServeMux
}

// NewServer - функция-конструктор для создания сервера.
func NewServer(mux *http.ServeMux) *Server {
	return &Server{mux: mux}
}

func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writerm request)
}

// Init инициализация сервра (регистрирует все Handler'ы)
func (s *Serevr) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}