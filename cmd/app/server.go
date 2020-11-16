package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Firdavs2002/http/pkg/banners"
)

// Server представляет собой логический сервер нашего приложения.
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

// NewServer - функция-конструктор для создания сервера.
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{
		mux:        mux,
		bannersSvc: bannersSvc,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init инициализация сервра (регистрирует все Handler'ы)
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

func (s *Server) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {

	//берем все баннеры из сервиса
	banners, err := s.bannersSvc.All(r.Context())

	//если получили какую нибуд ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//преобразуем данные в JSON
	data, err := json.Marshal(banners)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//вызываем функцию для ответа в формате JSON
	sendJSON(w, data)

}

func (s *Server) handleGetBannerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP := r.URL.Query().Get("id")

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest)
		return
	}

	//получаем баннер из сервиса
	banner, err := s.bannersSvc.ByID(r.Context(), id)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//преобразуем данные в JSON
	data, err := json.Marshal(banner)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//вызываем функцию для ответа в формате JSON
	sendJSON(w, data)
}

func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {

	//получаем данные из параметра запроса
	idP := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")
	button := r.URL.Query().Get("button")
	link := r.URL.Query().Get("link")

	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest)
		return
	}
	//Здесь опционалная проверка то что если все данные приходит пустыми то вернем ошибку
	if title == "" && content == "" && button == "" && link == "" {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest)
		return
	}

	//создаём указател на структуру баннера
	item := &banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}

	file, fileHeader, err := r.FormFile("image")

	//если нет ошибки значит файл пришел то берем его имя тоест расширения
	if err == nil {
		//Получаем расширеную файла например global.jpg берём только jpg а осталное будем генерироват  в сервисе
		//здес разделяем имя файла по "."
		var name = strings.Split(fileHeader.Filename, ".")
		// берем последный элемент из массива тоест jpg и вставляем его в item.Image (в методе Save его будем менят)
		item.Image = name[len(name)-1]

	}

	//вызываем метод Save тоест сохраняем или обновляем его
	banner, err := s.bannersSvc.Save(r.Context(), item, file)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//преобразуем данные в JSON
	data, err := json.Marshal(banner)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}
	//вызываем функцию для ответа в формате JSON
	sendJSON(w, data)
}

func (s *Server) handleRemoveByID(w http.ResponseWriter, r *http.Request) {

	//извлекаем из параметра запроса ID
	idP := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {

		//печатаем ошибку
		log.Print(err)

		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest)
		return
	}

	banner, err := s.bannersSvc.RemoveByID(r.Context(), id)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	//преобразуем данные в JSON
	data, err := json.Marshal(banner)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//печатаем ошибку
		log.Print(err)
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError)
		return
	}
	//вызываем функцию для ответа в формате JSON
	sendJSON(w, data)
}

//это фукция для записывание ошибки в responseWriter или просто для ответа с ошиками
func errorWriter(w http.ResponseWriter, httpSts int) {
	http.Error(w, http.StatusText(httpSts), httpSts)
}

//это функция для ответа в формате JSON
func sendJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		//печатаем ошибку
		log.Print(err)
	}
}
