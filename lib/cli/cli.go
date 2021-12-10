package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/paulosimao/ports-api/lib/proto"

	"google.golang.org/grpc"
)

//IfErr - simplifies error handling for HTTP.
func IfErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Printf("Err: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

//Cli - creates a new GRPC Client.
func Cli() (pb.PortDbClient, error) {
	// dial server
	conn, err := grpc.Dial(os.Getenv("GRPC_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// create stream
	client := pb.NewPortDbClient(conn)
	return client, nil
}

//uploadFile - processes file upload.
//Considering we may get a stream of map code:port objects
//It would be better to put code insde of the object and use an array like format
//then we could stream one by one, instead of reading all in memory
//if the format is a strong requirement, then a more elaborated parser will be required.
func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	//MAX 10MB
	r.ParseMultipartForm(10 << 20)

	f, _, err := r.FormFile("ports")
	if IfErr(w, err) {
		return
	}
	defer f.Close()

	client, err := Cli()
	if IfErr(w, err) {
		return
	}

	dec := json.NewDecoder(f)

	for dec.More() {
		in := make(map[string]interface{})
		err = dec.Decode(&in)
		if IfErr(w, err) {
			return
		}
		for k, v := range in {
			datastr, err := json.Marshal(v)
			if IfErr(w, err) {
				return
			}

			port := &pb.PortData{Code: k, Data: string(datastr)}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			_, err = client.PutPort(ctx, port)
			cancel()
			if IfErr(w, err) {

				return
			}

		}
	}

}

//handleGetPorts allows getting port data back
func handleGetPorts(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	client, err := Cli()
	if IfErr(rw, err) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	cli, err := client.GetPorts(ctx, &pb.GetRequest{})
	cancel()
	if IfErr(rw, err) {
		return
	}
	enc := json.NewEncoder(rw)
	for {
		res, err := cli.Recv()
		if err == io.EOF {
			log.Printf("Ended data")
			return
		}
		if IfErr(rw, err) {
			return
		}
		log.Printf("Server sending: %#v", res)
		enc.Encode(res)
	}
}

//Run - Executes the service as a whole
func Run() error {

	addr := os.Getenv("ADDR")
	if addr == "" {
		return errors.New("no ADDR env provided")
	}
	grpc_addr := os.Getenv("GRPC_ADDR")
	if grpc_addr == "" {
		return errors.New("no GRPC_ADDR env provided")
	}

	http.HandleFunc("/get", handleGetPorts)

	http.HandleFunc("/put", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		client, err := Cli()
		if IfErr(rw, err) {
			return
		}

		f, err := os.OpenFile("in/ports.json", os.O_RDONLY, 0600)
		if IfErr(rw, err) {
			return
		}
		dec := json.NewDecoder(f)

		//Considering we may get a stream of map code:port objects
		//It would be better to put code insde of the object and use an array like format
		//then we could stream one by one, instead of reading all in memory
		//if the format is a strong requirement, then a more elaborated parser will be required
		for dec.More() {
			in := make(map[string]interface{})
			err = dec.Decode(&in)
			if IfErr(rw, err) {
				return
			}
			for k, v := range in {

				datastr, err := json.Marshal(v)
				if IfErr(rw, err) {
					return
				}

				port := &pb.PortData{Code: k, Data: string(datastr)}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				_, err = client.PutPort(ctx, port)
				cancel()
				if IfErr(rw, err) {
					return
				}

			}
		}

	})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPorts(rw, r)
		case http.MethodPost:
			handleUploadFile(rw, r)
		default:
			http.Error(rw, fmt.Sprintf("Method %s not implemented.", r.Method), http.StatusNotImplemented)
		}
	})

	return http.ListenAndServe(os.Getenv("ADDR"), nil)
}
