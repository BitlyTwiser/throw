package pufs_client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"
	"github.com/BitlyTwiser/tinychunk"

	//	"github.com/BitlyTwiser/tinycrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IpfsClient struct {
	Id               int64
	Client           pufs_pb.IpfsFileSystemClient
	Files            []string
	FileUpload       chan string
	DeletedFile      chan string
	FileDeleted      chan bool
	FileDownloadPath string
}

func (c *IpfsClient) UploadFileStream(fileData *os.File, fileSize int64, fileName string) error {
	var wg sync.WaitGroup
	log.Printf("Sending large file.. File Size: %v", fileSize)
	// Look to make the time variables depending on file size as well.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileUpload, err := c.Client.UploadFileStream(ctx)

	if err != nil {
		return err
	}

	//No IPFS hash here, that will not be known until we upload on server
	metadata := &pufs_pb.File{
		Filename:   fileName,
		FileSize:   fileSize,
		IpfsHash:   "",
		UploadedAt: timestamppb.New(time.Now()),
	}

	data := make([]byte, fileSize)
	_, err = fileData.Read(data)

	if err != nil {
		return err
	}

	log.Println("Sending first request")
	// Send metadata request first then data.
	m := &pufs_pb.UploadFileStreamRequest{Data: &pufs_pb.UploadFileStreamRequest_FileMetadata{
		FileMetadata: metadata,
	}}

	if err := fileUpload.Send(m); err != nil {
		log.Printf("Error sending first request: %v", err)
	}

	// Total chunks is utilized here to ensure we add enough wait groups.
	// In this particular case, we are chunking the data into 2MB chunks. If alloted amount was altered, we would be forced to revisit this logic.
	totalChunks := uint(math.Floor(float64(fileSize) / float64((2 << 20))))

	wg.Add(int(totalChunks))
	err = tinychunk.Chunk(data, 2, func(chunkedData []byte) error {
		defer wg.Done()

		log.Println("Sending chunked data")
		if err := fileUpload.Send(&pufs_pb.UploadFileStreamRequest{Data: &pufs_pb.UploadFileStreamRequest_FileData{FileData: chunkedData}}); err != nil {
			log.Printf("Error sending file: %v", err)
			return err
		}

		return nil
	})

	wg.Wait()

	if err != nil {
		log.Printf("Error chunking and sending data: %v", err)
		return err
	}

	resp, err := fileUpload.CloseAndRecv()

	if err != nil {
		log.Printf("No response from server")
		return err
	}

	if resp.GetSucessful() {
		log.Println("File has been uploaded")
	} else {
		return errors.New("server did not say successful")
	}

	return nil
}

func (c *IpfsClient) UploadFile(path, fileName string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)

	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	//gRPC data size cap at 4MB
	if fileSize >= (2 << 21) {
		log.Println("Sending big file")
		err = c.UploadFileStream(file, fileSize, fileName)

		if err != nil {
			return err
		}
	} else {
		log.Printf("Sending file of size: %v", fileSize)

		fileData := make([]byte, fileSize)
		_, err := file.Read(fileData)

		if err != nil {
			return err
		}

		err = c.UploadFileData(fileData, fileSize, fileName)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *IpfsClient) DeleteFile(fileName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := c.Client.DeleteFile(ctx, &pufs_pb.DeleteFileRequest{FileName: fileName})

	if err != nil {
		return err
	}

	if resp.Successful {
		log.Println("File deleted")
	} else {
		return fmt.Errorf("error occured deleting file: %v", resp)
	}

	// Push onto a different channel for refreshing files.
	c.DeletedFile <- fileName
	c.FileDeleted <- true

	return nil
}

func (c *IpfsClient) DownloadCappedFile(fileName, path string) error {
	log.Printf("Downloading larger file: %v", fileName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pufs_pb.DownloadFileRequest{FileName: fileName}

	download, err := c.Client.DownloadFile(ctx, req)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("%v/%v", path, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		log.Printf("error opening file to store downloaded data: %v", err)
	}

	for {
		fileChunk, err := download.Recv()

		if err == io.EOF {
			log.Printf("All data downloaded")

			break
		}

		if err != nil {
			log.Printf("Error downloading capped file: %v", err)
			return err
		}

		n, err := file.Write(fileChunk.GetFileData())

		if err != nil {
			return err
		}

		if n == 0 {
			return errors.New("no bytes were written to file")
		}
	}

	return nil
}

// We must chunk the file here if its over the 4MB limit.
func (c *IpfsClient) DownloadFile(fileName, path string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Downloading file: %v", fileName)

	fileResp, err := c.Client.DownloadUncappedFile(ctx, &pufs_pb.DownloadFileRequest{FileName: fileName})

	if err != nil {
		return err
	}

	fileData, fileMetadata := fileResp.FileData, fileResp.FileMetadata

	log.Println("Downloading file and saving to disk...")

	err = os.WriteFile(fmt.Sprintf("%v/%v", path, fileMetadata.Filename), fileData, 0600)

	if err != nil {
		return err
	}

	return nil
}

//Uploads a file stream that is under the 4MB gRPC file size cap
func (c *IpfsClient) UploadFileData(fileData []byte, fileSize int64, fileName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file := &pufs_pb.File{
		Filename:   fileName,
		FileSize:   fileSize,
		IpfsHash:   "",
		UploadedAt: timestamppb.New(time.Now()),
	}

	log.Println("Uploading file")

	request := &pufs_pb.UploadFileRequest{FileData: fileData, FileMetadata: file}
	resp, err := c.Client.UploadFile(ctx, request)

	if err != nil {
		return err
	}

	if !resp.Sucessful {
		return errors.New("something went wrong uploading file")
	}

	return nil
}

// Load files upon client start.
func (c *IpfsClient) LoadFiles() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := c.Client.ListFiles(ctx, &pufs_pb.FilesRequest{})

	if err != nil {
		fyne.NewNotification("Error", fmt.Sprintf("Error loading files from server. Error: %v", err))

		return
	}

	for {
		file, err := req.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading file stream. Error: %v", err)
			break
		}

		c.Files = append(c.Files, file.Files.Filename)
	}
}

// Listen for file changes realtime.
// Take ID and store this upstream.
func (c *IpfsClient) SubscribeFileStream() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		stream, err := c.Client.ListFilesEventStream(ctx, &pufs_pb.FilesRequest{Id: c.Id})

		if err != nil || stream == nil {
			log.Println("Error or stream not empty, waiting for 5 seconds")
			time.Sleep(time.Second * 5)

			continue
		} else {
			file, err := stream.Recv()

			if err == io.EOF {
				log.Println("All files read, awaiting..")
				stream = nil
				break
			}

			if err != nil {
				log.Println("error encountered, retrying..")
				stream = nil
				break
			}
			log.Printf("Pushing file.. Filename: %v", file.Files.Filename)

			if len(c.FileDeleted) == 0 {
				c.FileUpload <- file.Files.Filename
			} else {
				<-c.FileDeleted
			}
		}
	}
}

func (c *IpfsClient) ChunkFile(fileName string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	size, err := c.Client.FileSize(ctx, &pufs_pb.FileSizeRequest{FileName: fileName})

	if err != nil {
		log.Printf("Could not get file size. Error: %v", err)
	}

	return size.FileSize >= (2 << 20)
}

func (c *IpfsClient) UnsubscribeClient() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.Client.UnsubscribeFileStream(ctx, &pufs_pb.FilesRequest{Id: c.Id})
}
