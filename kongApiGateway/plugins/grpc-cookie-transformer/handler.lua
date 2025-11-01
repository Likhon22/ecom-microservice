local GrpcCookieTransformer = {
  PRIORITY =1000,
  VERSION= "1.0.0"
}

function GrpcCookieTransformer:header_filter(conf)
  local headers=kong.response.get_headers()
  local cookies={}
  for key, value in ipairs(headers) do
    if key:lower():match("^grpc%-metadata%-set%-cookie") then
      table.insert(cookies, value)
    end
    
  end
    if #cookies > 0 then
    kong.response.set_header("Set-Cookie", cookies)
  end
end
return GrpcCookieTransformer