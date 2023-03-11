package Files

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Middleware"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
	"strings"
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

func (s *S3FileHttpService) UploadImagesBase64(w http.ResponseWriter, r *http.Request) {
	user := Middleware.GetFirebaseUser(r.Context())

	var imgData []Images

	err := json.NewDecoder(r.Body).Decode(&imgData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, img := range imgData {
		data := make([]byte, base64.StdEncoding.DecodedLen(len(img.Data)))
		_, err := base64.StdEncoding.Decode(data, []byte(img.Data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imgData[i].Bytes = data
		imgData[i].Data = ""
		name := user.UserId + "_" + imgData[i].FileName
		imgData[i].FileName = name
	}

	urls, err := s.service.UploadImages(imgData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

		buff := make([]byte, 10*1024*1024)
		n, err := file.Read(buff)
		if err != nil {
			errNew = err.Error()
			httpStatus = http.StatusInternalServerError
			break
		}

		buff = buff[:n]

		fmt.Println(n, len(buff))

		// checking the content type
		// so we don't allow files other than images
		filetype := http.DetectContentType(buff)
		fmt.Println(!strings.Contains(filetype, "image/"), filetype)
		if !strings.Contains(filetype, "image/") && filetype != "application/octet-stream" {
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

	router.Use(settings.MiddlewareService.ValidateToken)

	router.HandleFunc("/uploadImages", server.UploadImages).Methods("POST", "OPTIONS")
	router.HandleFunc("/uploadImagesBase64", server.UploadImagesBase64).Methods("POST", "OPTIONS")
	//router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
}
