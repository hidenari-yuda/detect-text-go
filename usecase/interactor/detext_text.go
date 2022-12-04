package interactor

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/api/option"
)

type DetectTextRoute struct {
}

//Mid

// uploadCloudStorage 画像をCloud Storageにアップロード
func UploadCloudStorage(w io.Writer, filePath string) error {
	// key.jsonはサービスアカウントを作成してゲットする
	credentialFilePath := "credential.json"
	// credentialFilePath := Cfg.Google.ApplicationCredentials
	fmt.Println(credentialFilePath)

	// クライアントを作成する
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	// GCSオブジェクトを書き込む空のsample.txtファイルをローカルで作成
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// オブジェクトのReaderを作成
	bucketName := "experimental-bucket-tt"   // e.g. example-bucket
	objectPath := "sample-object/sample.txt" // e.g. foo/var/sample.txt

	writer := client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)
	if _, err := io.Copy(writer, f); err != nil {
		panic(err)
	}

	if err := writer.Close(); err != nil {
		panic(err)
	}

	log.Println("done")

	return nil
}

// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
func GetLabel() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	filename := "/rectangle_large_type_2_e18726aa9c3084139fe5af4c22187e68.webp"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
	}
}

// detectText gets text from the Vision API for an image at the given file path.
func DetectText(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		for index, annotation := range annotations {
			fmt.Fprintf(w, "%q\n", annotation.Description)
			//mid
			fmt.Fprintf(w, "Mid is %v index: %d", annotation.Mid, index)
			//Locale
			// The language code for the locale in which the entity textual
			// `description` is expressed.
			fmt.Fprintf(w, "Locale is %v index: %d", annotation.Locale, index)

			//Description
			// Entity textual description, expressed in its `locale` language.
			fmt.Fprintf(w, "Description is %v index: %d", annotation.Description, index)

			//Score
			// Overall score of the result. Range [0, 1].
			fmt.Fprintf(w, "Score is %v index: %d", annotation.Score, index)

			//Confidence
			// **Deprecated. Use `score` instead.**
			// The accuracy of the entity detection in an image.
			// For example, for an image in which the "Eiffel Tower" entity is detected,
			// this field represents the confidence that there is a tower in the query
			// image. Range [0, 1].
			// Deprecated: Do not use.
			fmt.Fprintf(w, "Confidence is %v index: %d", annotation.Confidence, index)

			//Topicality
			// The relevancy of the ICA (Image Content Annotation) label to the
			// image. For example, the relevancy of "tower" is likely higher to an image
			// containing the detected "Eiffel Tower" than to an image containing a
			// detected distant towering building, even though the confidence that
			// there is a tower in each image may be the same. Range [0, 1].
			fmt.Fprintf(w, "Topicality is %v index: %d", annotation.Topicality, index)

			//BoundingPoly
			// Image region to which this entity belongs. Not produced
			// for `LABEL_DETECTION` features.
			fmt.Fprintf(w, "BoundingPoly is %v index: %d", annotation.BoundingPoly, index)

			//Locations
			// The location information for the detected entity. Multiple
			// `LocationInfo` elements can be present because one location may
			// indicate the location of the scene in the image, and another location
			// may indicate the location of the place where the image was taken.
			// Location information is usually present for landmarks.
			fmt.Fprintf(w, "Locations is %v index: %d", annotation.Locations, index)

			//Properties
			// Some entities may have optional user-supplied `Property` (name/value)
			// fields, such a score or string that qualifies the entity.
			fmt.Fprintf(w, "Properties is %v index: %d", annotation.Properties, index)
		}

		// if annotation.BoundingPoly != nil {
		// 	fmt.Fprintf(w, "\tBounds: %v\n", annotation.BoundingPoly)
		// }
		// }
	}

	return nil
}

// detectAsyncDocumentURI performs Optical Character Recognition (OCR) on a
// PDF file stored in GCS.
func DetectAsyncDocumentURI(w io.Writer, gcsSourceURI, gcsDestinationURI string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	request := &visionpb.AsyncBatchAnnotateFilesRequest{
		Requests: []*visionpb.AsyncAnnotateFileRequest{
			{
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
					},
				},
				InputConfig: &visionpb.InputConfig{
					GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
					// Supported MimeTypes are: "application/pdf" and "image/tiff".
					MimeType: "application/pdf",
				},
				OutputConfig: &visionpb.OutputConfig{
					GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
					// How many pages should be grouped into each json output file.
					BatchSize: 2,
				},
			},
		},
	}

	operation, err := client.AsyncBatchAnnotateFiles(ctx, request)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Waiting for the operation to finish.")

	resp, err := operation.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%v", resp)

	return nil
}

//BatchAnnotateFiles ファイルパスからテキスト検出。多分無理っぽい
func DetectBatchAnnotateFiles(w io.Writer, filePath string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	request := &visionpb.BatchAnnotateFilesRequest{
		Requests: []*visionpb.AnnotateFileRequest{
			{
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
					},
				},
				InputConfig: &visionpb.InputConfig{
					Content: []byte("Hello, world!"),
					// GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
					// Supported MimeTypes are: "application/pdf" and "image/tiff".
					MimeType: "application/pdf",
				},
				// OutputConfig: &visionpb.OutputConfig{
				// 	GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
				// 	// How many pages should be grouped into each json output file.
				// 	BatchSize: 2,
				// },
			},
		},
	}

	operation, err := client.BatchAnnotateFiles(ctx, request)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Waiting for the operation to finish.")

	fmt.Fprintf(w, "%v", operation.Responses)

	return nil
}

// detectDocumentText gets the full document text from the Vision API for an image at the given file path.
func DetectDocumentText(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return err
	}

	if annotation == nil {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Document Text:")
		fmt.Fprintf(w, "%q\n", annotation.Text)

		fmt.Fprintln(w, "Pages:")
		for _, page := range annotation.Pages {
			fmt.Fprintf(w, "\tConfidence: %f, Width: %d, Height: %d\n", page.Confidence, page.Width, page.Height)
			fmt.Fprintln(w, "\tBlocks:")
			for _, block := range page.Blocks {
				fmt.Fprintf(w, "\t\tConfidence: %f, Block type: %v\n", block.Confidence, block.BlockType)
				fmt.Fprintln(w, "\t\tParagraphs:")
				for _, paragraph := range block.Paragraphs {
					fmt.Fprintf(w, "\t\t\tConfidence: %f", paragraph.Confidence)
					fmt.Fprintln(w, "\t\t\tWords:")
					for _, word := range paragraph.Words {
						symbols := make([]string, len(word.Symbols))
						for i, s := range word.Symbols {
							symbols[i] = s.Text
						}
						wordText := strings.Join(symbols, "")
						fmt.Fprintf(w, "\t\t\t\tConfidence: %f, Symbols: %s\n", word.Confidence, wordText)
					}
				}
			}
		}
	}

	return nil
}

// setEndpoint changes your endpoint.
func SetEndpoint(endpoint string) error {
	// endpoint := "eu-vision.googleapis.com:443"

	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("NewImageAnnotatorClient: %v", err)
	}
	defer client.Close()

	return nil
}
