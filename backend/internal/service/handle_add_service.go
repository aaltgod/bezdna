package service

import (
	"net/http"
)

func (s *service) AddService(w http.ResponseWriter, req *http.Request) {
	// if err := s.dbRepository.InsertService(); err != nil {
	// 	log.Error(err)

	// 	// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	// 	return
	// }

	// if err := s.sniffer.AddConfig(sniffer.Config{
	// 	ServiceName: "pop",
	// 	Port:        8080,
	// }); err != nil {
	// 	log.Println(err)

	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	// 	return
	// }

	// if err := s.sniffer.AddConfig(sniffer.Config{
	// 	ServiceName: "fig",
	// 	Port:        4554,
	// }); err != nil {
	// 	log.Println(err)

	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	// 	return
	// }
}
