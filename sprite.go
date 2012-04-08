package glutils

func BuildSpriteBuffer() *MeshBuffer {
	build := NewMeshBuilder(4, 6, RENDER_POLYGONS, BUF_TEX_COORD0)
	i0 := build.StartVertex()
	build.AddPosition(-0.5, -0.5, 0)
	build.AddTexCoord(0, 0)

	i1 := build.StartVertex()
	build.AddPosition(0.5, -0.5, 0)
	build.AddTexCoord(1, 0)

	i2 := build.StartVertex()
	build.AddPosition(0.5, 0.5, 0)
	build.AddTexCoord(1, 1)

	i3 := build.StartVertex()
	build.AddPosition(-0.5, 0.5, 0)
	build.AddTexCoord(0, 1)

	build.AddIndice3(i0, i1, i2)
	build.AddIndice3(i2, i3, i0)

	return build.Finalize(true)
}
