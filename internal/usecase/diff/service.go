package diff

import (
	"code/internal/domain/diff"
	"code/internal/interfaces"
)

type Service struct {
	parser    interfaces.Parser
	formatter interfaces.Formatter
}

func NewService(parser interfaces.Parser, formatter interfaces.Formatter) *Service {
	return &Service{parser: parser, formatter: formatter}
}

func (s *Service) GenerateDiff(leftRaw, rightRaw []byte) (string, error) {
	left, err := s.parser.Parse(leftRaw)
	if err != nil {
		return "", err
	}

	right, err := s.parser.Parse(rightRaw)
	if err != nil {
		return "", err
	}

	nodes := diff.BuildDiff(left, right)

	return s.formatter.Format(nodes)
}
