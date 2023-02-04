package echo

func (r *rest) routing() {
	v1 := r.echo.Group("/v1")

	v1.POST("/login", r.adminController.login)
	v1.GET("/users", r.userController.FindAllUsers, r.adminMiddleware.OnlyAdmin)
	v1.GET("/users/:id", r.userController.FindUserByID, r.adminMiddleware.OnlyAdmin)
	v1.POST("/users", r.userController.CreateUser, r.adminMiddleware.OnlyAdmin)
	v1.PUT("/users/:id", r.userController.UpdateUser, r.adminMiddleware.OnlyAdmin)
	v1.DELETE("/users/:id", r.userController.DeleteUser, r.adminMiddleware.OnlyAdmin)
	v1.PATCH("/users/:id/activate", r.userController.ActiveUser, r.adminMiddleware.OnlyAdmin)
	v1.PATCH("/users/:id/disable", r.userController.DisableUser, r.adminMiddleware.OnlyAdmin)
	v1.POST("/servers/:id/heartbeat", r.serverController.HeartBeat, r.adminMiddleware.OnlyAdmin)
}
