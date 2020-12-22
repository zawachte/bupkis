package formatter


const (
	defaultImageTableFormat = "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{if .CreatedSince }}{{.CreatedSince}}{{else}}N/A{{end}}\t{{.Size}}"
	imageIDHeader           = "IMAGE ID"
	repositoryHeader        = "REPOSITORY"
	tagHeader               = "TAG"
	digestHeader            = "DIGEST"
	createdSinceHeader      = "CREATED"
	createdAtHeader         = "CREATED AT"
	sizeHeader              = "SIZE"
	labelsHeader            = "LABELS"
	nameHeader              = "NAME"
)

