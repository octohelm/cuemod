package golang

import (
	"strconv"

	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
)

func newInt(i int) *cueast.BasicLit {
	return cueast.NewLit(cuetoken.INT, strconv.Itoa(i))
}

func oneOf(types ...cueast.Expr) cueast.Expr {
	return cueast.NewBinExpr(cuetoken.OR, types...)
}

func allOf(types ...cueast.Expr) cueast.Expr {
	return cueast.NewBinExpr(cuetoken.AND, types...)
}

func any() *cueast.Ident {
	return cueast.NewIdent("_")
}

func addComments(node cueast.Node, comments ...*cueast.CommentGroup) {
	for i := range comments {
		cg := comments[i]
		if cg == nil {
			continue
		}
		cueast.AddComment(node, comments[i])
	}
}
