package pdfx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoveWatermarks(t *testing.T) {
	tests := []struct {
		name          string
		inputFile     string
		outputFile    string
		expectedError string
	}{
		{
			name:          "Successful watermark removal",
			inputFile:     "testdata/in2.pdf",
			outputFile:    "testdata/out2.pdf",
			expectedError: "",
		},
		{
			name:          "Successful watermark removal",
			inputFile:     "testdata/in3.pdf",
			outputFile:    "testdata/out3.pdf",
			expectedError: "",
		},
		{
			name:          "Successful watermark removal",
			inputFile:     "testdata/in4.pdf",
			outputFile:    "testdata/out4.pdf",
			expectedError: "",
		},
		{
			name:          "Successful watermark removal",
			inputFile:     "testdata/in5.pdf",
			outputFile:    "testdata/out5.pdf",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

			pdfProcessor, err := New(ctx, tt.inputFile, tt.outputFile)
			assert.NoError(t, err)

			// Call RemoveWatermarks and check the result
			err = pdfProcessor.RemoveWatermarks()
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
