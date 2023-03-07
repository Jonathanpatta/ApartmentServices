package Files

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Middleware"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type S3FileHttpService struct {
	service           *S3FileService
	middlewareService *Middleware.MiddlwareService
}

func NewS3FileHttpService(settings *Settings.Settings) (*S3FileHttpService, error) {
	service, err := NewS3FileService(settings)
	if err != nil {
		return nil, err
	}

	return &S3FileHttpService{
		service:           service,
		middlewareService: settings.MiddlewareService,
	}, nil
}

// handler to handle the image upload
func (s *S3FileHttpService) UploadImages(w http.ResponseWriter, r *http.Request) {
	user := Middleware.GetFirebaseUser(r.Context())

	// 32 MB is the default used by FormFile() function
	BulkFileSize := int64(4 * 1024 * 1024)
	if err := r.ParseMultipartForm(BulkFileSize); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	var errNew string
	var httpStatus int

	var images []Images

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			errNew = err.Error()
			httpStatus = http.StatusInternalServerError
			break
		}

		defer file.Close()

		buff := make([]byte, 5*1024*1024)
		_, err = file.Read(buff)
		if err != nil {
			errNew = err.Error()
			httpStatus = http.StatusInternalServerError
			break
		}

		// checking the content type
		// so we don't allow files other than images
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" && filetype != "application/octet-stream" {
			errNew = "The provided file format is not allowed. Please upload a JPEG,JPG or PNG image"
			httpStatus = http.StatusBadRequest
			break
		}

		name := user.UserId + "_" + fileHeader.Filename
		images = append(images, Images{
			Bytes:    buff,
			FileType: filetype,
			FileName: name,
		})
	}

	if errNew != "" {
		message := errNew
		messageType := "E"

		resp := map[string]interface{}{
			"messageType": messageType,
			"message":     message,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(resp)
		return
	}

	urls, err := s.service.UploadImages(images)
	if err != nil {
		errNew = err.Error()
		httpStatus = http.StatusInternalServerError
	}

	outData, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func AddSubrouter(r *mux.Router, settings *Settings.Settings) {
	server, err := NewS3FileHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/files").Subrouter()

	//router.Use(settings.MiddlewareService.ValidateToken)

	router.HandleFunc("/uploadImages", server.UploadImages).Methods("POST", "OPTIONS")
	//router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
}
