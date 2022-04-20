package codegen_go

import (
	"dog/cfg"
	"dog/storage"
	gast "go/ast"
)

//check if should add star
func (this *Translation) checkStar(src gast.Expr) (dst gast.Expr) {
	if !cfg.StarClassTypeParam {
		return src
	}

	switch tp := src.(type) {
	case *gast.SelectorExpr:
		cl, err := storage.FindByName(tp.Sel.Name)
		if err == nil {
			if cl.Kind == 0 {
				dst = &gast.StarExpr{
					X: tp,
				}
				return dst
			}

		} else if tp.Sel.Name != "BigDecimal" {
			dst = &gast.StarExpr{
				X: tp,
			}
			return dst
		}
	case *gast.Ident:
		if this.currentClass != nil && tp.Name == this.currentClass.GetName() {
			dst = &gast.StarExpr{
				X: tp,
			}
			return dst
		} else {

			_, err := storage.FindByName(tp.Name)
			if err == nil {
				dst = &gast.StarExpr{
					X: tp,
				}
				return dst
			}

		}
	}
	return src
}
