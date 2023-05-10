package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

func uploadFile(uploader *s3manager.Uploader, filePath string, bucketName string, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	return err
}

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "log/logfile.log",
		MaxSize:    1, // megabytes
		MaxBackups: 10,
		MaxAge:     90,   //days
		Compress:   true, // disabled by default
	})
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Example start", "v0.0.5")

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Viper read config err", err)
		log.Println("Viper read config err", err)
		return
	}

	awsRegion := viper.GetString("AWS_REGION")
	awsID := viper.GetString("AWS_ID")
	awsSecret := viper.GetString("AWS_SECRET")
	awsBucketName := viper.GetString("AWS_BUCKET_NAME")

	if awsRegion == "" && awsID == "" && awsSecret == "" && awsBucketName == "" {
		fmt.Println("Please check env file")
		log.Println("Please check env file")
		return
	}

	flag.String("file", "", "please select file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println("Viper bind flag err", err)
		log.Println("Viper bind flag err", err)
		return
	}

	//todo config
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsID, awsSecret, ""),
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		log.Printf("Failed to initialize new session: %v", err)
		return
	}

	file := viper.GetString("file")
	if file == "" {
		fmt.Println("please select file")
		log.Println("please select file")
		return
	}
	bucketName := awsBucketName
	uploader := s3manager.NewUploader(sess)
	filename := file

	err = uploadFile(uploader, file, bucketName, filename)
	if err != nil {
		fmt.Printf("Failed to upload file: %v", err)
		log.Printf("Failed to upload file: %v", err)
		return
	}
	fmt.Println("Successfully uploaded file!")
	log.Println("Successfully uploaded file!")
}
