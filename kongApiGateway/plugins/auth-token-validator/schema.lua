return {
  name = "auth-token-validator",
  fields = {
    { config = {
        type = "record",
        fields = {
          { jwt_secret = { 
              type = "string", 
              required = false,
              description = "JWT secret for access token signature verification"
            } 
          },
        },
      },
    },
  },
}