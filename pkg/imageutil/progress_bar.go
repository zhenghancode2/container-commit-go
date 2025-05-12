package imageutil

import (
	"io"

	"github.com/cheggaaa/pb/v3"
)

// ProgressBarReader wraps an io.Reader and displays a progress bar.
type ProgressBarReader struct {
	io.Reader
	bar *pb.ProgressBar
}

func NewProgressBarReader(r io.Reader, total int64, prefix string) *ProgressBarReader {
	bar := pb.New64(total).Set(pb.Bytes, true)
	if prefix != "" {
		bar.Set("prefix", prefix+" ")
	}
	bar.Start()
	return &ProgressBarReader{
		Reader: r,
		bar:    bar,
	}
}

func (pr *ProgressBarReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.bar.Add(n)
	if err == io.EOF {
		pr.bar.Finish()
	}
	return n, err
}

// ProgressBarWriter wraps an io.Writer and displays a progress bar.
type ProgressBarWriter struct {
	io.Writer
	bar *pb.ProgressBar
}

func NewProgressBarWriter(w io.Writer, total int64, prefix string) *ProgressBarWriter {
	bar := pb.New64(total).Set(pb.Bytes, true)
	if prefix != "" {
		bar.Set("prefix", prefix+" ")
	}
	bar.Start()
	return &ProgressBarWriter{
		Writer: w,
		bar:    bar,
	}
}

func (pw *ProgressBarWriter) Write(p []byte) (int, error) {
	n, err := pw.Writer.Write(p)
	pw.bar.Add(n)
	if pw.bar.Current() >= pw.bar.Total() {
		pw.bar.Finish()
	}
	return n, err
}
