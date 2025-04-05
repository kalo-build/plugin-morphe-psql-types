package compile

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

// CompiledMorpheViews maps Morphe.Name -> MorpheView.Name -> CompiledView
type CompiledMorpheViews map[string]map[string]CompiledView

func (views CompiledMorpheViews) AddCompiledMorpheView(morpheName string, viewDef *psqldef.View, viewContents []byte) {
	if views[morpheName] == nil {
		views[morpheName] = make(map[string]CompiledView)
	}
	views[morpheName][viewDef.Name] = CompiledView{
		View:         viewDef,
		ViewContents: viewContents,
	}
}

func (views CompiledMorpheViews) GetAllCompiledMorpheViews(morpheName string) map[string]CompiledView {
	morpheViews, morpheViewsExist := views[morpheName]
	if !morpheViewsExist {
		return nil
	}
	return morpheViews
}

func (views CompiledMorpheViews) GetCompiledMorpheView(morpheName string, viewName string) CompiledView {
	morpheViews, morpheViewsExist := views[morpheName]
	if !morpheViewsExist {
		return CompiledView{}
	}
	compiledView, compiledViewExists := morpheViews[viewName]
	if !compiledViewExists {
		return CompiledView{}
	}
	return compiledView
}
