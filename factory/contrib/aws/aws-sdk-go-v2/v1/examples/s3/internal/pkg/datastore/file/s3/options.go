package s3

// Options of the s3 repository.
// @Config
type Options struct {
	Bucket string `default:"example" desc:"bucket the files are kept in"`
}
