package echo

func (r *rest) routing() {
	v1 := r.echo.Group("/v1")

	v1.POST("/login", r.adminController.login)

}
