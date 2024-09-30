package controllers

func (s *Server) initializeRoutes() {
	v1 := s.Router.Group("/api/v1")
	{
		//CREATE USER ROUTE
		v1.POST("users", s.CreateUser)
		v1.GET("users", s.GetUsers)
		v1.GET("users/:id", s.GetUser)
		v1.PATCH("users/:id", s.UpdateUser)
		v1.POST("users/login", s.Login)
	}

}
