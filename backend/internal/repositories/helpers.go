package repositories

import "modulate/backend/internal/models"

func BuildCommentTree(comments []models.Comment) []models.Comment {
	children := make(map[int64][]models.Comment)
	var roots []models.Comment

	// Group comments by parent
	for _, c := range comments {
		if c.ParentCommentID == nil {
			roots = append(roots, c)
		} else {
			children[*c.ParentCommentID] =
				append(children[*c.ParentCommentID], c)
		}
	}

	// Recursively attach children to comments
	var attach func(*models.Comment)
	attach = func(c *models.Comment) {
		c.Children = children[c.ID]
		for i := range c.Children {
			attach(&c.Children[i])
		}
	}

	for i := range roots {
		attach(&roots[i])
	}

	return roots
}
