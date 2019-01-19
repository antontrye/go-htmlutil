package htmlutil

import (
	"golang.org/x/net/html"
)

type (
	filterNodesConfig struct {
		Node    *html.Node
		Filters []func(node *html.Node) bool
		Find    bool
	}
)

func copyValidFilters(input []func(node *html.Node) bool) (output []func(node *html.Node) bool) {
	for _, filter := range input {
		if filter != nil {
			output = append(output, filter)
		}
	}
	return
}

func filterNodes(config filterNodesConfig) []*html.Node {
	config.Filters = copyValidFilters(config.Filters)

	var (
		result []*html.Node
		fn     func(config filterNodesConfig)
	)

	fn = func(config filterNodesConfig) {
		if config.Node == nil {
			return
		}

		if config.Find && len(result) != 0 {
			return
		}

		if len(config.Filters) == 0 {
			result = append(result, config.Node)
			return
		}

		start := len(result)

		func(config filterNodesConfig) {
			var filter func(node *html.Node) bool

			for filter == nil && len(config.Filters) != 0 {
				filter = config.Filters[0]
				config.Filters = config.Filters[1:]
			}

			if filter != nil && !filter(config.Node) {
				return
			}

			if len(config.Filters) == 0 {
				fn(config)

				return
			}

			for c := config.Node.FirstChild; c != nil; c = c.NextSibling {
				config.Node = c

				fn(config)
			}
		}(config)

		finish := len(result)

		for c := config.Node.FirstChild; c != nil; c = c.NextSibling {
			config.Node = c

			fn(config)

			for i := start; i < finish; i++ {
				for j := finish; j < len(result); j++ {
					if result[i] != result[j] {
						continue
					}

					copy(result[j:], result[j+1:])
					result[len(result)-1] = nil
					result = result[:len(result)-1]
					j--
				}
			}
		}
	}

	fn(config)

	return result
}

func encodeText(node *html.Node) []byte {
	if node == nil {
		return nil
	}
	if node.Type == html.TextNode {
		return []byte(node.Data)
	}
	var b []byte
	for node := node.FirstChild; node != nil; node = node.NextSibling {
		b = append(b, encodeText(node)...)
	}
	return b
}